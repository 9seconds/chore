package v1

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const (
	KeyLength = 32 // aes256

	ArgonTime    = 1
	ArgonMemory  = 64 * 1024
	ArgonThreads = 4
)

// constant IV is fine because we rotate keys each time, keys is
// defined by a password and KDF nonce. Since KDF nonce is regenerated,
// we do not have a persistent key. Thus, an attack on a first block is
// not relevan for us. it means, we can skip storing and generating IV
// https://stackoverflow.com/a/2648345
var ConstantIV [16]byte

func encryptMessage(key, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	encryptor := cipher.NewCBCEncrypter(block, ConstantIV[:])
	paddingLength := aes.BlockSize - len(data)%aes.BlockSize
	data = append(data, bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)...)

	encryptor.CryptBlocks(data, data)

	return data
}

func decryptMessage(key, data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, ErrShortData
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	decryptor := cipher.NewCBCDecrypter(block, ConstantIV[:])
	decryptor.CryptBlocks(data, data)

	length := int(data[len(data)-1])
	if length == 0 || length > aes.BlockSize {
		return nil, ErrIncorrectPadding
	}

	suffix := bytes.Repeat([]byte{byte(length)}, length)
	if !bytes.HasSuffix(data, suffix) {
		return nil, ErrIncorrectPadding
	}

	return data[:len(data)-length], nil
}

func generateNonce() []byte {
	data := make([]byte, NonceLength)

	if _, err := rand.Read(data); err != nil {
		panic(err.Error())
	}

	return data
}

func generateKey(password, nonce []byte) []byte {
	return argon2.IDKey(
		password,
		nonce,
		ArgonTime,
		ArgonMemory,
		ArgonThreads,
		KeyLength)
}
