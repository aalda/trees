package history

import (
	"github.com/aalda/trees/common"
)

type CacheResolver interface {
	ShouldBeInCache(pos common.Position) bool
	ShouldCache(pos common.Position) bool
}

type SingleTargetedCacheResolver struct {
	version uint64
}

func NewSingleTargetedCacheResolver(version uint64) *SingleTargetedCacheResolver {
	return &SingleTargetedCacheResolver{version}
}

func (r SingleTargetedCacheResolver) ShouldBeInCache(pos common.Position) bool {
	return r.version > pos.IndexAsUint64()+1<<pos.Height()-1
}

func (r SingleTargetedCacheResolver) ShouldCache(pos common.Position) bool {
	return r.version >= pos.IndexAsUint64()+1<<pos.Height()-1
}

type DoubleTargetedCacheResolver struct {
	start, end uint64
}

func NewDoubleTargetedCacheResolver(start, end uint64) *DoubleTargetedCacheResolver {
	return &DoubleTargetedCacheResolver{start, end}
}

func (r DoubleTargetedCacheResolver) ShouldBeInCache(pos common.Position) bool {
	if pos.Height() == 0 && pos.IndexAsUint64() == r.start { // TODO THIS SHOULD BE TRUE for inc proofs but not for membership
		return false
	}
	threshold := pos.IndexAsUint64() + 1<<pos.Height() - 1
	if r.start > threshold && r.end > threshold {
		return true
	}

	lastDescendantIndex := pos.IndexAsUint64() + 1<<pos.Height() - 1
	return pos.IndexAsUint64() > r.start && lastDescendantIndex <= r.end
}

func (r DoubleTargetedCacheResolver) ShouldCache(pos common.Position) bool {
	return r.end >= pos.IndexAsUint64()+1<<pos.Height()-1
}

type IncrementalCacheResolver struct {
	start, end uint64
}

func NewIncrementalCacheResolver(start, end uint64) *IncrementalCacheResolver {
	return &IncrementalCacheResolver{start, end}
}

func (r IncrementalCacheResolver) ShouldBeInCache(pos common.Position) bool {
	if pos.Height() == 0 && pos.IndexAsUint64() == r.start {
		return true
	}
	threshold := pos.IndexAsUint64() + 1<<pos.Height() - 1
	if r.start > threshold && r.end > threshold {
		return true
	}

	lastDescendantIndex := pos.IndexAsUint64() + 1<<pos.Height() - 1
	return pos.IndexAsUint64() > r.start && lastDescendantIndex <= r.end
}

func (r IncrementalCacheResolver) ShouldCache(pos common.Position) bool {
	return r.end >= pos.IndexAsUint64()+1<<pos.Height()-1
}

type IncrementalVerifyCacheResolver struct {
	start, end uint64
}

func NewIncrementalVerifyCacheResolver(start, end uint64) *IncrementalVerifyCacheResolver {
	return &IncrementalVerifyCacheResolver{start, end}
}

func (r IncrementalVerifyCacheResolver) ShouldBeInCache(pos common.Position) bool {
	if pos.Height() == 0 { // changes this
		return true
	}
	threshold := pos.IndexAsUint64() + 1<<pos.Height() - 1
	if r.start > threshold && r.end > threshold {
		return true
	}

	lastDescendantIndex := pos.IndexAsUint64() + 1<<pos.Height() - 1
	return pos.IndexAsUint64() > r.start && lastDescendantIndex < r.end // changes this
}

func (r IncrementalVerifyCacheResolver) ShouldCache(pos common.Position) bool {
	return r.end >= pos.IndexAsUint64()+1<<pos.Height()-1
}
