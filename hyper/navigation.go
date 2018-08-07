package hyper

import (
	"github.com/aalda/trees/common"
)

type HyperNavigator struct {
	target     common.Position
	numBits    uint16
	cacheLevel uint16
}

func NewHyperNavigator(target common.Position, numBits uint16, cacheLevel uint16) *HyperNavigator {
	return &HyperNavigator{
		target:     target,
		numBits:    numBits,
		cacheLevel: cacheLevel,
	}
}

func (n HyperNavigator) IsRoot(pos common.Position) bool {
	return pos.Height() == n.numBits
}

func (n HyperNavigator) IsLeaf(pos common.Position) bool {
	return pos.Height() == 0
}

func (n HyperNavigator) GoToLeft(pos common.Position) common.Position {
	if pos.Height() == 0 {
		return nil
	}
	return NewPosition(pos.Index(), pos.Height()-1)
}

func (n HyperNavigator) GoToRight(pos common.Position) common.Position {
	return NewPosition(n.splitBase(pos), pos.Height()-1)
}

func (n HyperNavigator) ShouldBeCached(pos common.Position) bool {
	return pos.Height() > n.cacheLevel && !n.isOnPath(pos)
}

func (n HyperNavigator) ShouldCache(pos common.Position) bool {
	return pos.Height() > n.cacheLevel
}

func (n HyperNavigator) isOnPath(pos common.Position) bool {
	bit := n.numBits - pos.Height()
	return bitGet(n.target.Index(), bit) == bitGet(pos.Index(), bit)
}

func (n HyperNavigator) splitBase(pos common.Position) []byte {
	splitBit := n.numBits - pos.Height()
	split := make([]byte, n.numBits/8)
	copy(split, pos.Index())
	if splitBit < n.numBits {
		bitSet(split, splitBit)
	}
	return split
}

func bitIsSet(bits []byte, i uint16) bool { return bits[i/8]&(1<<uint(7-i%8)) != 0 }
func bitGet(bits []byte, i uint16) byte   { return bits[i/8] & (1 << uint(7-i%8)) }
func bitSet(bits []byte, i uint16)        { bits[i/8] |= 1 << uint(7-i%8) }
func bitUnset(bits []byte, i uint16)      { bits[i/8] &= 0 << uint(7-i%8) }
