package common

type Navigator interface {
	IsRoot(Position) bool
	IsLeaf(Position) bool
	GoToLeft(Position) Position
	GoToRight(Position) Position
	ShouldBeCached(Position) bool
	ShouldCache(Position) bool
}

type TreeNavigator interface {
	Root() Position
	IsLeaf(Position) bool
	IsRoot(Position) bool
	GoToLeft(Position) Position
	GoToRight(Position) Position
	DescendToFirst(Position) Position
	DescendToLast(Position) Position
}
