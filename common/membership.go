package common

import (
	"fmt"
)

// TODO unify with incremental -> we need to unify cachedresolvers before

type AuditPathVisitor struct {
	decorated *ComputeHashVisitor
	auditPath AuditPath
}

func NewAuditPathVisitor(decorated *ComputeHashVisitor) *AuditPathVisitor {
	return &AuditPathVisitor{decorated, make(AuditPath)}
}

func (v AuditPathVisitor) Result() AuditPath {
	return v.auditPath
}

func (v *AuditPathVisitor) VisitRoot(pos Position, leftResult, rightResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitRoot(pos, leftResult, rightResult)
}

func (v *AuditPathVisitor) VisitNode(pos Position, leftResult, rightResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitNode(pos, leftResult, rightResult)
}

func (v *AuditPathVisitor) VisitPartialNode(pos Position, leftResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitPartialNode(pos, leftResult)
}

func (v *AuditPathVisitor) VisitLeaf(pos Position, eventDigest []byte) interface{} {
	// ignore. target leafs not included in path
	return v.decorated.VisitLeaf(pos, eventDigest)
}

func (v *AuditPathVisitor) VisitCached(pos Position, cachedDigest Digest) interface{} {
	digest := v.decorated.VisitCached(pos, cachedDigest)
	fmt.Printf("Adding cached to path in position: %v\n", pos)
	v.auditPath[pos.StringId()] = digest.(Digest)
	return digest
}

func (v *AuditPathVisitor) VisitCacheable(pos Position, result interface{}) interface{} {
	// by-pass
	return v.decorated.VisitCacheable(pos, result)
}
