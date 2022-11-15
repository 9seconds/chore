package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/9seconds/chore/internal/access"
)

var (
	errDoesNotExist = errors.New("does not exist")
)

type mixinPermissions struct {
	exists     bool
	readable   bool
	writable   bool
	executable bool
	mode       fs.FileMode
}

func (m mixinPermissions) String() string {
	return fmt.Sprintf(
		"mode=%s, exists=%t, readable=%t, writable=%t, executable=%t",
		m.mode,
		m.exists,
		m.readable,
		m.writable,
		m.executable)
}

func (m mixinPermissions) isExist() bool {
	return m.exists || m.readable || m.writable || m.executable || m.mode != 0
}

func (m mixinPermissions) validate(value string, exists bool) (fs.FileInfo, error) {
	stat, err := os.Stat(value)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		stat = nil

		if exists {
			return nil, errDoesNotExist
		}
	case err != nil:
		return nil, fmt.Errorf("cannot stat: %w", err)
	case m.mode != 0 && stat.Mode().Perm() != m.mode:
		return nil, fmt.Errorf(
			"incorrect mode: got %s, expected %s",
			stat.Mode().Perm(),
			m.mode)
	case m.readable, m.writable, m.executable:
		if err := access.Access(value, m.readable, m.writable, m.executable); err != nil {
			return nil, fmt.Errorf("incorrect user permissions: %w", err)
		}
	}

	return stat, nil
}

func makeMixinPermissions(spec map[string]string,
	existsName, readableName, writableName, executableName string) (mixinPermissions, error) {
	mixin := mixinPermissions{}

	if value, err := parseBool(spec, "exists"); err == nil {
		mixin.exists = value
	} else {
		return mixin, err
	}

	if value, err := parseBool(spec, "readable"); err == nil {
		mixin.readable = value
	} else {
		return mixin, err
	}

	if value, err := parseBool(spec, "writable"); err == nil {
		mixin.writable = value
	} else {
		return mixin, err
	}

	if value, err := parseBool(spec, "executable"); err == nil {
		mixin.executable = value
	} else {
		return mixin, err
	}

	if value, ok := spec["mode"]; ok {
		parsed, err := strconv.ParseUint(value, 8, 32)
		if err != nil {
			return mixin, fmt.Errorf("cannot parse mode: %w", err)
		}

		mixin.mode = fs.FileMode(parsed)
	}

	return mixin, nil
}
