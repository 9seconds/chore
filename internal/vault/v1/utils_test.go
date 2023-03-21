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

func TestGenerateKey(t *testing.T) {
	key := generateKey([]byte("aaa"), bytes.Repeat([]byte{1}, NonceLength))

	assert.Equal(
		t,
		"+3UvJQj/0KdzvXx1+5NVpjBFGX8iwh1oIoD5YwkA+cc=",
		base64.StdEncoding.EncodeToString(key))
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
