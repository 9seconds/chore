package config

import (
	"context"
	"errors"
)

var (
	errInvalidChoice = errors.New("invalid choice")
	errNoChoices     = errors.New("no choices are prodvided")
)

const ParameterEnum = "enum"

type paramEnum struct {
	baseParameter

	choices map[string]struct{}
}

func (p paramEnum) Type() string {
	return ParameterEnum
}

func (p paramEnum) Validate(_ context.Context, value string) error {
	if _, ok := p.choices[value]; !ok {
		return errInvalidChoice
	}

	return nil
}

func NewEnum(description string, required bool, spec map[string]string) (Parameter, error) {
	param := paramEnum{
		baseParameter: baseParameter{
			required:      required,
			description:   description,
			specification: spec,
		},
		choices: make(map[string]struct{}),
	}

	for _, v := range parseCSV(spec["choices"]) {
		param.choices[v] = struct{}{}
	}

	if len(param.choices) == 0 {
		return param, errNoChoices
	}

	return param, nil
}
