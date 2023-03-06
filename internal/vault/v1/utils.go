package v1

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/scrypt"
)

const (
	KeyLength = 32

	CipherKeyN = 1 << 12
	CipherKeyR = 2
	CipherKeyP = 8

	MACKeyN = 1 << 8
	MACKeyR = 4
	MACKeyP = 10
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

func generateCipherKey(password, nonce []byte) []byte {
	return generateKey(password, nonce, CipherKeyN, CipherKeyR, CipherKeyP)
}

func generateMacKey(password, nonce []byte) []byte {
	return generateKey(password, nonce, MACKeyN, MACKeyP, MACKeyP)
}

func generateKey(password, nonce []byte, n, r, p int) []byte {
	key, err := scrypt.Key(password, nonce, n, r, p, KeyLength)
	if err != nil {
		panic(err.Error())
	}

	return key
}
