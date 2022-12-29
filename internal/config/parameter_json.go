package config

import (
	"context"
	"encoding/json"
	"fmt"
)

const ParameterJSON = "json"

type paramJSON struct {
	baseParameter
}

func (p paramJSON) Type() string {
	return ParameterJSON
}

func (p paramJSON) String() string {
	return fmt.Sprintf("%q (required=%t)", p.description, p.required)
}

func (p paramJSON) Validate(_ context.Context, value string) error {
	var doc interface{}

	if err := json.Unmarshal([]byte(value), &doc); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	return nil
}

func NewJSON(description string, required bool, spec map[string]string) (Parameter, error) {
	return paramJSON{
		baseParameter: baseParameter{
			required:    required,
			description: description,
		},
	}, nil
}
