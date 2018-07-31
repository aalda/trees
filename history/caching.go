package history

import (
	"fmt"

	"github.com/aalda/trees/common"
)

type CachingVisitor struct {
	cache     common.Store
	version   uint64
	decorated common.Visitor
}

func NewCachingVisitor(version uint64, cache common.Store, decorated common.Visitor) *CachingVisitor {
	return &CachingVisitor{cache, version, decorated}
}

func (v *CachingVisitor) VisitRoot(pos *common.Position, leftResult, rightResult interface{}) interface{} {
	digest := v.decorated.VisitRoot(pos, leftResult, rightResult).(common.Digest)
	if v.shouldCache(pos) {
		fmt.Printf("Caching node with position: %v\n", pos)
		v.cache.Add(*pos, digest)
	}
	return digest
}

func (v *CachingVisitor) VisitNode(pos *common.Position, leftResult, rightResult interface{}) interface{} {
	digest := v.decorated.VisitNode(pos, leftResult, rightResult).(common.Digest)
	if v.shouldCache(pos) {
		fmt.Printf("Caching node with position: %v\n", pos)
		v.cache.Add(*pos, digest)
	}
	return digest
}

func (v *CachingVisitor) VisitPartialNode(pos *common.Position, leftResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitPartialNode(pos, leftResult)
}

func (v *CachingVisitor) VisitLeaf(pos *common.Position, eventDigest []byte) interface{} {
	digest := v.decorated.VisitLeaf(pos, eventDigest).(common.Digest)
	if v.shouldCache(pos) {
		fmt.Printf("Caching leaf with position: %v\n", pos)
		v.cache.Add(*pos, digest)
	}
	return digest
}

func (v *CachingVisitor) VisitCached(pos *common.Position) interface{} {
	// by-pass
	return v.decorated.VisitCached(pos)
}

func (v *CachingVisitor) shouldCache(pos *common.Position) bool {
	return v.version >= pos.Index+pow(2, pos.Height)-1
}
