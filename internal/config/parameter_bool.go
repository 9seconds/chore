package config

import (
	"context"
	"fmt"
	"strconv"
)

const ParameterBool = "bool"

type parameterBool struct {
	required bool
}

func (p parameterBool) Required() bool {
	return p.required
}

func (p parameterBool) Type() string {
	return ParameterBool
}

func (p parameterBool) String() string {
	return fmt.Sprintf("required=%t", p.required)
}

func (p parameterBool) Validate(_ context.Context, value string) error {
	_, err := strconv.ParseBool(value)

	return err
}

func NewBool(required bool, _ map[string]string) (Parameter, error) {
	return parameterBool{
		required: required,
	}, nil
}
