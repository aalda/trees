package common

type Cache interface {
	Get(pos Position) (Digest, bool)
}

type ModifiableCache interface {
	Put(pos Position, value Digest)
	Cache
}

type PassThroughCache struct {
	prefix byte
	store  Store
}

func NewPassThroughCache(prefix byte, store Store) *PassThroughCache {
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

const keySize = 34

type SimpleCache struct {
	cached map[[keySize]byte]Digest
}

func NewSimpleCache(size uint64) *SimpleCache {
	return &SimpleCache{make(map[[keySize]byte]Digest, size)}
}

func (c SimpleCache) Get(pos Position) (Digest, bool) {
	var key [keySize]byte
	copy(key[:], pos.Bytes())
	digest, ok := c.cached[key]
	return digest, ok
}

func (c *SimpleCache) Put(pos Position, value Digest) {
	var key [keySize]byte
	copy(key[:], pos.Bytes())
	c.cached[key] = value
}

type TwoLevelCache struct {
	decorated Cache
	cached    map[[keySize]byte]Digest
}

func NewTwoLevelCache(size uint64, decorated Cache) *TwoLevelCache {
	return &TwoLevelCache{
		decorated: decorated,
		cached:    make(map[[keySize]byte]Digest, size),
	}
}

func (c TwoLevelCache) Get(pos Position) (Digest, bool) {
	var key [keySize]byte
	copy(key[:], pos.Bytes())

	digest, ok := c.cached[key]
	if !ok {
		digest, ok = c.decorated.Get(pos)
		if ok {
			c.cached[key] = digest
		}
	}
	return digest, ok
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
	return c.defaultHashes[pos.Height()], true
}
