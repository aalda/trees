package common

import (
	"fmt"
)

type Visitor interface {
	VisitRoot(pos *Position, leftResult, rightResult interface{}) interface{}
	VisitNode(pos *Position, leftResult, rightResult interface{}) interface{}
	VisitPartialNode(pos *Position, leftResult interface{}) interface{}
	VisitLeaf(pos *Position, value []byte) interface{}
	VisitCached(pos *Position) interface{}
}

type Visitable interface {
	Accept(visitor Visitor) interface{}
	String() string
}

type Root struct {
	pos         *Position
	left, right Visitable
}

type Node struct {
	pos         *Position
	left, right Visitable
}

type PartialNode struct {
	pos  *Position
	left Visitable
}

type Leaf struct {
	pos         *Position
	eventDigest Digest
}

type Cached struct {
	pos *Position
}

func (r Root) Accept(visitor Visitor) interface{} {
	leftResult := r.left.Accept(visitor)
	rightResult := r.right.Accept(visitor)
	return visitor.VisitRoot(r.pos, leftResult, rightResult)
}

func (r Root) String() string {
	return fmt.Sprintf("Root(%d, %d)[ %v | %v ]", r.pos.Index, r.pos.Height, r.left, r.right)
}

func (n Node) Accept(visitor Visitor) interface{} {
	leftResult := n.left.Accept(visitor)
	rightResult := n.right.Accept(visitor)
	return visitor.VisitNode(n.pos, leftResult, rightResult)
}

func (n Node) String() string {
	return fmt.Sprintf("Node(%d, %d)[ %v | %v ]", n.pos.Index, n.pos.Height, n.left, n.right)
}

func (p PartialNode) Accept(visitor Visitor) interface{} {
	leftResult := p.left.Accept(visitor)
	return visitor.VisitPartialNode(p.pos, leftResult)
}

func (p PartialNode) String() string {
	return fmt.Sprintf("PartialNode(%d, %d)[ %v ]", p.pos.Index, p.pos.Height, p.left)
}

func (l Leaf) Accept(visitor Visitor) interface{} {
	return visitor.VisitLeaf(l.pos, l.eventDigest)
}

func (l Leaf) String() string {
	return fmt.Sprintf("Leaf(%d, %d)", l.pos.Index, l.pos.Height)
}

func (c Cached) Accept(visitor Visitor) interface{} {
	return visitor.VisitCached(c.pos)
}

func (c Cached) String() string {
	return fmt.Sprintf("Cached(%d, %d)", c.pos.Index, c.pos.Height)
}

func Traverse(pos *Position, navigator Navigator, eventDigest Digest) Visitable {
	if navigator.ShouldBeCached(pos) {
		return &Cached{pos}
	}
	if navigator.IsLeaf(pos) {
		return &Leaf{pos, eventDigest}
	} else {
		// we do a post-order traversal
		left := Traverse(navigator.GoToLeft(pos), navigator, eventDigest)
		rightPos := navigator.GoToRight(pos)
		if rightPos == nil {
			return &PartialNode{pos, left}
		} else {
			right := Traverse(rightPos, navigator, eventDigest)
			if navigator.IsRoot(pos) {
				return &Root{pos, left, right}
			}
			return &Node{pos, left, right}
		}
	}
}