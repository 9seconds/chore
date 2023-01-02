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
	Base64EncURL     = "url"
	Base64EncDefault = Base64EncStd
)

var base64Encodings = map[string]*base64.Encoding{
	Base64EncRawStd: base64.RawStdEncoding,
	Base64EncRawURL: base64.RawURLEncoding,
	Base64EncStd:    base64.StdEncoding,
	Base64EncURL:    base64.URLEncoding,
}

type paramBase64 struct {
	baseParameter
	mixinStringLength

	encodingName string
	encoding     *base64.Encoding
}

func (p paramBase64) Type() string {
	return ParameterBase64
}

func (p paramBase64) Validate(_ context.Context, value string) error {
	if _, err := p.encoding.DecodeString(value); err != nil {
		return fmt.Errorf("incorrectly encoded value: %w", err)
	}

	return p.mixinStringLength.validate(value)
}

func NewBase64(description string, required bool, spec map[string]string) (Parameter, error) {
	param := paramBase64{
		baseParameter: baseParameter{
			required:      required,
			description:   description,
			specification: spec,
		},
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

	if mixin, err := makeMixinStringLength(spec); err == nil {
		param.mixinStringLength = mixin
	} else {
		return nil, err
	}

	return param, nil
}
