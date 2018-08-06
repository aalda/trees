package hyper

import (
	"fmt"
	"testing"

	"github.com/aalda/trees/storage"
	"github.com/aalda/trees/storage/bplus"
)

func TestHyperTraverser(t *testing.T) {

	// we use a empty fake store
	store := bplus.NewBPlusTreeStorage()

	// we want to insert one new leaf
	targetPos := NewPosition([]byte{0x1}, 0)
	leaves := storage.NewKVRange()
	leaves = leaves.InsertSorted(storage.KVPair{targetPos.Index(), []byte("blah blah")})
	prev := NewPosition([]byte{0x0}, 0)
	leaves = leaves.InsertSorted(storage.KVPair{prev.Index(), []byte("blah blah")})

	// we use a 8 bits hash function and a cache with 2 levels
	numBits := uint16(8)
	cacheLevel := uint16(6)
	traverser := NewHyperTraverser(numBits, cacheLevel, leaves, store)
	navigator := NewHyperNavigator(targetPos, numBits, cacheLevel)

	root := traverser.Traverse(NewPosition([]byte{0x0}, numBits), navigator)
	fmt.Println(root)
}
