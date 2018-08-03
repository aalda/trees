package storage

import (
	"bytes"
	"sort"
)

const (
	VersionPrefix      = byte(0x0)
	IndexPrefix        = byte(0x1)
	HyperCachePrefix   = byte(0x2)
	HistoryCachePrefix = byte(0x3)
)

type Mutation struct {
	Prefix     byte
	Key, Value []byte
}

func NewMutation(prefix byte, key, value []byte) *Mutation {
	return &Mutation{prefix, key, value}
}

type KVPair struct {
	Key, Value []byte
}

type KVRange []KVPair

func (r KVRange) InsertSorted(p KVPair) KVRange {
	index := sort.Search(len(r), func(i int) bool {
		return bytes.Compare(r[i].Key, p.Key) > 0
	})
	r = append(r, p)
	copy(r[index+1:], r[index:])
	r[index] = p
	return r
}

func (r KVRange) Split(key []byte) (left, right KVRange) {
	// the smallest index i where r[i] >= index
	index := sort.Search(len(r), func(i int) bool {
		return bytes.Compare(r[i].Key, key) >= 0
	})
	return r[:index], r[index:]
}

func (r KVRange) Get(key []byte) KVPair {
	index := sort.Search(len(r), func(i int) bool {
		return bytes.Compare(r[i].Key, key) >= 0
	})
	if index < len(r) && bytes.Equal(r[index].Key, key) {
		return r[index]
	} else {
		panic("This should never happen")
	}
}

type Store interface {
	Mutate(mutations []Mutation) error
	GetRange(prefix byte, start, end []byte) (KVRange, error)
	Get(prefix byte, key []byte) (*KVPair, error)
}

type Cache interface {
	Get(key []byte) (*KVPair, error)
}
