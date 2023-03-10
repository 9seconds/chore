package config

import (
	"context"
	"fmt"
	"net"
)

const ParameterMac = "mac"

type paramMac struct {
	baseParameter
}

func (p paramMac) Type() string {
	return ParameterMac
}

func (p paramMac) Validate(_ context.Context, value string) error {
	if _, err := net.ParseMAC(value); err != nil {
		return fmt.Errorf("incorrect mac address: %w", err)
	}

	return nil
}

func NewMac(description string, required bool, spec map[string]string) (Parameter, error) {
	return paramMac{
		baseParameter: baseParameter{
			required:      required,
			description:   description,
			specification: spec,
		},
	}, nil
}
