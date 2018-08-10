package hyper

import "github.com/aalda/trees/common"

type CacheResolver interface {
	ShouldBeInCache(pos common.Position) bool
	ShouldCache(pos common.Position) bool
}

type SingleTargetedCacheResolver struct {
	numBits    uint16
	cacheLevel uint16
	targetKey  []byte
}

func NewSingleTargetedCacheResolver(numBits, cacheLevel uint16, targetKey []byte) *SingleTargetedCacheResolver {
	return &SingleTargetedCacheResolver{numBits, cacheLevel, targetKey}
}

func (r SingleTargetedCacheResolver) ShouldBeInCache(pos common.Position) bool {
	return pos.Height() != r.numBits && pos.Height() > r.cacheLevel && !r.isOnPath(pos)
}

func (r SingleTargetedCacheResolver) ShouldCache(pos common.Position) bool {
	return pos.Height() > r.cacheLevel
}

func (r SingleTargetedCacheResolver) isOnPath(pos common.Position) bool {
	bit := r.numBits - pos.Height() - 1
	return bitIsSet(r.targetKey, bit) == bitIsSet(pos.Index(), bit)
}

func bitIsSet(bits []byte, i uint16) bool {
	return bits[i/8]&(1<<uint(7-i%8)) != 0
}
