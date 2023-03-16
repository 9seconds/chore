package config

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/paths"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Env   map[string]map[string]string `toml:"env"`
	Vault map[string]string            `toml:"vault"`
}

func (c Config) Environ(namespace string) []string {
	values := c.Env[namespace]
	chunks := make([]string, 0, len(values))

	for key, value := range values {
		chunks = append(chunks, env.MakeValue(key, value))
	}

	return chunks
}

func ReadConfig(reader io.Reader) (Config, error) {
	decoder := toml.NewDecoder(reader)
	conf := Config{}

	if _, err := decoder.Decode(&conf); err != nil {
		return conf, fmt.Errorf("cannot parse TOML config: %w", err)
	}

	return conf, nil
}

func Get() (Config, error) {
	reader, err := os.Open(paths.AppConfigPath())

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return Config{}, nil
	case err != nil:
		return Config{}, fmt.Errorf("cannot open a path: %w", err)
	}

	defer reader.Close()

	return ReadConfig(reader)
}
