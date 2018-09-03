package history

import (
	"math"

	"github.com/aalda/trees/common"
)

type HistoryTreeNavigator struct {
	version uint64
	depth   uint16
}

func NewHistoryTreeNavigator(version uint64) *HistoryTreeNavigator {
	depth := uint16(uint64(math.Ceil(math.Log2(float64(version + 1)))))
	return &HistoryTreeNavigator{version, depth}
}

func (n HistoryTreeNavigator) Root() common.Position {
	return NewPosition(0, n.depth)
}

func (n HistoryTreeNavigator) IsLeaf(pos common.Position) bool {
	return pos.Height() == 0
}

func (n HistoryTreeNavigator) IsRoot(pos common.Position) bool {
	return pos.Height() == n.depth
}

func (n HistoryTreeNavigator) GoToLeft(pos common.Position) common.Position {
	if pos.Height() == 0 {
		return nil
	}
	return NewPosition(pos.IndexAsUint64(), pos.Height()-1)
}
func (n HistoryTreeNavigator) GoToRight(pos common.Position) common.Position {
	rightIndex := pos.IndexAsUint64() + 1<<(pos.Height()-1)
	if pos.Height() == 0 || rightIndex > n.version {
		return nil
	}
	return NewPosition(rightIndex, pos.Height()-1)
}

func (n HistoryTreeNavigator) DescendToFirst(pos common.Position) common.Position {
	return NewPosition(pos.IndexAsUint64(), 0)
}

func (n HistoryTreeNavigator) DescendToLast(pos common.Position) common.Position {
	lastDescendantIndex := pos.IndexAsUint64() + 1<<pos.Height() - 1
	return NewPosition(lastDescendantIndex, 0)
}
