package config

import (
	"encoding/json"
	"fmt"
	"io"
)

type RawConfig struct {
	Description string                  `json:"description"`
	Git         string                  `json:"git"`
	Network     bool                    `json:"network"`
	Parameters  map[string]RawParameter `json:"parameters"`
}

type RawParameter struct {
	Type     string            `json:"type"`
	Required bool              `json:"required"`
	Spec     map[string]string `json:"spec"`
}

func parseRaw(reader io.Reader) (RawConfig, error) {
	raw := RawConfig{}

	if err := json.NewDecoder(reader).Decode(&raw); err != nil {
		return raw, fmt.Errorf("cannot parse JSON config: %w", err)
	}

	return raw, nil
}
