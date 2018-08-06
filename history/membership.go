package history

import (
	"fmt"

	"github.com/aalda/trees/common"
)

// TODO unify with incremental -> we need to unify cachedresolvers before

type MembershipProof struct {
	AuditPath common.AuditPath
}

func NewMembershipProof(path common.AuditPath) *MembershipProof {
	return &MembershipProof{path}
}

type AuditPathVisitor struct {
	decorated *common.ComputeHashVisitor
	auditPath common.AuditPath
}

func NewAuditPathVisitor(decorated *common.ComputeHashVisitor) *AuditPathVisitor {
	return &AuditPathVisitor{decorated, make(common.AuditPath)}
}

func (v AuditPathVisitor) Result() common.AuditPath {
	return v.auditPath
}

func (v *AuditPathVisitor) VisitRoot(pos common.Position, leftResult, rightResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitRoot(pos, leftResult, rightResult)
}

func (v *AuditPathVisitor) VisitNode(pos common.Position, leftResult, rightResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitNode(pos, leftResult, rightResult)
}

func (v *AuditPathVisitor) VisitPartialNode(pos common.Position, leftResult interface{}) interface{} {
	// by-pass
	return v.decorated.VisitPartialNode(pos, leftResult)
}

func (v *AuditPathVisitor) VisitLeaf(pos common.Position, eventDigest []byte) interface{} {
	// ignore. target leafs not included in path
	return v.decorated.VisitLeaf(pos, eventDigest)
}

func (v *AuditPathVisitor) VisitCached(pos common.Position) interface{} {
	digest := v.decorated.VisitCached(pos)
	fmt.Printf("Adding cached to path in position: %v\n", pos)
	v.auditPath[pos.StringId()] = digest.(common.Digest)
	return digest
}
