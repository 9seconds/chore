package config

import (
	"context"
	"fmt"
	"regexp"
)

const ParameterString = "string"

type paramString struct {
	mixinStringLength

	required bool
	re       *regexp.Regexp
}

func (p paramString) Required() bool {
	return p.required
}

func (p paramString) Type() string {
	return ParameterString
}

func (p paramString) String() string {
	return fmt.Sprintf("required=%t, re=%v, %s", p.required, p.re, p.mixinStringLength)
}

func (p paramString) Validate(_ context.Context, value string) error {
	if p.re != nil && !p.re.MatchString(value) {
		return fmt.Errorf("value %s does not match %s", value, p.re.String())
	}

	return p.mixinStringLength.validate(value)
}

func NewString(required bool, spec map[string]string) (Parameter, error) {
	param := paramString{
		required: required,
	}

	if stringLength, err := makeMixinStringLength(spec); err == nil {
		param.mixinStringLength = stringLength
	} else {
		return nil, err
	}

	if value, ok := spec["regexp"]; ok {
		re, err := regexp.Compile(value)
		if err != nil {
			return nil, fmt.Errorf("cannot compile regexp: %w", err)
		}

		param.re = re
	}

	return param, nil
}
