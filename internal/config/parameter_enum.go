package config

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

var (
	errInvalidChoice = errors.New("invalid choice")
	errNoChoices     = errors.New("no choices are prodvided")
)

const ParameterEnum = "enum"

type paramEnum struct {
	required bool
	choices  map[string]struct{}
}

func (p paramEnum) Required() bool {
	return p.required
}

func (p paramEnum) Type() string {
	return ParameterEnum
}

func (p paramEnum) String() string {
	choices := make([]string, 0, len(p.choices))

	for k := range p.choices {
		choices = append(choices, k)
	}

	sort.Strings(choices)

	return fmt.Sprintf("required=%t, choices=%v", p.required, choices)
}

func (p paramEnum) Validate(_ context.Context, value string) error {
	if _, ok := p.choices[value]; !ok {
		return errInvalidChoice
	}

	return nil
}

func NewEnum(required bool, spec map[string]string) (Parameter, error) {
	param := paramEnum{
		required: required,
		choices:  make(map[string]struct{}),
	}

	for _, v := range parseCSV(spec["choices"]) {
		param.choices[v] = struct{}{}
	}

	if len(param.choices) == 0 {
		return param, errNoChoices
	}

	return param, nil
}
