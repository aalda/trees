package common

import (
	"fmt"
)

// TODO this could be the same ComputeHashVisitor if we abstract the AuditPath interface to make equal to Store
type RecomputeHashVisitor struct {
	decorated *ComputeHashVisitor
	auditPath AuditPath
}

func NewRecomputeHashVisitor(decorated *ComputeHashVisitor, auditPath AuditPath) *RecomputeHashVisitor {
	return &RecomputeHashVisitor{decorated, auditPath}
}

func (v *RecomputeHashVisitor) VisitRoot(pos Position, leftResult, rightResult interface{}) interface{} {
	return v.decorated.VisitRoot(pos, leftResult, rightResult)
}

func (v *RecomputeHashVisitor) VisitNode(pos Position, leftResult, rightResult interface{}) interface{} {
	return v.decorated.VisitNode(pos, leftResult, rightResult)
}

func (v *RecomputeHashVisitor) VisitPartialNode(pos Position, leftResult interface{}) interface{} {
	return v.decorated.VisitPartialNode(pos, leftResult)
}

func (v *RecomputeHashVisitor) VisitLeaf(pos Position, value []byte) interface{} {
	return v.decorated.VisitLeaf(pos, value)
}

func (v *RecomputeHashVisitor) VisitCached(pos Position) interface{} {
	fmt.Printf("Getting hash from path in position: %v\n", pos)
	return v.auditPath[pos.StringId()]
}
