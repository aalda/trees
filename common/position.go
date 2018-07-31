package common

import (
	"fmt"

	"github.com/aalda/trees/util"
)

type Position struct {
	Index  uint64 // THIS SHOULD BE A []byte
	Height uint64
}

func NewPosition(index, height uint64) *Position {
	return &Position{Index: index, Height: height}
}

func (p Position) String() string {
	return fmt.Sprintf("Pos(index: %d, height: %d)", p.Index, p.Height)
}

func (p Position) StringId() string {
	return fmt.Sprintf("%d|%d", p.Index, p.Height)
}

func (p Position) Bytes() []byte {
	b := make([]byte, 16) // idLen is the size of layer and height, which is 16 bytes
	copy(b, p.IndexBytes())
	copy(b[len(p.IndexBytes()):], p.HeightBytes())
	return b
}

func (p Position) IndexBytes() []byte {
	return util.Uint64AsBytes(p.Index)
}

func (p Position) HeightBytes() []byte {
	return util.Uint64AsBytes(p.Height)
}
