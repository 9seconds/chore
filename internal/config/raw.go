package config

import (
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"
)

type RawConfig struct {
	Description string                  `toml:"description"`
	Git         string                  `toml:"git"`
	Network     bool                    `toml:"network"`
	Parameters  map[string]RawParameter `toml:"parameters"`
	Flags       map[string]RawFlag      `toml:"flags"`
}

type RawParameter struct {
	Type        string            `toml:"type"`
	Required    bool              `toml:"required"`
	Description string            `toml:"description"`
	Spec        map[string]string `toml:"spec"`
}

type RawFlag struct {
	Required    bool   `toml:"required"`
	Description string `toml:"description"`
}

func parseRaw(reader io.Reader) (RawConfig, error) {
	raw := RawConfig{}

	decoder := toml.NewDecoder(reader)
	decoder = decoder.DisallowUnknownFields()

	if err := decoder.Decode(&raw); err != nil {
		return raw, fmt.Errorf("cannot parse TOML config: %w", err)
	}

	return raw, nil
}
