package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
)

const ParameterUUID = "uuid"

var validUUIDVersions = map[byte]bool{
	1: true,
	3: true,
	4: true,
	5: true,
	6: true,
	7: true,
}

type paramUUID struct {
	required bool
	version  byte
}

func (p paramUUID) Required() bool {
	return p.required
}

func (p paramUUID) Type() string {
	return ParameterUUID
}

func (p paramUUID) String() string {
	return fmt.Sprintf("required=%t, version=%d", p.required, p.version)
}

func (p paramUUID) Validate(_ context.Context, value string) error {
	parsed, err := uuid.FromString(value)
	if err != nil {
		return fmt.Errorf("cannot parse uuid: %w", err)
	}

	if p.version != 0 {
		if pv := parsed.Version(); pv != p.version {
			return fmt.Errorf("incorrect uuid version %d, expected %d", pv, p.version)
		}
	}

	return nil
}

func NewUUID(required bool, spec map[string]string) (Parameter, error) {
	param := paramUUID{
		required: required,
	}

	if value, ok := spec["version"]; ok {
		parsed, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("incorrect version: %w", err)
		}

		if _, ok := validUUIDVersions[byte(parsed)]; !ok {
			return nil, fmt.Errorf("incorrect version %d", parsed)
		}

		param.version = byte(parsed)
	}

	return param, nil
}
