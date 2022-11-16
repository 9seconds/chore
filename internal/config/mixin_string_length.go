package config

import (
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"
)

type mixinStringLength struct {
	minValue int
	maxValue int
}

func (m mixinStringLength) String() string {
	return fmt.Sprintf("min=%d, max=%d", m.minValue, m.maxValue)
}

func (m mixinStringLength) validate(value string) error {
	length := utf8.RuneCountInString(value)

	switch {
	case length < m.minValue:
		return fmt.Errorf("value length must be >= %d", m.minValue)
	case length > m.maxValue:
		return fmt.Errorf("value length must be <= %d", m.maxValue)
	}

	return nil
}

func makeMixinStringLength(spec map[string]string) (mixinStringLength, error) {
	rValue := mixinStringLength{
		minValue: 0,
		maxValue: math.MaxInt,
	}

	if min, ok := spec["min_length"]; ok {
		value, err := strconv.ParseUint(min, 10, 64)
		if err != nil {
			return rValue, fmt.Errorf("incorrect 'min_length' value: %w", err)
		}

		rValue.minValue = int(value)
	}

	if max, ok := spec["max_length"]; ok {
		value, err := strconv.ParseUint(max, 10, 64)
		if err != nil {
			return rValue, fmt.Errorf("incorrect 'max_length' value: %w", err)
		}

		rValue.maxValue = int(value)
	}

	if rValue.maxValue < rValue.minValue {
		return rValue, fmt.Errorf(
			"min %d should be <= max %d",
			rValue.minValue,
			rValue.maxValue)
	}

	return rValue, nil
}
