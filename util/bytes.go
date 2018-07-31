package util

import "encoding/binary"

func Uint64AsBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

func BytesAsUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}
