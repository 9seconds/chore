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
	Parameters  map[string]RawParameter `json:"parameters"`
}

type RawParameter struct {
	Type     string            `json:"type"`
	Required bool              `json:"required"`
	Spec     map[string]string `json:"spec"`
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
