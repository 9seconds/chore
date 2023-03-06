package vault

import (
	"bufio"
	"encoding"
	"errors"
	"fmt"
	"io"

	v1 "github.com/9seconds/chore/internal/vault/v1"
)

var (
	ErrEmptySecret             = errors.New("secret should not be empty")
	ErrUnsupportedVaultVersion = errors.New("vault version is not supported")
	LatestVersion              = 1
)

type Vault interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	Version() uint8
	List() []string
	Set(key, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

func Open(reader io.Reader, password string) (Vault, error) {
	bufReader := bufio.NewReader(reader)

	version, err := bufReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("cannot read version: %w", err)
	}

	data, err := io.ReadAll(bufReader)
	if err != nil {
		return nil, fmt.Errorf("cannot read data: %w", err)
	}

	var vault Vault

	switch version {
	case 1:
		vault, err = v1.NewVault(password)
	default:
		return nil, ErrUnsupportedVaultVersion
	}

	if err != nil {
		return nil, fmt.Errorf("cannot open vault of version %d: %w", int(version), err)
	}

	if err := vault.UnmarshalBinary(data); err != nil {
		return nil, fmt.Errorf("cannot unmarshal vault: %w", err)
	}

	return vault, nil
}

func New(password string) (Vault, error) {
	return v1.NewVault(password)
}

func Save(writer io.Writer, vault Vault) error {
	data, err := vault.MarshalBinary()
	if err != nil {
		return fmt.Errorf("cannot marshal vault: %w", err)
	}

	toWrite := []byte{byte(vault.Version())}
	toWrite = append(toWrite, data...)

	_, err = writer.Write(toWrite)

	return err
}
