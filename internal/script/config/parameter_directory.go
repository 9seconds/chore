package config

import (
	"context"
	"errors"
)

const ParameterDirectory = "directory"

var errIsNotDirectory = errors.New("is not a directory")

type paramDirectory struct {
	baseParameter
	mixinPermissions
}

func (p paramDirectory) Type() string {
	return ParameterDirectory
}

func (p paramDirectory) Validate(_ context.Context, value string) error {
	stat, err := p.mixinPermissions.validate(value, p.isExist())

	switch {
	case err != nil:
		return err
	case stat != nil && !stat.IsDir():
		return errIsNotDirectory
	}

	return nil
}

func NewDirectory(description string, required bool, spec map[string]string) (Parameter, error) {
	param := paramDirectory{
		baseParameter: baseParameter{
			description:   description,
			required:      required,
			specification: spec,
		},
	}

	mixin, err := makeMixinPermissions(spec)
	if err != nil {
		return param, err
	}

	param.mixinPermissions = mixin

	return param, nil
}
