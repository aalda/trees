package hyper

import (
	"sync"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/log"
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

	log.Debugf("Adding event %b with version %d\n", eventDigest, version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher)
	caching := common.NewCachingVisitor(computeHash)

	// build pruning context
	versionAsBytes := util.Uint64AsBytes(version)
	context := PruningContext{
		navigator:     NewHyperTreeNavigator(t.hasher.Len()),
		cacheResolver: NewSingleTargetedCacheResolver(t.hasher.Len(), t.cacheLevel, eventDigest),
		cache:         t.cache,
		store:         t.store,
		defaultHashes: t.defaultHashes,
		cacheLevel:    t.cacheLevel,
	}

	// traverse from root and generate a visitable pruned tree
	pruned := NewInsertPruner(eventDigest, versionAsBytes, context).Prune()

	// print := common.NewPrintVisitor(t.hasher.Len())
	// pruned.PreOrder(print)
	// log.Debugf("Pruned tree: %s", print.Result())

	// visit the pruned tree
	rh := pruned.PostOrder(caching).(common.Digest)

	// persist mutations
	cachedElements := caching.Result()
	mutations := make([]common.Mutation, len(cachedElements))
	for i, e := range cachedElements {
		mutations[i] = *common.NewMutation(common.HyperCachePrefix, e.Pos.Bytes(), e.Digest)
		// update cache
		t.cache.Put(e.Pos, e.Digest)
	}
	// create a mutation for the new leaf
	leafMutation := common.NewMutation(common.IndexPrefix, eventDigest, versionAsBytes)
	mutations = append(mutations, *leafMutation)
	t.store.Mutate(mutations) // TODO the mutations should be returned and persited at the balloon level

	log.Debugf("Mutations: %v", mutations)

	return common.NewCommitment(version, rh)
}

type MembershipProof struct {
	AuditPath common.AuditPath
}

func (t *HyperTree) Get(eventDigest common.Digest) (version uint64, proof *MembershipProof) {
	t.lock.Lock()
	defer t.lock.Unlock()

	log.Debugf("Getting version for event %b\n", eventDigest)

	// visitors
	//computeHash := common.NewComputeHashVisitor(t.hasher, t.cache)

	// navigator
	//targetPos := NewPosition(eventDigest, 0)
	//navigator := NewHyperNavigator(targetPos, t.hasher.Len(), t.cacheLevel)

	// traverse from root and generate a visitable pruned tree
	//traverser := NewHyperTraverser(t.hasher.Len(), t.cacheLevel, t.store, t.defaultHashes)
	//pruned := traverser.TraverseAudit(newRootPosition(t.hasher.Len()), navigator, t.cache)

	return 0, nil
}
