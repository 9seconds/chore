package config

import (
	"fmt"
	"regexp"
)

type parameterString struct {
	parameterBase

	re *regexp.Regexp
}

func (p parameterString) Type() string {
	return ParameterString
}

func (p parameterString) String() string {
	value := p.parameterBase.String()

	if p.re != nil {
		value += ", regexp=" + p.re.String()
	}

	return value

}

func (p parameterString) Validate(value string) error {
	if p.re == nil || p.re.MatchString(value) {
		return nil
	}

	return fmt.Errorf("value does not match %s", p.re.String())
}

func NewString(required bool, spec map[string]string) (Parameter, error) {
	rv := parameterString{
		parameterBase: parameterBase{
			required: required,
		},
	}

	if value, ok := spec["regexp"]; ok {
		re, err := regexp.Compile(value)
		if err != nil {
			return nil, fmt.Errorf("cannot compile regexp: %w", err)
		}

		rv.re = re
	}

	return rv, nil
}
