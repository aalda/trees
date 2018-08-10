package common

import "github.com/aalda/trees/log"

type AuditPath map[string]Digest

func (p AuditPath) Get(pos Position) (Digest, bool) {
	digest, ok := p[pos.StringId()]
	return digest, ok
}

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
	// by-pass
	return v.decorated.VisitCached(pos, cachedDigest)
}

func (v *AuditPathVisitor) VisitCacheable(pos Position, result interface{}) interface{} {
	digest := v.decorated.VisitCacheable(pos, result)
	log.Debugf("Adding cacheable to path in position: %v", pos)
	v.auditPath[pos.StringId()] = digest.(Digest)
	return digest
}
