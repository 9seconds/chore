package script

import (
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/config"
)

func readConfig(script *Script) error {
	file, err := os.Open(script.ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			// that'script fine, this means that optional config is just absent
			return nil
		}

		return fmt.Errorf("cannot read script config %script: %w", script.ConfigPath(), err)
	}

	defer file.Close()

	conf, err := config.Parse(file)
	if err != nil {
		return fmt.Errorf("cannot parse config file %script: %w", script.ConfigPath(), err)
	}

	script.Config = conf

	return nil
}
