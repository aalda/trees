package common

import (
	"fmt"
)

type Traversable interface {
	Traverse(pos Position) Visitable
}

type Visitor interface {
	VisitRoot(pos Position, leftResult, rightResult interface{}) interface{}
	VisitNode(pos Position, leftResult, rightResult interface{}) interface{}
	VisitPartialNode(pos Position, leftResult interface{}) interface{}
	VisitLeaf(pos Position, value []byte) interface{}
	VisitCached(pos Position) interface{}
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
	pos         Position
	eventDigest Digest
}

type Cached struct {
	pos Position
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
	return visitor.VisitLeaf(l.pos, l.eventDigest)
}

func (l Leaf) String() string {
	return fmt.Sprintf("Leaf(%v)", l.pos)
}

func NewCached(pos Position) *Cached {
	return &Cached{pos}
}

func (c Cached) Accept(visitor Visitor) interface{} {
	return visitor.VisitCached(c.pos)
}

func (c Cached) String() string {
	return fmt.Sprintf("Cached(%v)", c.pos)
}

func Traverse(pos Position, navigator Navigator, eventDigest Digest) Visitable {
	if navigator.ShouldBeCached(pos) {
		return NewCached(pos)
	}
	if navigator.IsLeaf(pos) {
		return NewLeaf(pos, eventDigest)
	}
	// we do a post-order traversal
	left := Traverse(navigator.GoToLeft(pos), navigator, eventDigest)
	rightPos := navigator.GoToRight(pos)
	if rightPos == nil {
		return NewPartialNode(pos, left)
	}
	right := Traverse(rightPos, navigator, eventDigest)
	if navigator.IsRoot(pos) {
		return NewRoot(pos, left, right)
	}
	return NewNode(pos, left, right)
}
