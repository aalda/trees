package common

import (
	"fmt"

	"github.com/aalda/trees/storage"
)

type Commitment struct {
	Version uint64
	Digest  Digest
}

func NewCommitment(version uint64, digest Digest) *Commitment {
	return &Commitment{Version: version, Digest: digest}
}

type ComputeHashVisitor struct {
	hasher Hasher
	cache  storage.Cache
}

func NewComputeHashVisitor(hasher Hasher, cache storage.Cache) *ComputeHashVisitor {
	return &ComputeHashVisitor{hasher, cache}
}

func (v *ComputeHashVisitor) VisitRoot(pos Position, leftResult, rightResult interface{}) interface{} {
	fmt.Printf("Computing root hash in position: %v\n", pos)
	return v.interiorHash(pos.Bytes(), leftResult.(Digest), rightResult.(Digest))
}

func (v *ComputeHashVisitor) VisitNode(pos Position, leftResult, rightResult interface{}) interface{} {
	fmt.Printf("Computing node hash in position: %v\n", pos)
	return v.interiorHash(pos.Bytes(), leftResult.(Digest), rightResult.(Digest))
}

func (v *ComputeHashVisitor) VisitPartialNode(pos Position, leftResult interface{}) interface{} {
	fmt.Printf("Computing partial node hash in position: %v\n", pos)
	return v.leafHash(pos.Bytes(), leftResult.(Digest))
}

func (v *ComputeHashVisitor) VisitLeaf(pos Position, value []byte) interface{} {
	fmt.Printf("Computing leaf hash in position: %v\n", pos)
	return v.leafHash(pos.Bytes(), value)
}

func (v *ComputeHashVisitor) VisitCached(pos Position) interface{} {
	fmt.Printf("Getting cached hash in position: %v\n", pos)
	pair, _ := v.cache.Get(pos.Bytes())
	return Digest(pair.Value)
}

func (v *ComputeHashVisitor) leafHash(id, leaf Digest) Digest {
	return v.hasher.Do(leaf)
}

func (v *ComputeHashVisitor) interiorHash(id, left, right Digest) Digest {
	return v.hasher.Do(left, right)
}
