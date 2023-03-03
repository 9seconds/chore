package binutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func Chain(base string, values ...string) string {
	mac := hmac.New(sha256.New, []byte(base))

	MixStringSlice(mac, values)

	return ToString(mac.Sum(nil))
}

func ToString(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
