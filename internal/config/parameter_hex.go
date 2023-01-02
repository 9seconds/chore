package config

import (
	"context"
	"encoding/hex"
	"fmt"
)

const ParameterHex = "hex"

type paramHex struct {
	baseParameter
	mixinStringLength
}

func (p paramHex) Type() string {
	return ParameterHex
}

func (p paramHex) Validate(_ context.Context, value string) error {
	if _, err := hex.DecodeString(value); err != nil {
		return fmt.Errorf("incorrectly encoded hex value: %w", err)
	}

	return p.mixinStringLength.validate(value)
}

func NewHex(description string, required bool, spec map[string]string) (Parameter, error) {
	param := paramHex{
		baseParameter: baseParameter{
			required:      required,
			description:   description,
			specification: spec,
		},
	}

	if stringLength, err := makeMixinStringLength(spec); err == nil {
		param.mixinStringLength = stringLength
	} else {
		return nil, err
	}

	return param, nil
}
