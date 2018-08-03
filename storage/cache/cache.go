package cache

import "github.com/aalda/trees/storage"

type PassThroughCache struct {
	prefix byte
	store  storage.Store
}

func NewPassThroughCache(prefix byte, store storage.Store) *PassThroughCache {
	return &PassThroughCache{prefix, store}
}

func (c PassThroughCache) Get(key []byte) (*storage.KVPair, error) {
	return c.store.Get(c.prefix, key)
}

type TwoLevelCache struct {
	prefix byte
	store  *storage.Store
	//cached map[[36]byte]Digest
}

type FallbackCache struct{}
