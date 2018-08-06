package hyper

import (
	"fmt"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/storage"
)

type CachingVisitor struct {
	cacheLevel uint16
	decorated  common.Visitor
	mutations  []storage.Mutation
}

func NewCachingVisitor(cacheLevel uint16, decorated common.Visitor) *CachingVisitor {
	return &CachingVisitor{
		cacheLevel: cacheLevel,
		decorated:  decorated,
		mutations:  make([]storage.Mutation, 0),
	}
}

func (v *CachingVisitor) Result() []storage.Mutation {
	return v.mutations
}

func (v *CachingVisitor) VisitRoot(pos common.Position, leftResult, rightResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitRoot(pos, leftResult, rightResult).(common.Digest)
}

func (v *CachingVisitor) VisitNode(pos common.Position, leftResult, rightResult interface{}) interface{} {
	digest := v.decorated.VisitNode(pos, leftResult, rightResult).(common.Digest)
	if v.shouldCache(pos) {
		fmt.Printf("Caching node with position: %v\n", pos)
		mutation := storage.NewMutation(storage.HyperCachePrefix, pos.Bytes(), digest)
		v.mutations = append(v.mutations, *mutation)
	}
	return digest
}

func (v *CachingVisitor) VisitPartialNode(pos common.Position, leftResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitPartialNode(pos, leftResult)
}

func (v *CachingVisitor) VisitLeaf(pos common.Position, eventDigest []byte) interface{} {
	digest := v.decorated.VisitLeaf(pos, eventDigest).(common.Digest)
	if v.shouldCache(pos) {
		fmt.Printf("Caching leaf with position: %v\n", pos)
		mutation := storage.NewMutation(storage.HyperCachePrefix, pos.Bytes(), digest)
		v.mutations = append(v.mutations, *mutation)
	}
	return digest
}

func (v *CachingVisitor) VisitCached(pos common.Position) interface{} {
	// by-pass
	return v.decorated.VisitCached(pos)
}

func (v CachingVisitor) shouldCache(pos common.Position) bool {
	return pos.Height() > v.cacheLevel
}
