package v1

import (
	"bytes"
	"testing"
)

var (
	benchmarkGenerateKeyPassword = []byte("correct-horse-battery-staple")
	benchmarkGenerateKeyNonce    = bytes.Repeat([]byte{1}, NonceLength)
	benchmarkGenerateKeyResult   []byte
)

func BenchmarkGenerateKey(b *testing.B) {
	var result []byte

	for i := 0; i < b.N; i++ {
		result = generateKey(
			benchmarkGenerateKeyPassword,
			benchmarkGenerateKeyNonce)
	}

	benchmarkGenerateKeyResult = result
}
