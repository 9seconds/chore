package binutils_test

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/9seconds/chore/internal/binutils"
	"github.com/stretchr/testify/assert"
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

func TestMixStringMap(t *testing.T) {
	testTable := map[string]string{
		"":       "0301000000000000000102000000000000006b5f010000000000000000",
		"aa":     "0301000000000000000104000000000000006b5f61610102000000000000006161",
		"aa,bb":  "0302000000000000000104000000000000006b5f616101020000000000000061610104000000000000006b5f62620102000000000000006262",
		"aa,b,b": "0302000000000000000104000000000000006b5f616101020000000000000061610103000000000000006b5f6201010000000000000062",
	}

	for testValue, expected := range testTable {
		testValue := testValue
		expected := expected

		t.Run(testValue, func(t *testing.T) {
			buf := &bytes.Buffer{}

			values := make(map[string]string)

			for _, v := range strings.Split(testValue, ",") {
				values["k_"+v] = v
			}

			assert.NoError(t, binutils.MixStringsMap(buf, values))
			assert.Equal(t, expected, hex.EncodeToString(buf.Bytes()))
		})
	}
}
