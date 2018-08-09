package common

import "github.com/aalda/trees/log"

type Commitment struct {
	Version uint64
	Digest  Digest
}

func NewCommitment(version uint64, digest Digest) *Commitment {
	return &Commitment{Version: version, Digest: digest}
}

type ComputeHashVisitor struct {
	hasher Hasher
	cache  Cache
}

func NewComputeHashVisitor(hasher Hasher, cache Cache) *ComputeHashVisitor {
	return &ComputeHashVisitor{hasher, cache}
}

func (v *ComputeHashVisitor) VisitRoot(pos Position, leftResult, rightResult interface{}) interface{} {
	log.Debugf("Computing root hash in position: %v", pos)
	return v.interiorHash(pos.Bytes(), leftResult.(Digest), rightResult.(Digest))
}

func (v *ComputeHashVisitor) VisitNode(pos Position, leftResult, rightResult interface{}) interface{} {
	log.Debugf("Computing node hash in position: %v", pos)
	return v.interiorHash(pos.Bytes(), leftResult.(Digest), rightResult.(Digest))
}

func (v *ComputeHashVisitor) VisitPartialNode(pos Position, leftResult interface{}) interface{} {
	log.Debugf("Computing partial node hash in position: %v", pos)
	return v.leafHash(pos.Bytes(), leftResult.(Digest))
}

func (v *ComputeHashVisitor) VisitLeaf(pos Position, value []byte) interface{} {
	log.Debugf("Computing leaf hash in position: %v", pos)
	return v.leafHash(pos.Bytes(), value)
}

func (v *ComputeHashVisitor) VisitCached(pos Position, cachedDigest Digest) interface{} {
	log.Debugf("Getting cached hash in position: %v", pos)
	return cachedDigest
}

func (v *ComputeHashVisitor) VisitCacheable(pos Position, result interface{}) interface{} {
	log.Debugf("Getting cacheable value in position: %v", pos)
	return result
}

func (v *ComputeHashVisitor) leafHash(id, leaf Digest) Digest {
	return v.hasher.Do(leaf)
}

func (v *ComputeHashVisitor) interiorHash(id, left, right Digest) Digest {
	return v.hasher.Do(left, right)
}
