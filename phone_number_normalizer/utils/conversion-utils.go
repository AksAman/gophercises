package utils

import (
	"encoding/binary"
)

func IntToBytes(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func BytesToInt(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func BytesToString(b []byte) string {
	return string(b)
}

func StringToBytes(s string) []byte {
	return []byte(s)
}
