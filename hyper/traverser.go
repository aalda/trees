package hyper

import (
	"github.com/aalda/trees/common"
)

type HyperTraverser struct {
	numBits       uint16
	cacheLevel    uint16
	store         common.Store
	defaultHashes []common.Digest
}

func NewHyperTraverser(numBits, cacheLevel uint16, store common.Store, defaultHashes []common.Digest) *HyperTraverser {
	return &HyperTraverser{
		numBits:       numBits,
		cacheLevel:    cacheLevel,
		store:         store,
		defaultHashes: defaultHashes,
	}
}

func (t HyperTraverser) Traverse(pos common.Position, navigator common.Navigator, cache common.Cache, leaves common.KVRange) common.Visitable {

	if navigator.ShouldBeCached(pos) {
		digest, ok := cache.Get(pos)
		if !ok {
			return common.NewCached(pos, t.defaultHashes[pos.Height()])
		}
		return common.NewCached(pos, digest)
	}

	// if we are over the cache level, we need to do a range query to get the leaves
	if pos.Height() < t.cacheLevel {
		first := t.descendToFirst(pos)
		last := t.descendToLast(pos)
		kvRange, _ := t.store.GetRange(common.IndexPrefix, first.Index(), last.Index())

		// replace leaves with new slice and append the previous to the new one
		for _, l := range leaves {
			kvRange = kvRange.InsertSorted(l)
		}
		leaves = kvRange

		return t.Traverse2(pos, navigator, leaves)
	}

	rightPos := navigator.GoToRight(pos)
	leftSlice, rightSlice := leaves.Split(rightPos.Index())
	left := t.Traverse(navigator.GoToLeft(pos), navigator, cache, leftSlice)
	right := t.Traverse(rightPos, navigator, cache, rightSlice)

	if navigator.IsRoot(pos) {
		return common.NewRoot(pos, left, right)
	}

	node := common.NewNode(pos, left, right)
	if navigator.ShouldCache(pos) {
		return common.NewCacheable(pos, node)
	}
	return node
}

func (t HyperTraverser) Traverse2(pos common.Position, navigator common.Navigator, leaves common.KVRange) common.Visitable {
	if navigator.IsLeaf(pos) && len(leaves) == 1 {
		return common.NewLeaf(pos, leaves[0].Value)
	}
	if !navigator.IsRoot(pos) && len(leaves) == 0 {
		return common.NewCached(pos, t.defaultHashes[pos.Height()])
	}
	if len(leaves) > 1 && navigator.IsLeaf(pos) {
		panic("this should never happen (unsorted LeavesSlice or broken split?)")
	}

	// we do a post-order traversal

	// split leaves
	rightPos := navigator.GoToRight(pos)
	leftSlice, rightSlice := leaves.Split(rightPos.Index())
	left := t.Traverse2(navigator.GoToLeft(pos), navigator, leftSlice)
	right := t.Traverse2(rightPos, navigator, rightSlice)
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
	base := make([]byte, t.numBits/8)
	copy(base, pos.Index())
	for bit := layer; bit < t.numBits; bit++ {
		bitSet(base, bit)
	}
	return NewPosition(base, 0)
}
