package binutils_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"

	"github.com/9seconds/chore/internal/binutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMixString(t *testing.T) {
	testTable := map[string]string{
		"":   "010000000000000000",
		"aa": "0102000000000000006161",
	}

	for testValue, expected := range testTable {
		testValue := testValue
		expected := expected

		t.Run(testValue, func(t *testing.T) {
			buf := &bytes.Buffer{}

			assert.NoError(t, binutils.MixString(buf, testValue))
			assert.Equal(t, expected, hex.EncodeToString(buf.Bytes()))
		})
	}
}

func TestMixStringSlice(t *testing.T) {
	testTable := map[string]string{
		"":       "020100000000000000010000000000000000",
		"aa":     "0201000000000000000102000000000000006161",
		"aa,bb":  "02020000000000000001020000000000000061610102000000000000006262",
		"aa,b,b": "02030000000000000001020000000000000061610101000000000000006201010000000000000062",
	}

	for testValue, expected := range testTable {
		testValue := testValue
		expected := expected

		t.Run(testValue, func(t *testing.T) {
			buf := &bytes.Buffer{}
			values := strings.Split(testValue, ",")

			assert.NoError(t, binutils.MixStringSlice(buf, values))
			assert.Equal(t, expected, hex.EncodeToString(buf.Bytes()))
		})
	}
}

func TestMixLength(t *testing.T) {
	testTable := map[int]string{
		0:      "0000000000000000",
		1:      "0100000000000000",
		100002: "a286010000000000",
	}

	for testValue, expected := range testTable {
		testValue := testValue
		expected := expected

		t.Run(strconv.Itoa(testValue), func(t *testing.T) {
			buf := &bytes.Buffer{}

			assert.NoError(t, binutils.MixLength(buf, testValue))
			assert.Equal(t, expected, hex.EncodeToString(buf.Bytes()))
		})
	}
}

func TestSortedMapKeys(t *testing.T) {
	data := make(map[string]bool)

	for i := 0; i < 5000; i++ {
		str := make([]byte, 32)
		_, err := rand.Read(str)

		require.NoError(t, err)

		data[base64.StdEncoding.EncodeToString(str)] = true
	}

	keys := binutils.SortedMapKeys(data)

	assert.IsIncreasing(t, keys)
	assert.Len(t, keys, 5000)
}
