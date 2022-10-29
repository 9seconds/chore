package config

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

var spaceRegexp = regexp.MustCompile(`\s`)

const (
	ParameterString  = "string"
	ParameterInteger = "integer"
)

type Config struct {
	Description string
	Network     bool
	Parameters  map[string]Parameter
}

type RawConfig struct {
	Description string                  `json:"description"`
	Network     bool                    `json:"network"`
	Parameters  map[string]RawParameter `json:"parameters"`
}

type RawParameter struct {
	Type     string            `json:"type"`
	Required bool              `json:"required"`
	Spec     map[string]string `json:"spec"`
}

func Parse(reader io.Reader) (Config, error) {
	raw := RawConfig{}
	conf := Config{}

	if err := json.NewDecoder(reader).Decode(&raw); err != nil {
		return conf, fmt.Errorf("cannot parse JSON config: %w", err)
	}

	conf.Description = raw.Description
	conf.Network = raw.Network
	conf.Parameters = make(map[string]Parameter)

	for name, param := range raw.Parameters {
		if found := spaceRegexp.FindStringIndex(name); found != nil {
			return conf, fmt.Errorf("incorrect parameter name %s", name)
		}

		var value Parameter
		var err error

		switch param.Type {
		case ParameterInteger:
			value, err = NewInteger(param.Required, param.Spec)
		case ParameterString:
			value, err = NewString(param.Required, param.Spec)
		default:
			return conf, fmt.Errorf("unknown parameter type %s for parameter %s", param.Type, name)
		}

		if err != nil {
			return conf, fmt.Errorf("cannot initialize parameter %s: %w", name, err)
		}

		conf.Parameters[name] = value
	}

	return conf, nil
}
