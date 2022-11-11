package config

import (
	"context"
	"encoding/xml"
	"fmt"
)

const ParameterXML = "xml"

type paramXML struct {
	required bool
}

func (p paramXML) Required() bool {
	return p.required
}

func (p paramXML) Type() string {
	return ParameterXML
}

func (p paramXML) String() string {
	return fmt.Sprintf("required=%t", p.required)
}

func (p paramXML) Validate(_ context.Context, value string) error {
	var doc interface{}

	if err := xml.Unmarshal([]byte(value), &doc); err != nil {
		return fmt.Errorf("invalid xml: %w", err)
	}

	return nil
}

func NewXML(required bool, spec map[string]string) (Parameter, error) {
	return paramXML{
		required: required,
	}, nil
}
