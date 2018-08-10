package hyper

import (
	"github.com/aalda/trees/common"
)

type PruningContext struct {
	navigator     common.TreeNavigator
	cacheResolver CacheResolver
	cache         common.Cache
	store         common.Store
	defaultHashes []common.Digest
}

type Pruner interface {
	Prune() common.Visitable
}

type InsertPruner struct {
	key   common.Digest
	value []byte
	PruningContext
}

func NewInsertPruner(key, value []byte, context PruningContext) *InsertPruner {
	return &InsertPruner{key, value, context}
}

func (p *InsertPruner) Prune() common.Visitable {
	leaves := common.KVRange{common.NewKVPair(p.key, p.value)}
	return p.traverse(p.navigator.Root(), leaves)
}

func (p *InsertPruner) traverse(pos common.Position, leaves common.KVRange) common.Visitable {
	if p.cacheResolver.ShouldBeInCache(pos) {
		digest, ok := p.cache.Get(pos)
		if !ok {
			return common.NewCached(pos, p.defaultHashes[pos.Height()])
		}
		return common.NewCached(pos, digest)
	}

	// if we are over the cache level, we need to do a range query to get the leaves
	if !p.cacheResolver.ShouldCache(pos) {
		first := p.navigator.DescendToFirst(pos)
		last := p.navigator.DescendToLast(pos)
		kvRange, _ := p.store.GetRange(common.IndexPrefix, first.Index(), last.Index())

		// replace leaves with new slice and append the previous to the new one
		for _, l := range leaves {
			kvRange = kvRange.InsertSorted(l)
		}
		leaves = kvRange

		return p.traverseWithoutCache(pos, leaves)
	}

	rightPos := p.navigator.GoToRight(pos)
	leftSlice, rightSlice := leaves.Split(rightPos.Index())
	left := p.traverse(p.navigator.GoToLeft(pos), leftSlice)
	right := p.traverse(rightPos, rightSlice)

	if p.navigator.IsRoot(pos) {
		return common.NewRoot(pos, left, right)
	}
	return common.NewCacheable(pos, common.NewNode(pos, left, right))
}

func (p *InsertPruner) traverseWithoutCache(pos common.Position, leaves common.KVRange) common.Visitable {
	if p.navigator.IsLeaf(pos) && len(leaves) == 1 {
		return common.NewLeaf(pos, leaves[0].Value)
	}
	if !p.navigator.IsRoot(pos) && len(leaves) == 0 {
		return common.NewCached(pos, p.defaultHashes[pos.Height()])
	}
	if len(leaves) > 1 && p.navigator.IsLeaf(pos) {
		panic("this should never happen (unsorted LeavesSlice or broken split?)")
	}

	// we do a post-order traversal

	// split leaves
	rightPos := p.navigator.GoToRight(pos)
	leftSlice, rightSlice := leaves.Split(rightPos.Index())
	left := p.traverseWithoutCache(p.navigator.GoToLeft(pos), leftSlice)
	right := p.traverseWithoutCache(rightPos, rightSlice)

	if p.navigator.IsRoot(pos) {
		return common.NewRoot(pos, left, right)
	}
	return common.NewNode(pos, left, right)
}

type SearchPruner struct {
	key []byte
	PruningContext
}

func NewSearchPruner(key []byte, context PruningContext) *SearchPruner {
	return &SearchPruner{key, context}
}

func (p *SearchPruner) Prune() common.Visitable {
	return p.traverse(p.navigator.Root(), common.NewKVRange())
}

func (p *SearchPruner) traverse(pos common.Position, leaves common.KVRange) common.Visitable {
	if p.cacheResolver.ShouldBeInCache(pos) {
		digest, ok := p.cache.Get(pos)
		if !ok {
			cached := common.NewCached(pos, p.defaultHashes[pos.Height()])
			return common.NewCacheable(pos, cached)
		}
		return common.NewCacheable(pos, common.NewCached(pos, digest))
	}

	// if we are over the cache level, we need to do a range query to get the leaves
	if !p.cacheResolver.ShouldCache(pos) {
		first := p.navigator.DescendToFirst(pos)
		last := p.navigator.DescendToLast(pos)
		kvRange, _ := p.store.GetRange(common.IndexPrefix, first.Index(), last.Index())

		// replace leaves with new slice and append the previous to the new one
		for _, l := range leaves {
			kvRange = kvRange.InsertSorted(l)
		}
		leaves = kvRange
		return p.traverseWithoutCache(pos, leaves)
	}

	rightPos := p.navigator.GoToRight(pos)
	leftSlice, rightSlice := leaves.Split(rightPos.Index())
	left := p.traverse(p.navigator.GoToLeft(pos), leftSlice)
	right := p.traverse(rightPos, rightSlice)

	if p.navigator.IsRoot(pos) {
		return common.NewRoot(pos, left, right)
	}

	return common.NewNode(pos, left, right)

}

func (p *SearchPruner) traverseWithoutCache(pos common.Position, leaves common.KVRange) common.Visitable {
	if p.navigator.IsLeaf(pos) && len(leaves) == 1 {
		return common.NewLeaf(pos, leaves[0].Value)
	}
	if !p.navigator.IsRoot(pos) && len(leaves) == 0 {
		cached := common.NewCached(pos, p.defaultHashes[pos.Height()])
		return common.NewCacheable(pos, cached)
	}
	if len(leaves) > 1 && p.navigator.IsLeaf(pos) {
		panic("this should never happen (unsorted LeavesSlice or broken split?)")
	}

	// we do a post-order traversal

	// split leaves
	rightPos := p.navigator.GoToRight(pos)
	leftSlice, rightSlice := leaves.Split(rightPos.Index())
	left := p.traverseWithoutCache(p.navigator.GoToLeft(pos), leftSlice)
	right := p.traverseWithoutCache(rightPos, rightSlice)
	if p.navigator.IsRoot(pos) {
		return common.NewRoot(pos, left, right)
	}
	return common.NewNode(pos, left, right)
}
