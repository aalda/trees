package history

import (
	"bytes"
	"fmt"
	"math"
	"sync"
	"trees/common"
)

type HistoryTree struct {
	lock   sync.RWMutex
	frozen common.Store
	hasher common.Hasher
}

func NewHistoryTree(hasher common.Hasher, frozen common.Store) *HistoryTree {
	var lock sync.RWMutex
	return &HistoryTree{lock, frozen, hasher}
}

func (t *HistoryTree) newRootPosition(version uint64) *common.Position {
	return &common.Position{0, t.getDepth(version)}
}

func (t *HistoryTree) getDepth(version uint64) uint64 {
	return uint64(math.Ceil(math.Log2(float64(version + 1))))
}

func (t *HistoryTree) Add(eventDigest common.Digest, version uint64) *common.Commitment {
	t.lock.Lock()
	defer t.lock.Unlock()
	fmt.Printf("Adding event %b with version %d\n", eventDigest, version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher, t.frozen)
	caching := NewCachingVisitor(version, t.frozen, computeHash)

	// navigator
	targetPos := common.NewPosition(version, 0)
	resolver := NewMembershipCachedResolver(targetPos)
	navigator := NewHistoryNavigator(resolver, targetPos, targetPos, t.getDepth(version))

	// traverse from root and generate a visitable pruned tree
	root := common.Traverse(t.newRootPosition(version), navigator, eventDigest)
	fmt.Printf("Pruned tree: %v\n", root)

	// visit the pruned tree
	rh := root.Accept(caching).(common.Digest)
	return common.NewCommitment(version, rh)
}

func (t *HistoryTree) ProveMembership(index, version uint64) *MembershipProof {
	t.lock.Lock()
	defer t.lock.Unlock()
	fmt.Printf("Proving membership for index %d with version %d\n", index, version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher, t.frozen)
	calcAuditPath := NewAuditPathVisitor(computeHash)

	// navigator
	startPos := common.NewPosition(index, 0)
	endPos := common.NewPosition(version, 0)
	var resolver CachedResolver
	switch index == version {
	case true:
		resolver = NewMembershipCachedResolver(startPos)
	case false:
		resolver = NewIncrementalCachedResolver(startPos, endPos)
	}
	navigator := NewHistoryNavigator(resolver, startPos, endPos, t.getDepth(version))

	// traverse from root and generate a visitable pruned tree
	root := common.Traverse(t.newRootPosition(version), navigator, nil)
	fmt.Printf("Pruned tree: %v\n", root)

	// visit the pruned tree
	root.Accept(calcAuditPath)
	return NewMembershipProof(calcAuditPath.Result())
}

func (t *HistoryTree) VerifyMembership(proof *MembershipProof, version uint64, eventDigest, expectedDigest common.Digest) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	fmt.Printf("Verifying membership for version %d\n", version)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher, t.frozen)
	recomputeHash := common.NewRecomputeHashVisitor(computeHash, proof.AuditPath)

	// navigator
	targetPos := common.NewPosition(version, 0)
	resolver := NewMembershipCachedResolver(targetPos)
	navigator := NewHistoryNavigator(resolver, targetPos, targetPos, t.getDepth(version))

	// traverse from root and generate a visitable pruned tree
	root := common.Traverse(t.newRootPosition(version), navigator, eventDigest)
	fmt.Printf("Pruned tree: %v\n", root)

	// visit the pruned tree
	recomputed := root.Accept(recomputeHash).(common.Digest)
	return bytes.Equal(recomputed, expectedDigest)
}

func (t *HistoryTree) ProveConsistency(start, end uint64) *IncrementalProof {
	t.lock.Lock()
	defer t.lock.Unlock()
	fmt.Printf("Proving consistency between versions %d and %d\n", start, end)

	// visitors
	computeHash := common.NewComputeHashVisitor(t.hasher, t.frozen)
	calcAuditPath := NewIncAuditPathVisitor(computeHash)

	// navigator
	startPos := common.NewPosition(start, 0)
	endPos := common.NewPosition(end, 0)
	resolver := NewIncrementalCachedResolver(startPos, endPos)
	navigator := NewHistoryNavigator(resolver, startPos, endPos, t.getDepth(end))

	// traverse from root and generate a visitable pruned tree
	root := common.Traverse(t.newRootPosition(end), navigator, nil)
	fmt.Printf("Pruned tree: %v\n", root)

	// visit the pruned tree
	root.Accept(calcAuditPath)
	return NewIncrementalProof(calcAuditPath.Result())
}
