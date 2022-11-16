package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
)

const ParameterFile = "file"

var (
	errIsNotFile = errors.New("is not a file")
)

type paramFile struct {
	mixinPermissions

	required  bool
	mimetypes []string
}

func (p paramFile) Required() bool {
	return p.required
}

func (p paramFile) Type() string {
	return ParameterFile
}

func (p paramFile) String() string {
	return fmt.Sprintf(
		"required=%t, mimetypes=%v, %s",
		p.required,
		p.mimetypes,
		p.mixinPermissions)
}

func (p paramFile) isExist() bool {
	return p.mixinPermissions.isExist() || len(p.mimetypes) > 0
}

func (p paramFile) Validate(_ context.Context, value string) error {
	stat, err := p.mixinPermissions.validate(value, p.isExist())

	switch {
	case err != nil:
		return err
	case stat == nil:
		return nil
	case !stat.Mode().IsRegular():
		return errIsNotFile
	case len(p.mimetypes) > 0:
		mtype, err := mimetype.DetectFile(value)
		if err != nil {
			return fmt.Errorf("cannot detect mimetype: %w", err)
		}

		if !mimetype.EqualsAny(mtype.String(), p.mimetypes...) {
			return fmt.Errorf("unexpected mimetype %s", mtype.String())
		}
	}

	return nil
}

func NewFile(required bool, spec map[string]string) (Parameter, error) {
	param := paramFile{
		required: required,
	}

	mixin, err := makeMixinPermissions(
		spec,
		"exists",
		"readable",
		"writable",
		"executable")
	if err != nil {
		return param, err
	}

	param.mixinPermissions = mixin

	for _, v := range parseCSV(spec["mimetypes"]) {
		if mimetype.Lookup(v) == nil {
			return nil, fmt.Errorf("unknown mimetype %s", v)
		}

		param.mimetypes = append(param.mimetypes, v)
	}

	return param, nil
}