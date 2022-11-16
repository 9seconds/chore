package config

import (
	"context"
	"errors"
	"fmt"
)

const ParameterDirectory = "directory"

var errIsNotDirectory = errors.New("is not a directory")

type paramDirectory struct {
	mixinPermissions

	required bool
}

func (p paramDirectory) Required() bool {
	return p.required
}

func (p paramDirectory) Type() string {
	return ParameterDirectory
}

func (p paramDirectory) String() string {
	return fmt.Sprintf("required=%t, %s", p.required, p.mixinPermissions)
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

func NewDirectory(required bool, spec map[string]string) (Parameter, error) {
	param := paramDirectory{
		required: required,
	}

	mixin, err := makeMixinPermissions(spec)
	if err != nil {
		return param, err
	}

	param.mixinPermissions = mixin

	return param, nil
}
