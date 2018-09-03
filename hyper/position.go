package hyper

import (
	"encoding/binary"
	"fmt"

	"github.com/aalda/trees/util"
)

type HyperPosition struct {
	index  []byte
	height uint16
}

func NewPosition(index []byte, height uint16) *HyperPosition {
	return &HyperPosition{
		index:  index,
		height: height,
	}
}

func (p HyperPosition) Index() []byte {
	return p.index
}

func (p HyperPosition) Height() uint16 {
	return p.height
}

func (p HyperPosition) Bytes() []byte {
	return append(p.index, util.Uint16AsBytes(p.height)...)
}

func (p HyperPosition) String() string {
	return fmt.Sprintf("Pos(%x, %d)", p.index, p.height)
}

func (p HyperPosition) StringId() string {
	return fmt.Sprintf("%x|%d", p.index, p.height)
}

func (p HyperPosition) IndexAsUint64() uint64 {
	return binary.LittleEndian.Uint64(p.index)
}
