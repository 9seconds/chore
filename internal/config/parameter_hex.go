package config

import (
	"context"
	"encoding/hex"
	"fmt"
)

const ParameterHex = "hex"

type paramHex struct {
	mixinStringLength

	required bool
}

func (p paramHex) Required() bool {
	return p.required
}

func (p paramHex) Type() string {
	return ParameterHex
}

func (p paramHex) String() string {
	return fmt.Sprintf("required=%t, %s", p.required, p.mixinStringLength)
}

func (p paramHex) Validate(_ context.Context, value string) error {
	if _, err := hex.DecodeString(value); err != nil {
		return fmt.Errorf("incorrectly encoded hex value: %w", err)
	}

	return p.mixinStringLength.Validate(value)
}

func NewHex(required bool, spec map[string]string) (Parameter, error) {
	param := paramHex{
		required: required,
	}

	if stringLength, err := makeMixinStringLength(spec, "min_length", "max_length"); err == nil {
		param.mixinStringLength = stringLength
	} else {
		return nil, err
	}

	return param, nil
}
