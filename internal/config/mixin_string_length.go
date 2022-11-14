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

func (m mixinStringLength) Validate(value string) error {
	length := utf8.RuneCountInString(value)

	switch {
	case length < m.minValue:
		return fmt.Errorf("value length must be >= %d", m.minValue)
	case length > m.maxValue:
		return fmt.Errorf("value length must be <= %d", m.maxValue)
	}

	return nil
}

func makeMixinStringLength(spec map[string]string, minName, maxName string) (mixinStringLength, error) {
	rValue := mixinStringLength{
		minValue: 0,
		maxValue: math.MaxInt,
	}

	if min, ok := spec[minName]; ok {
		if value, err := strconv.ParseUint(min, 10, 64); err != nil {
			return rValue, fmt.Errorf("incorrect %s value: %w", minName, err)
		} else {
			rValue.minValue = int(value)
		}
	}

	if max, ok := spec[maxName]; ok {
		if value, err := strconv.ParseUint(max, 10, 64); err != nil {
			return rValue, fmt.Errorf("incorrect %s value: %w", maxName, err)
		} else {
			rValue.maxValue = int(value)
		}
	}

	if rValue.maxValue < rValue.minValue {
		return rValue, fmt.Errorf(
			"min %d should be <= max %d",
			rValue.minValue,
			rValue.maxValue)
	}

	return rValue, nil
}
