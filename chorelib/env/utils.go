package env

import (
	"context"
	"crypto/rand"
	"encoding/base64"
)

func generateRandomString(length int) string {
	id := make([]byte, length)

	if _, err := rand.Read(id); err != nil {
		panic(err)
	}

	return encodeBytes(id)
}

func encodeBytes(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func makeEnvValue(name, value string) string {
	return name + "=" + value
}

func sendEnvValue(ctx context.Context, result chan<- string, name, value string) {
	select {
	case <-ctx.Done():
	case result <- makeEnvValue(name, value):
	}
}
