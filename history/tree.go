package history

import (
	"bytes"
	"math"
	"sync"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/log"
)

type HistoryTree struct {
	lock   sync.RWMutex
	frozen common.Store
	cache  common.Cache
	hasher common.Hasher
}

func NewHistoryTree(hasher common.Hasher, frozen common.Store, cache common.Cache) *HistoryTree {
	var lock sync.RWMutex
	return &HistoryTree{lock, frozen, cache, hasher}
}

func (t *HistoryTree) newRootPosition(version uint64) *HistoryPosition {
	return NewPosition(0, t.getDepth(version))
}

func (t *HistoryTree) getDepth(version uint64) uint16 {
	return uint16(uint64(math.Ceil(math.Log2(float64(version + 1)))))
}

func (t *HistoryTree) Add(eventDigest common.Digest, version uint64) *common.Commitment {
	t.lock.Lock()
	defer t.lock.Unlock()

	log.Debugf("Adding event %b with version %d\n", eventDigest, version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher)
	caching := common.NewCachingVisitor(computeHash)

	// build pruning context
	context := PruningContext{
		navigator:     NewHistoryTreeNavigator(version),
		cacheResolver: NewSingleTargetedCacheResolver(version),
		cache:         t.cache,
	}

	// traverse from root and generate a visitable pruned tree
	pruned := NewInsertPruner(eventDigest, context).Prune()

	print := common.NewPrintVisitor(t.getDepth(version))
	pruned.PreOrder(print)
	log.Debugf("Pruned tree: %s", print.Result())

	// visit the pruned tree
	rh := pruned.PostOrder(caching).(common.Digest)

	// persist mutations
	cachedElements := caching.Result()
	mutations := make([]common.Mutation, len(cachedElements))
	for _, e := range cachedElements {
		mutation := common.NewMutation(common.HistoryCachePrefix, e.Pos.Bytes(), e.Digest)
		mutations = append(mutations, *mutation)
	}
	t.frozen.Mutate(mutations)

	return common.NewCommitment(version, rh)
}

type MembershipProof struct {
	AuditPath common.AuditPath
}

func NewMembershipProof(path common.AuditPath) *MembershipProof {
	return &MembershipProof{path}
}

func (t *HistoryTree) ProveMembership(index, version uint64) *MembershipProof {
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Debugf("Proving membership for index %d with version %d", index, version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher)
	calcAuditPath := common.NewAuditPathVisitor(computeHash)

	// build pruning context
	var resolver CacheResolver
	switch index == version {
	case true:
		resolver = NewSingleTargetedCacheResolver(version)
	case false:
		resolver = NewDoubleTargetedCacheResolver(index, version)
	}
	context := PruningContext{
		navigator:     NewHistoryTreeNavigator(version),
		cacheResolver: resolver,
		cache:         t.cache,
	}

	// traverse from root and generate a visitable pruned tree
	pruned := NewSearchPruner(context).Prune()

	print := common.NewPrintVisitor(t.getDepth(version))
	pruned.PreOrder(print)
	log.Debugf("Pruned tree: %s", print.Result())

	// visit the pruned tree
	pruned.PostOrder(calcAuditPath)
	return NewMembershipProof(calcAuditPath.Result())
}

func (t *HistoryTree) VerifyMembership(proof *MembershipProof, version uint64, eventDigest, expectedDigest common.Digest) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Debugf("Verifying membership for version %d", version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher)

	// build pruning context
	context := PruningContext{
		navigator:     NewHistoryTreeNavigator(version),
		cacheResolver: NewSingleTargetedCacheResolver(version),
		cache:         proof.AuditPath,
	}

	// traverse from root and generate a visitable pruned tree
	pruned := NewVerifyPruner(eventDigest, context).Prune()

	print := common.NewPrintVisitor(t.getDepth(version))
	pruned.PreOrder(print)
	log.Debugf("Pruned tree: %s", print.Result())

	// visit the pruned tree
	recomputed := pruned.PostOrder(computeHash).(common.Digest)
	return bytes.Equal(recomputed, expectedDigest)
}

type IncrementalProof struct {
	AuditPath common.AuditPath
}

func NewIncrementalProof(path common.AuditPath) *IncrementalProof {
	return &IncrementalProof{path}
}

func (t *HistoryTree) ProveConsistency(start, end uint64) *IncrementalProof {
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Debugf("Proving consistency between versions %d and %d", start, end)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher)
	calcAuditPath := common.NewAuditPathVisitor(computeHash)

	// build pruning context
	context := PruningContext{
		navigator:     NewHistoryTreeNavigator(end),
		cacheResolver: NewIncrementalCacheResolver(start, end),
		cache:         t.cache,
	}

	// traverse from root and generate a visitable pruned tree
	pruned := NewSearchPruner(context).Prune()

	print := common.NewPrintVisitor(t.getDepth(end))
	pruned.PreOrder(print)
	log.Debugf("Pruned tree: %s", print.Result())

	// visit the pruned tree
	pruned.PostOrder(calcAuditPath)
	return NewIncrementalProof(calcAuditPath.Result())
}
