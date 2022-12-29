package config

import (
	"context"
	"fmt"
	"regexp"
	"unicode"
)

const ParameterString = "string"

type paramString struct {
	baseParameter
	mixinStringLength

	ascii bool
	re    *regexp.Regexp
}

func (p paramString) Type() string {
	return ParameterString
}

func (p paramString) String() string {
	return fmt.Sprintf(
		"%q (required=%t, ascii=%t, re=%v, %s)",
		p.description,
		p.required,
		p.ascii,
		p.re,
		p.mixinStringLength)
}

func (p paramString) Validate(_ context.Context, value string) error {
	if p.re != nil && !p.re.MatchString(value) {
		return fmt.Errorf("value %s does not match %s", value, p.re.String())
	}

	if p.ascii {
		for _, char := range value {
			if char > unicode.MaxASCII {
				return fmt.Errorf("value %s contains non-ascii character", value)
			}
		}
	}

	return p.mixinStringLength.validate(value)
}

func NewString(description string, required bool, spec map[string]string) (Parameter, error) {
	param := paramString{
		baseParameter: baseParameter{
			required:    required,
			description: description,
		},
	}

	if stringLength, err := makeMixinStringLength(spec); err == nil {
		param.mixinStringLength = stringLength
	} else {
		return nil, err
	}

	if value, err := parseRegexp(spec, "regexp"); err == nil {
		param.re = value
	} else {
		return nil, err
	}

	if value, err := parseBool(spec, "ascii"); err == nil {
		param.ascii = value
	} else {
		return nil, err
	}

	return param, nil
}
