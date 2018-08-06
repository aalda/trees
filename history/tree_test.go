package history

import (
	"testing"

	"github.com/aalda/trees/common"
	"github.com/aalda/trees/storage"
	"github.com/aalda/trees/storage/bplus"
	"github.com/aalda/trees/util"

	"github.com/bbva/qed/testutils/rand"
	"github.com/stretchr/testify/assert"
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
	cache := common.NewPassThroughCache(storage.HistoryCachePrefix, store)
	tree := NewHistoryTree(new(common.XorHasher), store, cache)

	for i, c := range testCases {
		index := uint64(i)
		commitment := tree.Add(c.eventDigest, index)
		require.Equalf(t, c.expectedRootHash, commitment.Digest, "Incorrect root hash for index %d", i)
	}

}

func TestProveMembership(t *testing.T) {

	testCases := []struct {
		eventDigest common.Digest
		auditPath   common.AuditPath
	}{
		{
			common.Digest{0x0},
			common.AuditPath{},
		},
		{
			common.Digest{0x1},
			common.AuditPath{"0|0": common.Digest{0x0}},
		},
		{
			common.Digest{0x2},
			common.AuditPath{"0|1": common.Digest{0x1}},
		},
		{
			common.Digest{0x3},
			common.AuditPath{"0|1": common.Digest{0x1}, "2|0": common.Digest{0x2}},
		},
		{
			common.Digest{0x4},
			common.AuditPath{"0|2": common.Digest{0x0}},
		},
		{
			common.Digest{0x5},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|0": common.Digest{0x4}},
		},
		{
			common.Digest{0x6},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|1": common.Digest{0x1}},
		},
		{
			common.Digest{0x7},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|1": common.Digest{0x1}, "6|0": common.Digest{0x6}},
		},
		{
			common.Digest{0x8},
			common.AuditPath{"0|3": common.Digest{0x0}},
		},
		{
			common.Digest{0x9},
			common.AuditPath{"0|3": common.Digest{0x0}, "8|0": common.Digest{0x8}},
		},
	}

	store := bplus.NewBPlusTreeStorage()
	cache := common.NewPassThroughCache(storage.HistoryCachePrefix, store)
	tree := NewHistoryTree(new(common.XorHasher), store, cache)

	for i, c := range testCases {
		index := uint64(i)
		tree.Add(c.eventDigest, index)
		pf := tree.ProveMembership(index, index)
		require.Equalf(t, c.auditPath, pf.AuditPath, "Incorrect audit path for index %d", i)
	}
}

func TestProveMembershipNonConsecutive(t *testing.T) {
	store := bplus.NewBPlusTreeStorage()
	cache := common.NewPassThroughCache(storage.HistoryCachePrefix, store)
	tree := NewHistoryTree(new(common.XorHasher), store, cache)

	// add nine events
	for i := uint64(0); i < 9; i++ {
		eventDigest := util.Uint64AsBytes(i)
		tree.Add(eventDigest, i)
	}

	// query for membership with event 0 and version 8
	proof := tree.ProveMembership(0, 8)
	expectedAuditPath := common.AuditPath{"1|0": common.Digest{0x1}, "2|1": common.Digest{0x1}, "4|2": common.Digest{0x0}, "8|0": common.Digest{0x8}}
	assert.Equal(t, expectedAuditPath, proof.AuditPath, "Invalid audit path")
}

func TestVerify(t *testing.T) {

	testCases := []struct {
		eventDigest    common.Digest
		expectedDigest common.Digest
		auditPath      common.AuditPath
	}{
		{
			common.Digest{0x0},
			common.Digest{0x0},
			common.AuditPath{},
		},
		{
			common.Digest{0x1},
			common.Digest{0x1},
			common.AuditPath{"0|0": common.Digest{0x0}},
		},
		{
			common.Digest{0x2},
			common.Digest{0x3},
			common.AuditPath{"0|1": common.Digest{0x1}},
		},
		{
			common.Digest{0x3},
			common.Digest{0x0},
			common.AuditPath{"0|1": common.Digest{0x1}, "2|0": common.Digest{0x2}},
		},
		{
			common.Digest{0x4},
			common.Digest{0x4},
			common.AuditPath{"0|2": common.Digest{0x0}},
		},
		{
			common.Digest{0x5},
			common.Digest{0x1},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|0": common.Digest{0x4}},
		},
		{
			common.Digest{0x6},
			common.Digest{0x7},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|1": common.Digest{0x1}},
		},
		{
			common.Digest{0x7},
			common.Digest{0x0},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|1": common.Digest{0x1}, "6|0": common.Digest{0x6}},
		},
		{
			common.Digest{0x8},
			common.Digest{0x8},
			common.AuditPath{"0|3": common.Digest{0x0}},
		},
		{
			common.Digest{0x9},
			common.Digest{0x1},
			common.AuditPath{"0|3": common.Digest{0x0}, "8|0": common.Digest{0x8}},
		},
	}

	store := bplus.NewBPlusTreeStorage()
	cache := common.NewPassThroughCache(storage.HistoryCachePrefix, store)
	tree := NewHistoryTree(new(common.XorHasher), store, cache)

	for i, c := range testCases {
		index := uint64(i)
		proof := NewMembershipProof(c.auditPath)
		correct := tree.VerifyMembership(proof, index, c.eventDigest, c.expectedDigest)
		require.Truef(t, correct, "Event with index %d should be a member", index)
	}

}

