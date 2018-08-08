package hyper

import (
	"sync"

	"github.com/aalda/trees/common"
	//. "github.com/aalda/trees/logging"
	"github.com/aalda/trees/util"
)

type HyperTree struct {
	lock          sync.RWMutex
	store         common.Store
	cache         common.ModifiableCache
	hasher        common.Hasher
	cacheLevel    uint16
	defaultHashes []common.Digest
}

func NewHyperTree(hasher common.Hasher, store common.Store, cache common.ModifiableCache, cacheLevel uint16) *HyperTree {
	var lock sync.RWMutex
	tree := &HyperTree{
		lock:          lock,
		store:         store,
		cache:         cache,
		hasher:        hasher,
		cacheLevel:    cacheLevel,
		defaultHashes: make([]common.Digest, hasher.Len()),
	}

	tree.defaultHashes[0] = tree.hasher.Do([]byte{0x0}, []byte{0x0})
	for i := uint16(1); i < hasher.Len(); i++ {
		tree.defaultHashes[i] = tree.hasher.Do(tree.defaultHashes[i-1], tree.defaultHashes[i-1])
	}
	return tree
}

func newRootPosition(numBits uint16) common.Position {
	index := make([]byte, numBits/8)
	return NewPosition(index, numBits)
}

func (t *HyperTree) Add(eventDigest common.Digest, version uint64) *common.Commitment {
	t.lock.Lock()
	defer t.lock.Unlock()
	//fmt.Printf("Adding event %b with version %d\n", eventDigest, version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher, t.cache)
	caching := common.NewCachingVisitor(computeHash)

	// navigator
	targetPos := NewPosition(eventDigest, 0)
	navigator := NewHyperNavigator(targetPos, t.hasher.Len(), t.cacheLevel)

	// create a mutation for the new leaf
	versionAsBytes := util.Uint64AsBytes(version)
	leafMutation := common.NewMutation(common.IndexPrefix, eventDigest, versionAsBytes)

	// create a leaves range with the new leaf to insert
	leaves := common.NewKVRange()
	leaf := common.NewKVPair(eventDigest, versionAsBytes)
	leaves = leaves.InsertSorted(leaf)

	// traverse from root and generate a visitable pruned tree
	traverser := NewHyperTraverser(t.hasher.Len(), t.cacheLevel, t.store, t.defaultHashes)
	root := traverser.Traverse(newRootPosition(t.hasher.Len()), navigator, t.cache, leaves)

	// visit the pruned tree
	rh := root.Accept(caching).(common.Digest)

	//Trace.Println(root)

	// persist mutations
	cachedElements := caching.Result()
	mutations := make([]common.Mutation, len(cachedElements))
	for _, e := range cachedElements {
		mutation := common.NewMutation(common.HyperCachePrefix, e.Pos.Bytes(), e.Digest)
		mutations = append(mutations, *mutation)

		// update cache
		t.cache.Put(e.Pos, e.Digest)
	}
	mutations = append(mutations, *leafMutation)
	t.store.Mutate(mutations)

	//fmt.Println(mutations)

	return common.NewCommitment(version, rh)
}
