package common

type Navigator interface {
	IsRoot(Position) bool
	IsLeaf(Position) bool
	GoToLeft(Position) Position
	GoToRight(Position) Position
	ShouldBeCached(Position) bool
	ShouldCache(Position) bool
}
