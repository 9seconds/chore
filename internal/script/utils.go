package script

import (
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/env"
)

func ensureTempDir(s *Script) error {
	dir, err := os.MkdirTemp(
		"",
		fmt.Sprintf("%s-%s-%s", env.ChoreDir, s.Namespace, s.Executable))
	if err != nil {
		return err
	}

	s.tmpDir = dir

	return nil
}

func readConfig(s *Script) error {
	file, err := os.Open(s.ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			// that's fine, this means that optional config is just absent
			return nil
		}

		return fmt.Errorf("cannot read script config %s: %w", s.ConfigPath(), err)
	}

	defer file.Close()

	conf, err := config.Parse(file)
	if err != nil {
		return fmt.Errorf("cannot parse config file %s: %w", s.ConfigPath(), err)
	}

	s.Config = conf

	return nil
}
