package env

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func MakeValue(name, value string) string {
	return name + "=" + value
}

func EncodeBytes(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func chainValues(value, upstream string) string {
	mac := hmac.New(sha256.New, []byte(value))
	mac.Write([]byte(upstream))

	return EncodeBytes(mac.Sum(nil))
}

func generateRandomString(length int) string {
	id := make([]byte, length)

	if _, err := rand.Read(id); err != nil {
		panic(err)
	}

	return EncodeBytes(id)
}

func sendValue(ctx context.Context, results chan<- string, name, value string) {
	select {
	case <-ctx.Done():
	case results <- MakeValue(name, value):
	}
}
