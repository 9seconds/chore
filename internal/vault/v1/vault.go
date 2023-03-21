package v1

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

// format is very simple:
// [kdf nonce | hmac | encrypted message length | encrypted message]
//
// kdf nonce is 16 bytes. Used to produce AES key from a password
// hmac is HMAC-SHA256 of the rest of the message. Message is encrypted by AES-256-CBC. 32 bytes
// encrypted message length is little endian uint32 of a message length. If you can't fit
//    into this size, you are probably doing something very wrong.
// encrypted message is PKCS7-padded JSON encoded mapping of string to string

const (
	NonceLength = 16
	LenLength   = 4  // uint32
	MACLength   = 32 // sha256
)

var (
	ErrEmptyPassword    = errors.New("password is empty")
	ErrBadPassword      = errors.New("bad password")
	ErrShortData        = errors.New("encrypted data is short")
	ErrIncorrectPadding = errors.New("data is incorrectly padded")
)

type Vault struct {
	password []byte
	data     map[string]string
}

func (v *Vault) UnmarshalBinary(data []byte) error {
	if len(data) < NonceLength {
		return fmt.Errorf("cannot read KDF nonce: %w", ErrShortData)
	}

	kdfNonce, data := data[:NonceLength], data[NonceLength:]

	if len(data) < MACLength {
		return fmt.Errorf("cannot read MAC: %w", ErrShortData)
	}

	mac, data := data[:MACLength], data[MACLength:]

	if len(data) < LenLength {
		return fmt.Errorf("cannot read length: %w", ErrShortData)
	}

	length := int(binary.LittleEndian.Uint32(data[:LenLength]))
	if len(data)-LenLength != length {
		return fmt.Errorf("message length mismatch: %w", ErrShortData)
	}

	key := generateKey(v.password, kdfNonce)

	macMixer := hmac.New(sha256.New, key)
	macMixer.Write(data)

	if subtle.ConstantTimeCompare(mac, macMixer.Sum(nil)) != 1 {
		return ErrBadPassword
	}

	message, err := decryptMessage(key, data[LenLength:])
	if err != nil {
		return ErrBadPassword
	}

	v.data = make(map[string]string)

	return json.Unmarshal(message, &v.data)
}

func (v *Vault) MarshalBinary() ([]byte, error) {
	kdfNonce := generateNonce()
	key := generateKey(v.password, kdfNonce)

	data, err := json.Marshal(v.data)
	if err != nil {
		panic(err.Error())
	}

	data = encryptMessage(key, data)

	length := make([]byte, LenLength)
	binary.LittleEndian.PutUint32(length, uint32(len(data)))

	macMixer := hmac.New(sha256.New, key)
	macMixer.Write(length)
	macMixer.Write(data)

	result := make([]byte, 0, NonceLength+MACLength+LenLength+len(data))
	result = append(result, kdfNonce...)
	result = append(result, macMixer.Sum(nil)...)
	result = append(result, length...)
	result = append(result, data...)

	return result, nil
}

func (v *Vault) Version() uint8 {
	return 1
}

func (v *Vault) List() []string {
	items := make([]string, 0, len(v.data))

	for k := range v.data {
		items = append(items, k)
	}

	return items
}

func (v *Vault) Set(key, value string) {
	v.data[key] = value
}

func (v *Vault) Get(key string) (string, bool) {
	value, ok := v.data[key]

	return value, ok
}

func (v *Vault) Delete(key string) {
	delete(v.data, key)
}

func NewVault(password string) (*Vault, error) {
	if password == "" {
		return nil, ErrEmptyPassword
	}

	return &Vault{
		password: []byte(password),
		data:     make(map[string]string),
	}, nil
}
