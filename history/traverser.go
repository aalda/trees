package history

import "github.com/aalda/trees/common"

type HistoryTraverser struct {
	eventDigest common.Digest
}

func NewHistoryTraverser(eventDigest common.Digest) *HistoryTraverser {
	return &HistoryTraverser{eventDigest}
}

func (t HistoryTraverser) Traverse(pos common.Position, navigator common.Navigator) common.Visitable {
	if navigator.ShouldBeCached(pos) {
		return common.NewCached(pos)
	}
	if navigator.IsLeaf(pos) {
		leaf := common.NewLeaf(pos, t.eventDigest)
		if navigator.ShouldCache(pos) {
			return common.NewCacheable(pos, leaf)
		}
		return leaf
	}
	// we do a post-order traversal
	left := t.Traverse(navigator.GoToLeft(pos), navigator)
	rightPos := navigator.GoToRight(pos)
	if rightPos == nil {
		return common.NewPartialNode(pos, left)
	}
	right := t.Traverse(rightPos, navigator)
	var result common.Visitable
	if navigator.IsRoot(pos) {
		result = common.NewRoot(pos, left, right)
	} else {
		result = common.NewNode(pos, left, right)
	}
	if navigator.ShouldCache(pos) {
		return common.NewCacheable(pos, result)
	}
	return result
}
