package config

import (
	"context"
	"fmt"
	"net"
)

const ParameterMac = "mac"

type paramMac struct {
	required bool
}

func (p paramMac) Required() bool {
	return p.required
}

func (p paramMac) Type() string {
	return ParameterMac
}

func (p paramMac) String() string {
	return fmt.Sprintf("required=%t", p.required)
}

func (p paramMac) Validate(_ context.Context, value string) error {
	if _, err := net.ParseMAC(value); err != nil {
		return fmt.Errorf("incorrect mac address: %w", err)
	}

	return nil
}

func NewMac(required bool, _ map[string]string) (Parameter, error) {
	return paramMac{
		required: required,
	}, nil
}
