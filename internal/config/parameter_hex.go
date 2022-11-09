package config

import (
	"context"
	"encoding/hex"
	"fmt"
)

const ParameterHex = "hex"

type paramHex struct {
	required bool
}

func (p paramHex) Required() bool {
	return p.required
}

func (p paramHex) Type() string {
	return ParameterHex
}

func (p paramHex) String() string {
	return fmt.Sprintf("required=%t", p.required)
}

func (p paramHex) Validate(_ context.Context, value string) error {
	if _, err := hex.DecodeString(value); err != nil {
		return fmt.Errorf("incorrectly encoded hex value: %w", err)
	}

	return nil
}

func NewHex(required bool, _ map[string]string) (Parameter, error) {
	return paramHex{
		required: required,
	}, nil
}
