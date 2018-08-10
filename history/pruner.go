package history

import (
	"github.com/aalda/trees/common"
)

type PruningContext struct {
	navigator     common.TreeNavigator
	cacheResolver CacheResolver
	cache         common.Cache
}

type Pruner interface {
	Prune() common.Visitable
}

type InsertPruner struct {
	eventDigest common.Digest
	PruningContext
}

func NewInsertPruner(eventDigest common.Digest, context PruningContext) *InsertPruner {
	return &InsertPruner{eventDigest, context}
}

func (p *InsertPruner) Prune() common.Visitable {
	return p.traverse(p.navigator.Root(), p.eventDigest)
}

func (p *InsertPruner) traverse(pos common.Position, eventDigest common.Digest) common.Visitable {
	if p.cacheResolver.ShouldBeInCache(pos) {
		digest, ok := p.cache.Get(pos)
		if !ok {
			panic("this digest should be in cache")
		}
		return common.NewCached(pos, digest)
	}
	if p.navigator.IsLeaf(pos) {
		leaf := common.NewLeaf(pos, eventDigest)
		if p.cacheResolver.ShouldCache(pos) {
			return common.NewCacheable(pos, leaf)
		}
		return leaf
	}
	// we do a post-order traversal
	left := p.traverse(p.navigator.GoToLeft(pos), eventDigest)
	rightPos := p.navigator.GoToRight(pos)
	if rightPos == nil {
		return common.NewPartialNode(pos, left)
	}
	right := p.traverse(rightPos, eventDigest)
	var result common.Visitable
	if p.navigator.IsRoot(pos) {
		result = common.NewRoot(pos, left, right)
	} else {
		result = common.NewNode(pos, left, right)
	}
	if p.cacheResolver.ShouldCache(pos) {
		return common.NewCacheable(pos, result)
	}
	return result
}

type SearchPruner struct {
	PruningContext
}

func NewSearchPruner(context PruningContext) *SearchPruner {
	return &SearchPruner{context}
}

func (p *SearchPruner) Prune() common.Visitable {
	return p.traverse(p.navigator.Root())
}

func (p *SearchPruner) traverse(pos common.Position) common.Visitable {
	if p.cacheResolver.ShouldBeInCache(pos) {
		digest, ok := p.cache.Get(pos)
		if !ok {
			panic("this digest should be in cache")
		}
		return common.NewCacheable(pos, common.NewCached(pos, digest))
	}
	if p.navigator.IsLeaf(pos) {
		return common.NewLeaf(pos, nil)
	}
	// we do a post-order traversal
	left := p.traverse(p.navigator.GoToLeft(pos))
	rightPos := p.navigator.GoToRight(pos)
	if rightPos == nil {
		return common.NewPartialNode(pos, left)
	}
	right := p.traverse(rightPos)
	if p.navigator.IsRoot(pos) {
		return common.NewRoot(pos, left, right)
	}
	return common.NewNode(pos, left, right)
}

type VerifyPruner struct {
	eventDigest common.Digest
	PruningContext
}

func NewVerifyPruner(eventDigest common.Digest, context PruningContext) *VerifyPruner {
	return &VerifyPruner{eventDigest, context}
}

func (p *VerifyPruner) Prune() common.Visitable {
	return p.traverse(p.navigator.Root(), p.eventDigest)
}

func (p *VerifyPruner) traverse(pos common.Position, eventDigest common.Digest) common.Visitable {
	if p.cacheResolver.ShouldBeInCache(pos) {
		digest, ok := p.cache.Get(pos)
		if !ok {
			panic("this digest should be in cache")
		}
		return common.NewCached(pos, digest)
	}
	if p.navigator.IsLeaf(pos) {
		return common.NewLeaf(pos, eventDigest)
	}
	// we do a post-order traversal
	left := p.traverse(p.navigator.GoToLeft(pos), eventDigest)
	rightPos := p.navigator.GoToRight(pos)
	if rightPos == nil {
		return common.NewPartialNode(pos, left)
	}
	right := p.traverse(rightPos, eventDigest)
	if p.navigator.IsRoot(pos) {
		return common.NewRoot(pos, left, right)
	}
	return common.NewNode(pos, left, right)

}
