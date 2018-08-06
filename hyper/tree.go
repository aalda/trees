package hyper

import (
	"fmt"
	"sync"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/storage"
	"github.com/aalda/trees/util"
)

type HyperTree struct {
	lock       sync.RWMutex
	store      storage.Store
	cache      storage.Cache
	hasher     common.Hasher
	cacheLevel uint16
}

func NewHyperTree(hasher common.Hasher, store storage.Store, cache storage.Cache, cacheLevel uint16) *HyperTree {
	var lock sync.RWMutex
	return &HyperTree{lock, store, cache, hasher, cacheLevel}
}

func newRootPosition(numBits uint16) common.Position {
	index := make([]byte, numBits/8)
	return NewPosition(index, numBits)
}

func (t *HyperTree) Add(eventDigest common.Digest, version uint64) *common.Commitment {
	t.lock.Lock()
	defer t.lock.Unlock()
	fmt.Printf("Adding event %b with version %d\n", eventDigest, version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher, t.cache)
	caching := NewCachingVisitor(t.cacheLevel, computeHash)

	// navigator
	targetPos := NewPosition(eventDigest, 0)
	navigator := NewHyperNavigator(targetPos, t.hasher.Len(), t.cacheLevel)

	// create a mutation for the new leaf
	versionAsBytes := util.Uint64AsBytes(version)
	leafMutation := storage.NewMutation(storage.IndexPrefix, eventDigest, versionAsBytes)

	// create a leaves range with the new leaf to insert
	leaves := storage.NewKVRange()
	leaf := storage.NewKVPair(eventDigest, versionAsBytes)
	leaves = leaves.InsertSorted(leaf)

	// traverse from root and generate a visitable pruned tree
	traverser := NewHyperTraverser(t.hasher.Len(), t.cacheLevel, leaves, t.store)
	root := traverser.Traverse(newRootPosition(t.hasher.Len()), navigator)

	// visit the pruned tree
	rh := root.Accept(caching).(common.Digest)

	fmt.Println(root)

	// persiste mutations
	mutations := caching.Result()
	mutations = append(mutations, *leafMutation)
	t.store.Mutate(mutations)

	fmt.Println(mutations)

	return common.NewCommitment(version, rh)
}
