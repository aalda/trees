package common

type Store interface {
	Add(key Position, value []byte)
	Get(key Position) []byte
}

type InMemoryStore struct {
	elems map[Position]Digest
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{make(map[Position]Digest)}
}

func (s *InMemoryStore) Add(key Position, value []byte) {
	s.elems[key] = value
}

func (s *InMemoryStore) Get(key Position) []byte {
	e, ok := s.elems[key]
	if !ok {
		e = make([]byte, 0)
	}
	return e
}
