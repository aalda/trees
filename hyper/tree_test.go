package hyper

import (
	"testing"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/storage"
	"github.com/aalda/trees/storage/bplus"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {

	testCases := []struct {
		eventDigest      common.Digest
		expectedRootHash common.Digest
	}{
		{common.Digest{0x0}, common.Digest{0x0}},
		{common.Digest{0x1}, common.Digest{0x1}},
		{common.Digest{0x2}, common.Digest{0x3}},
		{common.Digest{0x3}, common.Digest{0x0}},
		{common.Digest{0x4}, common.Digest{0x4}},
		{common.Digest{0x5}, common.Digest{0x1}},
		{common.Digest{0x6}, common.Digest{0x7}},
		{common.Digest{0x7}, common.Digest{0x0}},
		{common.Digest{0x8}, common.Digest{0x8}},
		{common.Digest{0x9}, common.Digest{0x1}},
	}

	store := bplus.NewBPlusTreeStorage()
	cache := common.NewPassThroughCache(storage.HyperCachePrefix, store)
	twoLevel := common.NewTwoLevelCache(cache)
	fallback := common.NewFallbackCache([]byte("blah"), 8, new(common.XorHasher), twoLevel)
	tree := NewHyperTree(new(common.XorHasher), store, fallback, 4)

	for i, c := range testCases {
		index := uint64(i)
		commitment := tree.Add(c.eventDigest, index)
		require.Equalf(t, c.expectedRootHash, commitment.Digest, "Incorrect root hash for index %d", i)
	}
}
