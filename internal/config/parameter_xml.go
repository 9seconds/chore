package config

import (
	"context"
	"encoding/xml"
	"fmt"
)

const ParameterXML = "xml"

type paramXML struct {
	baseParameter
}

func (p paramXML) Type() string {
	return ParameterXML
}

func (p paramXML) Validate(_ context.Context, value string) error {
	var doc interface{}

	if err := xml.Unmarshal([]byte(value), &doc); err != nil {
		return fmt.Errorf("invalid xml: %w", err)
	}

	return nil
}

func NewXML(description string, required bool, spec map[string]string) (Parameter, error) {
	return paramXML{
		baseParameter: baseParameter{
			required:      required,
			description:   description,
			specification: spec,
		},
	}, nil
}