func TestProveConsistency(t *testing.T) {

	testCases := []struct {
		eventDigest common.Digest
		auditPath   common.AuditPath
	}{
		{
			common.Digest{0x0},
			common.AuditPath{"0|0": common.Digest{0x0}},
		},
		{
			common.Digest{0x1},
			common.AuditPath{"0|0": common.Digest{0x0}, "1|0": common.Digest{0x1}},
		},
		{
			common.Digest{0x2},
			common.AuditPath{"0|0": common.Digest{0x0}, "1|0": common.Digest{0x1}, "2|0": common.Digest{0x2}},
		},
		{
			common.Digest{0x3},
			common.AuditPath{"0|1": common.Digest{0x1}, "2|0": common.Digest{0x2}, "3|0": common.Digest{0x3}},
		},
		{
			common.Digest{0x4},
			common.AuditPath{"0|1": common.Digest{0x1}, "2|0": common.Digest{0x2}, "3|0": common.Digest{0x3}, "4|0": common.Digest{0x4}},
		},
		{
			common.Digest{0x5},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|0": common.Digest{0x4}, "5|0": common.Digest{0x5}},
		},
		{
			common.Digest{0x6},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|0": common.Digest{0x4}, "5|0": common.Digest{0x5}, "6|0": common.Digest{0x6}},
		},
		{
			common.Digest{0x7},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|1": common.Digest{0x1}, "6|0": common.Digest{0x6}, "7|0": common.Digest{0x7}},
		},
		{
			common.Digest{0x8},
			common.AuditPath{"0|2": common.Digest{0x0}, "4|1": common.Digest{0x1}, "6|0": common.Digest{0x6}, "7|0": common.Digest{0x7}, "8|0": common.Digest{0x8}},
		},
		{
			common.Digest{0x9},
			common.AuditPath{"0|3": common.Digest{0x0}, "8|0": common.Digest{0x8}, "9|0": common.Digest{0x9}},
		},
	}

	store := bplus.NewBPlusTreeStorage()
	cache := common.NewPassThroughCache(storage.HistoryCachePrefix, store)
	tree := NewHistoryTree(new(common.XorHasher), store, cache)

	for i, c := range testCases {
		index := uint64(i)
		tree.Add(c.eventDigest, index)

		proof := tree.ProveConsistency(uint64(max(0, i-1)), index)
		require.Equal(t, c.auditPath, proof.AuditPath, "Invalid audit path in test case: %d", i)
	}

}

func TestProveConsistencyNonConsecutive(t *testing.T) {
	store := bplus.NewBPlusTreeStorage()
	cache := common.NewPassThroughCache(storage.HistoryCachePrefix, store)
	tree := NewHistoryTree(new(common.XorHasher), store, cache)

	// add nine events
	for i := uint64(0); i < 9; i++ {
		eventDigest := util.Uint64AsBytes(i)
		tree.Add(common.Digest(eventDigest), i)
	}

	// query for consistency with event 2 and version 8
	proof := tree.ProveConsistency(uint64(2), uint64(8))
	expectedAuditPath := common.AuditPath{
		"0|1": common.Digest{0x1}, "2|0": common.Digest{0x2}, "3|0": common.Digest{0x3},
		"4|2": common.Digest{0x0}, "8|0": common.Digest{0x8},
	}
	require.Equal(t, expectedAuditPath, proof.AuditPath, "Invalid audit path")
}

func TestProveConsistencySameVersions(t *testing.T) {
	store := bplus.NewBPlusTreeStorage()
	cache := common.NewPassThroughCache(storage.HistoryCachePrefix, store)
	tree := NewHistoryTree(new(common.XorHasher), store, cache)

	// add nine events
	for i := uint64(0); i < 9; i++ {
		eventDigest := util.Uint64AsBytes(i)
		tree.Add(common.Digest(eventDigest), i)
	}

	// query for consistency with event 8 and version 8
	proof := tree.ProveConsistency(uint64(8), uint64(8))
	expectedAuditPath := common.AuditPath{"0|3": common.Digest{0x0}, "8|0": common.Digest{0x8}}
	require.Equal(t, expectedAuditPath, proof.AuditPath, "Invalid audit path")
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func BenchmarkAdd(b *testing.B) {
	store := bplus.NewBPlusTreeStorage()
	cache := common.NewPassThroughCache(storage.HistoryCachePrefix, store)
	tree := NewHistoryTree(common.NewSha256Hasher(), store, cache)
	b.N = 100000
	for i := uint64(0); i < uint64(b.N); i++ {
		key := rand.Bytes(64)
		tree.Add(key, i)
	}
}
