package hyper

import (
	"testing"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/log"
	"github.com/aalda/trees/storage/bplus"
	"github.com/aalda/trees/util"
	"github.com/bbva/qed/testutils/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {

	log.SetLogger("TestAdd", log.DEBUG)

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

func TestProveMembership(t *testing.T) {

	log.SetLogger("TestProveMembership", log.DEBUG)

	hasher := new(common.XorHasher)
	digest := hasher.Do(common.Digest{0x0})
	index := uint64(0)

	store := bplus.NewBPlusTreeStorage()
	simpleCache := common.NewSimpleCache(10)
	tree := NewHyperTree(new(common.XorHasher), store, simpleCache, 2)

	rh := tree.Add(digest, index)
	assert.Equal(t, rh.Digest, common.Digest{0x0}, "Incorrect root hash")

	_, pf, err := tree.Get(digest)
	assert.Nil(t, err, "Error adding to the tree: %v", err)

	ap := common.AuditPath{
		"10|4": common.Digest{0x0},
		"04|2": common.Digest{0x0},
		"80|7": common.Digest{0x0},
		"40|6": common.Digest{0x0},
		"20|5": common.Digest{0x0},
		"08|3": common.Digest{0x0},
		"02|1": common.Digest{0x0},
		"01|0": common.Digest{0x0},
	}
	assert.Equal(t, ap, pf.AuditPath, "Incorrect audit path")

}

func TestProveMembershipConsecutive(t *testing.T) {

	log.SetLogger("TestProveMembership", log.DEBUG)

	hasher := new(common.XorHasher)
	digest := hasher.Do(common.Digest{0x0})
	index := uint64(0)

	store := bplus.NewBPlusTreeStorage()
	simpleCache := common.NewSimpleCache(10)
	tree := NewHyperTree(new(common.XorHasher), store, simpleCache, 2)

	rh := tree.Add(digest, index)
	assert.Equal(t, rh.Digest, common.Digest{0x0}, "Incorrect root hash")

	tree.Add(hasher.Do(common.Digest{0x1}), uint64(1))
	tree.Add(hasher.Do(common.Digest{0x2}), uint64(2))

	_, pf, err := tree.Get(digest)
	assert.Nil(t, err, "Error adding to the tree: %v", err)

	ap := common.AuditPath{
		"10|4": common.Digest{0x0},
		"04|2": common.Digest{0x0},
		"80|7": common.Digest{0x0},
		"40|6": common.Digest{0x0},
		"20|5": common.Digest{0x0},
		"08|3": common.Digest{0x0},
		"02|1": common.Digest{0x2},
		"01|0": common.Digest{0x1},
	}
	assert.Equal(t, ap, pf.AuditPath, "Incorrect audit path")

}

func TestAddAndVerifyXor(t *testing.T) {

	log.SetLogger("TestAddAndVerifyXor", log.DEBUG)

	hasher := new(common.XorHasher)
	store := bplus.NewBPlusTreeStorage()
	simpleCache := common.NewSimpleCache(10)
	tree := NewHyperTree(new(common.XorHasher), store, simpleCache, 2)

	key := hasher.Do(common.Digest("a test event"))
	value := uint64(0)

	commitment := tree.Add(key, value)

	actualValue, proof, err := tree.Get(key)
	assert.Nil(t, err, "Error must be nil")

	assert.Equal(t, util.Uint64AsBytes(value), actualValue, "Incorrect actual value")

	correct := tree.VerifyMembership(proof, value, key, commitment.Digest)

	if !correct {
		t.Errorf("Key %x should be a member", key)
	}
}

func TestAddAndVerifySha256(t *testing.T) {

	log.SetLogger("TestAddAndVerifySha256", log.DEBUG)

	hasher := common.NewSha256Hasher()
	store := bplus.NewBPlusTreeStorage()
	simpleCache := common.NewSimpleCache(10)
	tree := NewHyperTree(common.NewSha256Hasher(), store, simpleCache, 2)

	key := hasher.Do(common.Digest("a test event"))
	value := uint64(0)

	commitment := tree.Add(key, value)

	actualValue, proof, err := tree.Get(key)
	assert.Nil(t, err, "Error must be nil")

	assert.Equal(t, util.Uint64AsBytes(value), actualValue, "Incorrect actual value")

	correct := tree.VerifyMembership(proof, value, key, commitment.Digest)

	if !correct {
		t.Errorf("Key %x should be a member", key)
	}
}

func BenchmarkAdd(b *testing.B) {

	log.SetLogger("BenchmarkAdd", log.SILENT)

	store, closeF := openBadgerStore("/var/tmp/hyper_tree_test.db")
	defer closeF()

	hasher := common.NewSha256Hasher()
	simpleCache := common.NewSimpleCache(0)
	tree := NewHyperTree(common.NewSha256Hasher(), store, simpleCache, hasher.Len()-25)

	b.ResetTimer()
	b.N = 100000
	for i := 0; i < b.N; i++ {
		key := hasher.Do(rand.Bytes(32))
		tree.Add(key, uint64(i))
	}
}
