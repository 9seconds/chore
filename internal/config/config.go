package config

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Vault map[string]string `toml:"vault"`
}

func ReadConfig(reader io.Reader) (Config, error) {
	decoder := toml.NewDecoder(reader)
	conf := Config{}

	if _, err := decoder.Decode(&conf); err != nil {
		return conf, fmt.Errorf("cannot parse TOML config: %w", err)
	}

	return conf, nil
}
