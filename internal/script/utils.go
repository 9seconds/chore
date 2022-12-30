package script

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/9seconds/chore/internal/access"
	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/env"
	"github.com/adrg/xdg"
)

const dirPermission = 0o700

func EnsureDir(path string) error {
	return os.MkdirAll(path, dirPermission)
}

func ListNamespaces(prefix string) ([]string, error) {
	dir := filepath.Join(xdg.ConfigHome, env.ChoreDir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read chore dir %s: %w", dir, err)
	}

	names := make([]string, 0, len(entries))

	for _, v := range entries {
		name := v.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		path := filepath.Join(dir, name)

		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			names = append(names, name)
		}
	}

	return names, nil
}

func ListScripts(namespace, prefix string) ([]string, error) {
	dir := filepath.Join(env.RootPathConfig(), namespace)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot list scripts in namespace %s: %w", namespace, err)
	}

	names := make([]string, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()

		if !strings.HasPrefix(name, prefix) {
			continue
		}

		scr := &Script{
			Namespace:  namespace,
			Executable: name,
		}

		if err := ValidateScript(scr.Path()); err == nil {
			names = append(names, name)
		}
	}

	return names, nil
}

func SearchScripts(namespacePrefix, scriptPrefix string) ([]*Script, error) {
	scripts := []*Script{}

	namespaces, err := ListNamespaces(namespacePrefix)
	if err != nil {
		return nil, fmt.Errorf("cannot find out namespaces: %w", err)
	}

	for _, namespace := range namespaces {
		names, err := ListScripts(namespace, scriptPrefix)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot find out scripts for %s namespace: %w",
				namespace,
				err)
		}

		for _, name := range names {
			scripts = append(scripts, &Script{
				Namespace:  namespace,
				Executable: name,
			})
		}
	}

	sort.Slice(scripts, func(i, j int) bool {
		if scripts[i].Namespace != scripts[j].Namespace {
			return scripts[i].Namespace < scripts[j].Namespace
		}

		return scripts[i].Executable < scripts[j].Executable
	})

	return scripts, nil
}

func FindScript(namespacePrefix, scriptPrefix string) (*Script, error) {
	scripts, err := SearchScripts(namespacePrefix, scriptPrefix)
	if err != nil {
		return nil, err
	}

	switch len(scripts) {
	case 0:
		return nil, errors.New("cannot find such script")
	case 1:
		return scripts[0], nil
	default:
		names := make([]string, len(scripts))

		for idx, v := range scripts {
			names[idx] = v.Namespace + "/" + v.Executable
		}

		return nil, fmt.Errorf("ambigous specification: do you mean %s?", strings.Join(names, ", "))
	}
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
		err = fmt.Errorf("cannot parse config file %script: %w", path, err)
	}

	return conf, err
}
