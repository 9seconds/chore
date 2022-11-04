package config

import (
	"fmt"
	"math"
	"strconv"
)

const ParameterFloat = "float"

var (
	paramFloatMax = func() float64 {
		val, err := strconv.ParseFloat("+inf", 64)
		if err != nil {
			panic(err)
		}

		return val
	}()

	paramFloatMin = func() float64 {
		val, err := strconv.ParseFloat("-inf", 64)
		if err != nil {
			panic(err)
		}

		return val
	}()
)

type paramFloat struct {
	required bool
	min      float64
	max      float64
}

func (p paramFloat) Required() bool {
	return p.required
}

func (p paramFloat) Type() string {
	return ParameterFloat
}

func (p paramFloat) String() string {
	return fmt.Sprintf("required=%t, min=%v, max=%v", p.required, p.min, p.max)
}

func (p paramFloat) Validate(value string) error {
	parsed, err := strconv.ParseFloat(value, 64)

	switch {
	case err != nil:
		return fmt.Errorf("cannot parse float: %w", err)
	case math.IsInf(parsed, 0), math.IsNaN(parsed):
		return fmt.Errorf("%s is unacceptable float", value)
	case parsed < p.min:
		return fmt.Errorf("value is less than minimum %v", p.min)
	case parsed > p.max:
		return fmt.Errorf("value is more than maximum %v", p.max)
	}

	return nil
}

func NewFloat(required bool, spec map[string]string) (Parameter, error) {
	rValue := paramFloat{
		required: required,
		min:      paramFloatMin,
		max:      paramFloatMax,
	}

	if strValue, ok := spec["min"]; ok {
		value, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'min' %s as float: %w", strValue, err)
		}

		if math.IsNaN(value) {
			return nil, fmt.Errorf("cannot use %s as min", strValue)
		}

		rValue.min = value
	}

	if strValue, ok := spec["max"]; ok {
		value, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'max' %s as float: %w", strValue, err)
		}

		if math.IsNaN(value) {
			return nil, fmt.Errorf("cannot use %s as max", strValue)
		}

		rValue.max = value
	}

	if rValue.min > rValue.max {
		return nil, fmt.Errorf("'max' %s value should be bigger than 'min' %s", spec["max"], spec["min"])
	}

	return rValue, nil
}
