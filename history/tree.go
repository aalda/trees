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
	computeHash := common.NewComputeHashVisitor(t.hasher, t.cache)
	caching := common.NewCachingVisitor(computeHash)

	// navigator
	targetPos := NewPosition(version, 0)
	resolver := NewMembershipCachedResolver(targetPos)
	navigator := NewHistoryNavigator(resolver, targetPos, targetPos, t.getDepth(version))

	// traverse from root and generate a visitable pruned tree
	traverser := NewHistoryTraverser(eventDigest)
	root := traverser.Traverse(t.newRootPosition(version), navigator, t.cache)

	log.Debugf("Pruned tree: %v", root)

	// visit the pruned tree
	rh := root.PostOrder(caching).(common.Digest)

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
	computeHash := common.NewComputeHashVisitor(t.hasher, t.cache)
	calcAuditPath := common.NewAuditPathVisitor(computeHash)

	// navigator
	startPos := NewPosition(index, 0)
	endPos := NewPosition(version, 0)
	var resolver CachedResolver
	switch index == version {
	case true:
		resolver = NewMembershipCachedResolver(startPos)
	case false:
		resolver = NewIncrementalCachedResolver(startPos, endPos)
	}
	navigator := NewHistoryNavigator(resolver, startPos, endPos, t.getDepth(version))

	// traverse from root and generate a visitable pruned tree
	traverser := NewHistoryTraverser(nil)
	root := traverser.Traverse(t.newRootPosition(version), navigator, t.cache)

	log.Debugf("Pruned tree: %v", root)

	// visit the pruned tree
	root.PostOrder(calcAuditPath)
	return NewMembershipProof(calcAuditPath.Result())
}

func (t *HistoryTree) VerifyMembership(proof *MembershipProof, version uint64, eventDigest, expectedDigest common.Digest) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Debugf("Verifying membership for version %d", version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher, t.cache)

	// navigator
	targetPos := NewPosition(version, 0)
	resolver := NewMembershipCachedResolver(targetPos)
	navigator := NewHistoryNavigator(resolver, targetPos, targetPos, t.getDepth(version))

	// traverse from root and generate a visitable pruned tree
	traverser := NewHistoryTraverser(eventDigest)
	root := traverser.Traverse(t.newRootPosition(version), navigator, proof.AuditPath)

	log.Debugf("Pruned tree: %v", root)

	// visit the pruned tree
	recomputed := root.PostOrder(computeHash).(common.Digest)
	return bytes.Equal(recomputed, expectedDigest)
}

func (t *HistoryTree) ProveConsistency(start, end uint64) *IncrementalProof {
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Debugf("Proving consistency between versions %d and %d", start, end)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher, t.cache)
	calcAuditPath := NewIncAuditPathVisitor(computeHash)

	// navigator
	startPos := NewPosition(start, 0)
	endPos := NewPosition(end, 0)
	resolver := NewIncrementalCachedResolver(startPos, endPos)
	navigator := NewHistoryNavigator(resolver, startPos, endPos, t.getDepth(end))

	// traverse from root and generate a visitable pruned tree
	traverser := NewHistoryTraverser(nil)
	root := traverser.Traverse(t.newRootPosition(end), navigator, t.cache)

	log.Debugf("Pruned tree: %v", root)

	// visit the pruned tree
	root.PostOrder(calcAuditPath)
	return NewIncrementalProof(calcAuditPath.Result())
}
