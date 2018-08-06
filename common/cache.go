package common

import "github.com/aalda/trees/storage"

type Cache interface {
	Get(pos Position) (Digest, bool)
}

type PassThroughCache struct {
	prefix byte
	store  storage.Store
}

func NewPassThroughCache(prefix byte, store storage.Store) *PassThroughCache {
	return &PassThroughCache{prefix, store}
}

func (c PassThroughCache) Get(pos Position) (Digest, bool) {
	pair, err := c.store.Get(c.prefix, pos.Bytes())
	if err != nil {
		return nil, false
	}
	if len(pair.Key) > 0 { // TODO FIX THIS
		return pair.Value, true
	}
	return nil, false
}

type TwoLevelCache struct {
	prefix byte
	store  *storage.Store
	//cached map[[36]byte]Digest
}

type FallbackCache struct {
	decorated     Cache
	defaultHashes []Digest
}

func NewFallbackCache(id []byte, height uint16, hasher Hasher, decorated Cache) *FallbackCache {
	hashes := make([]Digest, height)
	hashes[0] = hasher.Do(id, []byte{0x0})
	for i := uint16(1); i < height; i++ {
		hashes[i] = hasher.Do(hashes[i-1], hashes[i-1])
	}
	return &FallbackCache{
		decorated:     decorated,
		defaultHashes: hashes,
	}
}

func (c FallbackCache) Get(pos Position) (Digest, bool) {
	digest, ok := c.decorated.Get(pos)
	if ok {
		return digest, ok
	}
	return c.defaultHashes[pos.Height()], false
}
