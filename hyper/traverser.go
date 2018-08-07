package hyper

import (
	"github.com/aalda/trees/common"
	"github.com/aalda/trees/storage"
)

type HyperTraverser struct {
	numBits    uint16
	cacheLevel uint16
	leaves     storage.KVRange
	store      storage.Store
}

func NewHyperTraverser(numBits, cacheLevel uint16, leaves storage.KVRange, store storage.Store) *HyperTraverser {
	return &HyperTraverser{
		numBits:    numBits,
		cacheLevel: cacheLevel,
		leaves:     leaves,
		store:      store,
	}
}

func (t HyperTraverser) Traverse(pos common.Position, navigator common.Navigator) common.Visitable {
	if navigator.ShouldBeCached(pos) {
		return common.NewCached(pos)
	}
	if navigator.IsLeaf(pos) && len(t.leaves) == 1 {
		leaf := common.NewLeaf(pos, t.leaves[0].Value)
		if navigator.ShouldCache(pos) {
			return common.NewCacheable(pos, leaf)
		}
		return leaf
	}
	if !navigator.IsRoot(pos) && len(t.leaves) == 0 {
		return common.NewCached(pos) // it should resolve to a default hash because it actually won't be in cache
	}
	if len(t.leaves) > 1 && navigator.IsLeaf(pos) {
		panic("this should never happen (unsorted LeavesSlice or broken split?)")
	}

	// now we are over the cache level so we need to do a range query to get the leaves
	first := t.descendToFirst(pos)
	last := t.descendToLast(pos)
	kvRange, _ := t.store.GetRange(storage.IndexPrefix, first.Index(), last.Index())

	// replace leaves with new slice and append the previous to the new one
	for _, l := range t.leaves {
		kvRange = kvRange.InsertSorted(l)
	}
	t.leaves = kvRange

	return t.Traverse2(pos, navigator)
}

func (t HyperTraverser) Traverse2(pos common.Position, navigator common.Navigator) common.Visitable {
	if navigator.ShouldBeCached(pos) {
		return common.NewCached(pos)
	}
	if navigator.IsLeaf(pos) && len(t.leaves) == 1 {
		leaf := common.NewLeaf(pos, t.leaves[0].Value)
		if navigator.ShouldCache(pos) {
			return common.NewCacheable(pos, leaf)
		}
		return leaf
	}
	if !navigator.IsRoot(pos) && len(t.leaves) == 0 {
		return common.NewCached(pos) // it should resolve to a default hash because it actually won't be in cache
	}
	if len(t.leaves) > 1 && navigator.IsLeaf(pos) {
		panic("this should never happen (unsorted LeavesSlice or broken split?)")
	}

	// we do a post-order traversal

	// split leaves
	rightPos := navigator.GoToRight(pos)
	leftSlice, rightSlice := t.leaves.Split(rightPos.Index())
	left := NewHyperTraverser(t.numBits, t.cacheLevel, leftSlice, t.store).Traverse2(navigator.GoToLeft(pos), navigator)
	right := NewHyperTraverser(t.numBits, t.cacheLevel, rightSlice, t.store).Traverse2(rightPos, navigator)
	if navigator.IsRoot(pos) {
		return common.NewRoot(pos, left, right)
	}
	node := common.NewNode(pos, left, right)
	if navigator.ShouldCache(pos) {
		return common.NewCacheable(pos, node)
	}
	return node
}

func (t HyperTraverser) descendToFirst(pos common.Position) common.Position {
	return NewPosition(pos.Index(), 0)
}

func (t HyperTraverser) descendToLast(pos common.Position) common.Position {
	layer := t.numBits - pos.Height()
	base := make([]byte, max(uint16(1), t.numBits/8))
	copy(base, pos.Index())
	for bit := layer; bit < t.numBits; bit++ {
		bitSet(base, bit)
	}
	return NewPosition(base, 0)
}

func max(x, y uint16) uint16 {
	if x > y {
		return x
	}
	return y
}
