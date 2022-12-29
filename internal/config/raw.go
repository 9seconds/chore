package config

import (
	"fmt"
	"io"

	"github.com/hjson/hjson-go/v4"
)

type RawConfig struct {
	Description string                  `json:"description"`
	Git         string                  `json:"git"`
	Network     bool                    `json:"network"`
	AsUser      string                  `json:"as_user"`
	Parameters  map[string]RawParameter `json:"parameters"`
	Flags       map[string]RawFlag      `json:"flags"`
}

type RawParameter struct {
	Type        string            `json:"type"`
	Required    bool              `json:"required"`
	Description string            `json:"description"`
	Spec        map[string]string `json:"spec"`
}

type RawFlag struct {
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

func parseRaw(reader io.Reader) (RawConfig, error) {
	raw := RawConfig{}

	data, err := io.ReadAll(reader)
	if err != nil {
		return raw, fmt.Errorf("cannot read config: %w", err)
	}

	decoderOptions := hjson.DefaultDecoderOptions()
	decoderOptions.DisallowUnknownFields = true

	if err := hjson.UnmarshalWithOptions(data, &raw, decoderOptions); err != nil {
		return raw, fmt.Errorf("cannot parse JSON config: %w", err)
	}

	return raw, nil
}
