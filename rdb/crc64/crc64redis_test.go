package crc64_test

import (
	"testing"
	"github.com/mantyr/redis-replica/rdb/crc64"
)

func TestRedisCRC64(t *testing.T) {
	hash := crc64.CRC64Update(0, []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'})
	if hash != 0xe9c6d914c4b8d9ca {
		t.Errorf("crc64 doesn't match: crc64(\"123456789\") = %#v != 0xe9c6d914c4b8d9ca", hash)
	}
}

func BenchmarkRedisCRC64(b *testing.B) {
	data := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
	var crc uint64
	for i := 0; i < b.N; i++ {
		crc64.CRC64Update(crc, data)
	}
}
