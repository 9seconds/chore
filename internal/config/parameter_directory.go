package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/9seconds/chore/internal/access"
)

const ParameterDirectory = "directory"

var (
	errDoesNotExist   = errors.New("does not exist")
	errIsNotDirectory = errors.New("is not directory")
)

type paramDirectory struct {
	required   bool
	exists     bool
	readable   bool
	writable   bool
	executable bool
	mode       fs.FileMode
}

func (p paramDirectory) Required() bool {
	return p.required
}

func (p paramDirectory) Type() string {
	return ParameterDirectory
}

func (p paramDirectory) String() string {
	return fmt.Sprintf(
		"required=%t, exists=%t, readable=%t, writable=%t, executable=%t, mode=%s",
		p.required,
		p.exists,
		p.readable,
		p.writable,
		p.executable,
		p.mode)
}

func (p paramDirectory) Validate(_ context.Context, value string) error {
	stat, err := os.Stat(value)

	switch {
	case os.IsNotExist(err):
		if p.exists {
			return errDoesNotExist
		}

		return nil
	case err != nil:
		return fmt.Errorf("cannot stat %w", err)
	case !stat.IsDir():
		return errIsNotDirectory
	case p.mode != 0 && stat.Mode().Perm() != p.mode:
		return fmt.Errorf(
			"incorrect mode: got %s, expected %s",
			stat.Mode().Perm(),
			p.mode)
	case p.readable, p.writable, p.executable:
		if err := access.Access(value, p.readable, p.writable, p.executable); err != nil {
			return fmt.Errorf("incorrect user permissions: %w", err)
		}
	}

	return nil
}

func NewDirectory(required bool, spec map[string]string) (Parameter, error) {
	param := paramDirectory{
		required: required,
	}

	if value, err := parseBool(spec, "exists"); err == nil {
		param.exists = value
	} else {
		return nil, err
	}

	if value, err := parseBool(spec, "readable"); err == nil {
		param.readable = value
	} else {
		return nil, err
	}

	if value, err := parseBool(spec, "writable"); err == nil {
		param.writable = value
	} else {
		return nil, err
	}

	if value, err := parseBool(spec, "executable"); err == nil {
		param.executable = value
	} else {
		return nil, err
	}

	if value, ok := spec["mode"]; ok {
		parsed, err := strconv.ParseUint(value, 8, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot parse mode: %w", err)
		}

		param.mode = fs.FileMode(parsed)
	}

	return param, nil
}
