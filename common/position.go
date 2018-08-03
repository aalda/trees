package common

import (
	"fmt"

	"github.com/aalda/trees/util"
)

type Position struct {
	Index  []byte
	Height uint16
	// Memoizes the uint64 representation of Index.
	// We have a most 2^63-1 possible values which corresponds to a 8-size Index.
	// This should be enough for the history tree and it will truncate the Index value
	// of the hyper tree. But we don't need it.
	// Maybe we should consider using a BigInt
	indexUint64 uint64
}

func NewPosition(index []byte, height uint16) *Position {
	return &Position{
		Index:       index,
		Height:      height,
		indexUint64: util.BytesAsUint64(index),
	}
}

func NewPositionFixed(index []byte, height, numBits uint16) *Position {
	b := make([]byte, max(uint16(1), numBits/8))
	copy(b, index)
	return &Position{
		Index:       b,
		Height:      height,
		indexUint64: util.BytesAsUint64(index),
	}
}

func (p Position) String() string {
	return fmt.Sprintf("Pos(index: %d, height: %d)", p.Index, p.Height)
}

func (p Position) StringId() string {
	return fmt.Sprintf("%d|%d", p.Index, p.Height)
}

func (p Position) StringIdAsUint64() string {
	return fmt.Sprintf("%d|%d", p.IndexAsUint64(), p.Height)
}

func (p Position) Bytes() []byte {
	b := make([]byte, len(p.Index)+2) // Size of the index plus 2 bytes for the height
	copy(b, p.Index)
	copy(b[len(p.Index):], p.HeightBytes())
	return b
}

func (p Position) IndexAsUint64() uint64 {
	return p.indexUint64
}

func (p Position) HeightBytes() []byte {
	return util.Uint16AsBytes(p.Height)
}

func max(x, y uint16) uint16 {
	if x > y {
		return x
	}
	return y
}
