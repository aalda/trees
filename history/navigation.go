package history

import (
	"math"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/util"
)

type CachedResolver interface {
	ShouldBeCached(*common.Position) bool
}

type MembershipCachedResolver struct {
	target *common.Position
}

func NewMembershipCachedResolver(target *common.Position) *MembershipCachedResolver {
	return &MembershipCachedResolver{target}
}

func (r MembershipCachedResolver) ShouldBeCached(pos *common.Position) bool {
	return r.target.IndexAsUint64() > pos.IndexAsUint64()+pow(2, pos.Height)-1
}

type IncrementalCachedResolver struct {
	start *common.Position
	end   *common.Position
}

func NewIncrementalCachedResolver(start, end *common.Position) *IncrementalCachedResolver {
	return &IncrementalCachedResolver{start, end}
}

func (r IncrementalCachedResolver) ShouldBeCached(pos *common.Position) bool {
	if pos.Height == 0 && pos.IndexAsUint64() == r.start.IndexAsUint64() { // TODO THIS SHOULD BE TRUE for inc proofs but not for membership
		return false
	}
	threshold := pos.IndexAsUint64() + pow(2, pos.Height) - 1
	if r.start.IndexAsUint64() > threshold && r.end.IndexAsUint64() > threshold {
		return true
	}

	lastDescendantIndex := pos.IndexAsUint64() + pow(2, pos.Height) - 1
	return pos.IndexAsUint64() > r.start.IndexAsUint64() && lastDescendantIndex <= r.end.IndexAsUint64()
}

func pow(x, y uint16) uint64 {
	return uint64(math.Pow(float64(x), float64(y)))
}

type HistoryNavigator struct {
	resolver CachedResolver
	start    *common.Position
	end      *common.Position
	depth    uint16
}

func NewHistoryNavigator(resolver CachedResolver, start, end *common.Position, depth uint16) *HistoryNavigator {
	return &HistoryNavigator{resolver, start, end, depth}
}

func (n *HistoryNavigator) GoToLeft(pos *common.Position) *common.Position {
	if pos.Height == 0 {
		return nil
	}
	return common.NewPosition(pos.Index, pos.Height-1)
}

func (n *HistoryNavigator) GoToRight(pos *common.Position) *common.Position {
	rightIndex := pos.IndexAsUint64() + pow(2, pos.Height-1)
	if pos.Height == 0 || rightIndex > n.end.IndexAsUint64() {
		return nil
	}
	return common.NewPosition(util.Uint64AsBytes(rightIndex), pos.Height-1)
}

func (n *HistoryNavigator) IsLeaf(pos *common.Position) bool {
	return pos.Height == 0
}

func (n *HistoryNavigator) IsRoot(pos *common.Position) bool {
	return pos.Height == n.depth
}

func (n *HistoryNavigator) ShouldBeCached(pos *common.Position) bool {
	return n.resolver.ShouldBeCached(pos)
}
