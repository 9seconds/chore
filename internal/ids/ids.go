package ids

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"io"

	"github.com/rs/xid"
)

func New() string {
	return xid.New().String()
}

func Chain(base string, values ...string) string {
	mac := hmac.New(sha256.New, []byte(base))

	binary.Write(mac, binary.LittleEndian, len(values))

	for _, v := range values {
		binary.Write(mac, binary.LittleEndian, len(v))
	}

	for _, v := range values {
		io.WriteString(mac, v)
	}

	return Encode(mac.Sum(nil))
}

func Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
