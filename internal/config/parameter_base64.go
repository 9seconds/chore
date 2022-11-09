package config

import (
	"context"
	"encoding/base64"
	"fmt"
)

const (
	ParameterBase64 = "base64"

	Base64EncRawStd  = "raw_std"
	Base64EncRawURL  = "raw_url"
	Base64EncStd     = "std"
	Base64EncUrl     = "url"
	Base64EncDefault = Base64EncStd
)

var (
	base64Encodings = map[string]*base64.Encoding{
		Base64EncRawStd: base64.RawStdEncoding,
		Base64EncRawURL: base64.RawURLEncoding,
		Base64EncStd:    base64.StdEncoding,
		Base64EncUrl:    base64.URLEncoding,
	}
)

type paramBase64 struct {
	required     bool
	encodingName string
	encoding     *base64.Encoding
}

func (p paramBase64) Required() bool {
	return p.required
}

func (p paramBase64) Type() string {
	return ParameterBase64
}

func (p paramBase64) String() string {
	return fmt.Sprintf("required=%t, encoding=%s", p.required, p.encodingName)
}

func (p paramBase64) Validate(_ context.Context, value string) error {
	if _, err := p.encoding.DecodeString(value); err != nil {
		return fmt.Errorf("incorrectly encoded value: %w", err)
	}

	return nil
}

func NewBase64(required bool, spec map[string]string) (Parameter, error) {
	param := paramBase64{
		required:     required,
		encodingName: spec["encoding"],
	}

	encoding := Base64EncDefault

	if value, ok := spec["encoding"]; ok {
		encoding = value
	}

	if enc, ok := base64Encodings[encoding]; ok {
		param.encoding = enc
	} else {
		return nil, fmt.Errorf("incorrect encoding %s", encoding)
	}

	return param, nil
}
