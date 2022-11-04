package config

import (
	"fmt"
	"math"
	"strconv"
)

const ParameterInteger = "integer"

type paramInteger struct {
	required bool
	min      int64
	max      int64
}

func (p paramInteger) Required() bool {
	return p.required
}

func (p paramInteger) Type() string {
	return ParameterInteger
}

func (p paramInteger) String() string {
	return fmt.Sprintf("required=%t, min=%d, max=%d", p.required, p.min, p.max)
}

func (p paramInteger) Validate(value string) error {
	parsed, err := strconv.ParseInt(value, 10, 64)

	switch {
	case err != nil:
		return fmt.Errorf("cannot parse as integer: %w", err)
	case parsed < p.min:
		return fmt.Errorf("value is less than minimum %d", p.min)
	case parsed > p.max:
		return fmt.Errorf("value is bigger than maximum %d", p.max)
	}

	return nil
}

func NewInteger(required bool, spec map[string]string) (Parameter, error) {
	rValue := paramInteger{
		required: required,
		min:      math.MinInt64,
		max:      math.MaxInt64,
	}

	if strValue, ok := spec["min"]; ok {
		value, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'min' %s as integer: %w", strValue, err)
		}

		rValue.min = value
	}

	if strValue, ok := spec["max"]; ok {
		value, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'max' %s as integer: %w", strValue, err)
		}

		rValue.max = value
	}

	if rValue.min > rValue.max {
		return nil, fmt.Errorf("'max' %s value should be bigger than 'min' %s", spec["max"], spec["min"])
	}

	return rValue, nil
}
