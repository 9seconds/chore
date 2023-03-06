package v1

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNonce(t *testing.T) {
	assert.NotEqual(t, generateNonce(), generateNonce())
	assert.Len(t, generateNonce(), NonceLength)
}

func TestGenerateCipherKey(t *testing.T) {
	assert.Equal(
		t,
		"1Qzl39waGKJl5q+DgKT33WeSOlQcGDk913VI2yT6seM=",
		base64.StdEncoding.EncodeToString(
			generateCipherKey([]byte("aaa"), bytes.Repeat([]byte{1}, NonceLength)),
		))
}

func TestGenerateMacKey(t *testing.T) {
	assert.Equal(
		t,
		"Z0sx+YR7a+ZVfZtbArZ/wc4abHxGcutP8bSS3MCv1vE=",
		base64.StdEncoding.EncodeToString(
			generateMacKey([]byte("aaa"), bytes.Repeat([]byte{1}, NonceLength)),
		))
}

func TestEncryptDecrypt(t *testing.T) {
	testTable := map[string]string{
		"":                    "zOGkEVVmdYcSMHG+sw7VyA==",
		"12":                  "/NRzpIxIzjEQW27kKpA55w==",
		"1234":                "537ttjo3GSmrPhaju9gfnA==",
		"12456789":            "QIsOo6JjBS3D9thPAOFQgA==",
		"1234567890123456678": "mTL+dEOVv1wrR3cMYXNpTpX7Yp5fYcYifojHTX4P45U=",
	}

	key := bytes.Repeat([]byte{1}, KeyLength)

	for testValue, expected := range testTable {
		testValue := testValue
		expected := expected

		t.Run(testValue, func(t *testing.T) {
			encrypted := encryptMessage(key, []byte(testValue))
			assert.Equal(t, expected, base64.StdEncoding.EncodeToString(encrypted))

			decrypted, err := decryptMessage(key, encrypted)
			assert.NoError(t, err)
			assert.Equal(t, testValue, string(decrypted))
		})
	}
}
