package history

import (
	"fmt"
	"trees/common"
)

type IncrementalProof struct {
	AuditPath common.AuditPath
}

func NewIncrementalProof(path common.AuditPath) *IncrementalProof {
	return &IncrementalProof{path}
}

type IncAuditPathVisitor struct {
	decorated common.Visitor
	auditPath common.AuditPath
}

func NewIncAuditPathVisitor(decorated common.Visitor) *IncAuditPathVisitor {
	return &IncAuditPathVisitor{decorated, make(common.AuditPath)}
}

func (v IncAuditPathVisitor) Result() common.AuditPath {
	return v.auditPath
}

func (v *IncAuditPathVisitor) VisitRoot(pos *common.Position, leftResult, rightResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitNode(pos, leftResult, rightResult)
}

func (v *IncAuditPathVisitor) VisitNode(pos *common.Position, leftResult, rightResult interface{}) interface{} {
	return v.decorated.VisitNode(pos, leftResult, rightResult)
}

func (v *IncAuditPathVisitor) VisitPartialNode(pos *common.Position, leftResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitPartialNode(pos, leftResult)
}

func (v *IncAuditPathVisitor) VisitLeaf(pos *common.Position, eventDigest []byte) interface{} {
	// the leaf should be in cache
	return v.VisitCached(pos)
}

func (v *IncAuditPathVisitor) VisitCached(pos *common.Position) interface{} {
	digest := v.decorated.VisitCached(pos)
	fmt.Printf("Adding cached to path in position: %v\n", pos)
	v.auditPath[pos.StringId()] = digest.(common.Digest)
	return digest
}
