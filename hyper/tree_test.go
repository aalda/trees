package hyper

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/logging"
	"github.com/aalda/trees/storage/bplus"
	"github.com/bbva/qed/testutils/rand"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {

	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

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
	simpleCache := common.NewSimpleCache(10)
	tree := NewHyperTree(new(common.XorHasher), store, simpleCache, 2)

	for i, c := range testCases {
		index := uint64(i)
		commitment := tree.Add(c.eventDigest, index)
		require.Equalf(t, c.expectedRootHash, commitment.Digest, "Incorrect root hash for index %d", i)
	}
}

func BenchmarkAdd(b *testing.B) {

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	store, closeF := openBadgerStore("/var/tmp/hyper_tree_test.db")
	defer closeF()

	hasher := common.NewSha256Hasher()
	simpleCache := common.NewSimpleCache(1 << 25)
	tree := NewHyperTree(common.NewSha256Hasher(), store, simpleCache, hasher.Len()-25)

	b.ResetTimer()
	b.N = 100000
	for i := 0; i < b.N; i++ {
		key := hasher.Do(rand.Bytes(32))
		tree.Add(key, uint64(i))
	}
}
