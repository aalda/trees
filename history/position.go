package history

import (
	"fmt"

	"github.com/aalda/trees/util"
)

type HistoryPosition struct {
	index  uint64
	height uint16
}

func NewPosition(index uint64, height uint16) *HistoryPosition {
	return &HistoryPosition{
		index:  index,
		height: height,
	}
}

func (p HistoryPosition) Index() []byte {
	return util.Uint64AsBytes(p.index)
}

func (p HistoryPosition) Height() uint16 {
	return p.height
}

func (p HistoryPosition) IndexAsUint64() uint64 {
	return p.index
}

func (p HistoryPosition) Bytes() []byte {
	return append(p.Index(), util.Uint16AsBytes(p.height)...)
}

func (p HistoryPosition) String() string {
	return fmt.Sprintf("Pos(%d, %d)", p.IndexAsUint64(), p.height)
}

func (p HistoryPosition) StringId() string {
	return fmt.Sprintf("%d|%d", p.IndexAsUint64(), p.height)
}
