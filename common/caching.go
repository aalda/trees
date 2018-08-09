package common

import "github.com/aalda/trees/log"

type CachedElement struct {
	Pos    Position
	Digest Digest
}

type CachingVisitor struct {
	decorated PostOrderVisitor
	elements  []CachedElement
}

func NewCachingVisitor(decorated PostOrderVisitor) *CachingVisitor {
	return &CachingVisitor{
		decorated: decorated,
		elements:  make([]CachedElement, 0),
	}
}

func (v *CachingVisitor) Result() []CachedElement {
	return v.elements
}

func (v *CachingVisitor) VisitRoot(pos Position, leftResult, rightResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitRoot(pos, leftResult, rightResult).(Digest)
}

func (v *CachingVisitor) VisitNode(pos Position, leftResult, rightResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitNode(pos, leftResult, rightResult).(Digest)
}

func (v *CachingVisitor) VisitPartialNode(pos Position, leftResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitPartialNode(pos, leftResult)
}

func (v *CachingVisitor) VisitLeaf(pos Position, eventDigest []byte) interface{} {
	// by-pass
	return v.decorated.VisitLeaf(pos, eventDigest).(Digest)
}

func (v *CachingVisitor) VisitCached(pos Position, cachedDigest Digest) interface{} {
	// by-pass
	return v.decorated.VisitCached(pos, cachedDigest)
}

func (v *CachingVisitor) VisitCacheable(pos Position, result interface{}) interface{} {
	log.Debugf("Caching digest with position: %v", pos)
	element := &CachedElement{pos, result.(Digest)}
	v.elements = append(v.elements, *element)
	return result
}
