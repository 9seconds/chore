package config

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
)

const ParameterSemver = "semver"

type parameterSemver struct {
	baseParameter

	constraint *semver.Constraints
}

func (p parameterSemver) Type() string {
	return ParameterSemver
}

func (p parameterSemver) Validate(_ context.Context, value string) error {
	ver, err := semver.NewVersion(value)
	if err != nil {
		return fmt.Errorf("incorrect semver: %w", err)
	}

	if p.constraint == nil {
		return nil
	}

	if ok, errors := p.constraint.Validate(ver); !ok && len(errors) > 0 {
		return fmt.Errorf("invalid version: %w", errors[0])
	}

	return nil
}

func NewSemver(description string, required bool, spec map[string]string) (Parameter, error) {
	param := parameterSemver{
		baseParameter: baseParameter{
			required:      required,
			description:   description,
			specification: spec,
		},
	}

	if value, ok := spec["constraint"]; ok {
		cs, err := semver.NewConstraint(value)
		if err != nil {
			return nil, fmt.Errorf("incorrect constraint %s: %w", value, err)
		}

		param.constraint = cs
	}

	return param, nil
}
