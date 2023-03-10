package script

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"unicode"

	"github.com/9seconds/chore/internal/access"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script/config"
)

const dirPermission = 0o700

func EnsureDir(path string) error {
	return os.MkdirAll(path, dirPermission)
}

func ListNamespaces() ([]string, error) {
	dir := paths.ConfigRoot()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read chore dir %s: %w", dir, err)
	}

	names := make([]string, 0, len(entries))

	for _, v := range entries {
		name := v.Name()

		if stat, err := os.Stat(filepath.Join(dir, name)); err == nil && stat.IsDir() {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names, nil
}

func ListScripts(namespace string) ([]string, error) {
	entries, err := os.ReadDir(paths.ConfigNamespace(namespace))
	if err != nil {
		return nil, fmt.Errorf("cannot list scripts in namespace %s: %w", namespace, err)
	}

	names := make([]string, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()

		scr := &Script{
			Namespace:  namespace,
			Executable: name,
		}

		if err := ValidateScript(scr.Path()); err == nil {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names, nil
}

func ValidateScript(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat path: %w", err)
	}

	if stat.IsDir() {
		return fmt.Errorf("path is directory: %w", err)
	}

	if err := access.Access(path, false, false, true); err != nil {
		return fmt.Errorf("cannot find out executable %s: %w", path, err)
	}

	file, _ := os.Open(path)
	reader := bufio.NewReader(file)

	defer file.Close()

	for {
		char, _, err := reader.ReadRune()

		switch {
		case errors.Is(err, io.EOF):
			return errors.New("script is empty")
		case err != nil:
			return fmt.Errorf("cannot scan script: %w", err)
		case !unicode.IsSpace(char):
			return nil
		}
	}
}

func ValidateConfig(path string) (config.Config, error) {
	conf := config.Config{}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// that'script fine, this means that optional config is just absent
			return conf, nil
		}

		return conf, fmt.Errorf("cannot read script config %script: %w", path, err)
	}

	defer file.Close()

	conf, err = config.Parse(file)
	if err != nil {
		err = fmt.Errorf("cannot parse config file %s: %w", path, err)
	}

	return conf, err
}

func SearchScripts(_, _ string) ([]*Script, error) {
	return nil, nil
}
