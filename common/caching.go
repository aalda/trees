package common

import (
	"fmt"

	"github.com/aalda/trees/storage"
)

type CachingVisitor struct {
	storePrefix byte
	decorated   Visitor
	mutations   []storage.Mutation
}

func NewCachingVisitor(storePrefix byte, decorated Visitor) *CachingVisitor {
	return &CachingVisitor{
		storePrefix: storePrefix,
		decorated:   decorated,
		mutations:   make([]storage.Mutation, 0),
	}
}

func (v *CachingVisitor) Result() []storage.Mutation {
	return v.mutations
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

func (v *CachingVisitor) VisitCached(pos Position) interface{} {
	// by-pass
	return v.decorated.VisitCached(pos)
}

func (v *CachingVisitor) VisitCacheable(pos Position, result interface{}) interface{} {
	fmt.Printf("Caching digest with position: %v\n", pos)
	mutation := storage.NewMutation(v.storePrefix, pos.Bytes(), result.(Digest))
	v.mutations = append(v.mutations, *mutation)
	return result
}
