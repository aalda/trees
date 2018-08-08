package common

import (
	"fmt"
)

type Traversable interface {
	Traverse(pos Position, navigator Navigator, cache Cache) Visitable
}

type Visitor interface {
	VisitRoot(pos Position, leftResult, rightResult interface{}) interface{}
	VisitNode(pos Position, leftResult, rightResult interface{}) interface{}
	VisitPartialNode(pos Position, leftResult interface{}) interface{}
	VisitLeaf(pos Position, value []byte) interface{}
	VisitCached(pos Position, cachedDigest Digest) interface{}
	VisitCacheable(pos Position, result interface{}) interface{}
}

type Visitable interface {
	Accept(visitor Visitor) interface{}
	String() string
}

type Root struct {
	pos         Position
	left, right Visitable
}

type Node struct {
	pos         Position
	left, right Visitable
}

type PartialNode struct {
	pos  Position
	left Visitable
}

type Leaf struct {
	pos   Position
	value []byte
}

type Cached struct {
	pos    Position
	digest Digest
}

type Default struct {
	pos Position
}

type Cacheable struct {
	pos        Position
	underlying Visitable
}

func NewRoot(pos Position, left, right Visitable) *Root {
	return &Root{pos, left, right}
}

func (r Root) Accept(visitor Visitor) interface{} {
	leftResult := r.left.Accept(visitor)
	rightResult := r.right.Accept(visitor)
	return visitor.VisitRoot(r.pos, leftResult, rightResult)
}

func (r Root) String() string {
	return fmt.Sprintf("Root(%v)[ %v | %v ]", r.pos, r.left, r.right)
}

func NewNode(pos Position, left, right Visitable) *Node {
	return &Node{pos, left, right}
}

func (n Node) Accept(visitor Visitor) interface{} {
	leftResult := n.left.Accept(visitor)
	rightResult := n.right.Accept(visitor)
	return visitor.VisitNode(n.pos, leftResult, rightResult)
}

func (n Node) String() string {
	return fmt.Sprintf("Node(%v)[ %v | %v ]", n.pos, n.left, n.right)
}

func NewPartialNode(pos Position, left Visitable) *PartialNode {
	return &PartialNode{pos, left}
}

func (p PartialNode) Accept(visitor Visitor) interface{} {
	leftResult := p.left.Accept(visitor)
	return visitor.VisitPartialNode(p.pos, leftResult)
}

func (p PartialNode) String() string {
	return fmt.Sprintf("PartialNode(%v)[ %v ]", p.pos, p.left)
}

func NewLeaf(pos Position, value []byte) *Leaf {
	return &Leaf{pos, value}
}

func (l Leaf) Accept(visitor Visitor) interface{} {
	return visitor.VisitLeaf(l.pos, l.value)
}

func (l Leaf) String() string {
	return fmt.Sprintf("Leaf(%v)[ %x ]", l.pos, l.value)
}

func NewCached(pos Position, digest Digest) *Cached {
	return &Cached{pos, digest}
}

func (c Cached) Accept(visitor Visitor) interface{} {
	return visitor.VisitCached(c.pos, c.digest)
}

func (c Cached) String() string {
	return fmt.Sprintf("Cached(%v)[ %x ]", c.pos, c.digest)
}

func NewCacheable(pos Position, underlying Visitable) *Cacheable {
	return &Cacheable{pos, underlying}
}

func (c Cacheable) Accept(visitor Visitor) interface{} {
	result := c.underlying.Accept(visitor)
	return visitor.VisitCacheable(c.pos, result)
}

func (c Cacheable) String() string {
	return fmt.Sprintf("Cacheable[ %v ]", c.underlying)
}
