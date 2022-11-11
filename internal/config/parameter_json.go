package config

import (
	"context"
	"encoding/json"
	"fmt"
)

const ParameterJSON = "json"

type paramJSON struct {
	required bool
}

func (p paramJSON) Required() bool {
	return p.required
}

func (p paramJSON) Type() string {
	return ParameterJSON
}

func (p paramJSON) String() string {
	return fmt.Sprintf("required=%t", p.required)
}

func (p paramJSON) Validate(_ context.Context, value string) error {
	var doc interface{}

	if err := json.Unmarshal([]byte(value), &doc); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	return nil
}

func NewJSON(required bool, spec map[string]string) (Parameter, error) {
	return paramJSON{
		required: required,
	}, nil
}
