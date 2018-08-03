package common

import (
	"crypto/sha256"
	"hash"
)

type Digest []byte

type Hasher interface {
	Do(...[]byte) []byte
	Len() uint16
}

type XorHasher struct{}

func (x XorHasher) Do(data ...[]byte) []byte {
	var result byte
	for _, elem := range data {
		var sum byte
		for _, b := range elem {
			sum = sum ^ b
		}
		result = result ^ sum
	}
	return []byte{result}
}
func (s XorHasher) Len() uint16 { return uint16(8) }

type Sha256Hasher struct {
	underlying hash.Hash
}

func NewSha256Hasher() *Sha256Hasher {
	return &Sha256Hasher{underlying: sha256.New()}
}

func (s *Sha256Hasher) Do(data ...[]byte) []byte {
	s.underlying.Reset()
	for i := 0; i < len(data); i++ {
		s.underlying.Write(data[i])
	}
	return s.underlying.Sum(nil)[:]
}

func (s Sha256Hasher) Len() uint16 { return uint16(256) }
