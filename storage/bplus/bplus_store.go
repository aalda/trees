package bplus

import (
	"bytes"

	"github.com/aalda/trees/storage"
	"github.com/google/btree"
)

type BPlusTreeStore struct {
	db *btree.BTree
}

func NewBPlusTreeStorage() *BPlusTreeStore {
	return &BPlusTreeStore{btree.New(2)}
}

type KVItem struct {
	Key, Value []byte
}

func (p KVItem) Less(b btree.Item) bool {
	return bytes.Compare(p.Key, b.(KVItem).Key) < 0
}

func (s *BPlusTreeStore) Mutate(mutations []storage.Mutation) error {
	for _, m := range mutations {
		key := append([]byte{m.Prefix}, m.Key...)
		s.db.ReplaceOrInsert(KVItem{key, m.Value})
	}
	return nil
}

func (s *BPlusTreeStore) GetRange(prefix byte, start, end []byte) (storage.KVRange, error) {
	result := make(storage.KVRange, 0)
	startKey := append([]byte{prefix}, start...)
	endKey := append([]byte{prefix}, end...)
	s.db.AscendGreaterOrEqual(KVItem{startKey, nil}, func(i btree.Item) bool {
		key := i.(KVItem).Key
		if bytes.Compare(key, endKey) > 0 {
			return false
		}
		result = append(result, storage.KVPair{key[1:], i.(KVItem).Value})
		return true
	})
	return result, nil
}

func (s *BPlusTreeStore) Get(prefix byte, key []byte) (*storage.KVPair, error) {
	result := new(storage.KVPair)
	result.Key = key
	k := append([]byte{prefix}, key...)
	item := s.db.Get(KVItem{k, nil})
	if item != nil {
		result.Value = item.(KVItem).Value
	}
	return result, nil
}

func (s BPlusTreeStore) Close() error {
	s.db.Clear(false)
	return nil
}
