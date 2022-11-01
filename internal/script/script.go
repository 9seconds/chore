package script

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/env"
	"github.com/adrg/xdg"
)

type Script struct {
	Namespace  string
	Executable string
	Config     config.Config

	tmpDir string
}

func (s Script) String() string {
	return s.Namespace + "/" + s.Executable
}

func (s Script) buildPath(base string) string {
	return filepath.Join(base, env.ChoreDir, s.Namespace, s.Executable)
}

func (s Script) Path() string {
	return s.buildPath(xdg.ConfigHome)
}

func (s Script) ConfigPath() string {
	return s.Path() + ".json"
}

func (s Script) DataPath() string {
	return s.buildPath(xdg.DataHome)
}

func (s Script) CachePath() string {
	return s.buildPath(xdg.CacheHome)
}

func (s Script) StatePath() string {
	return s.buildPath(xdg.StateHome)
}

func (s Script) RuntimePath() string {
	return s.buildPath(xdg.RuntimeDir)
}

func (s Script) TempPath() string {
	return s.tmpDir
}

func New(namespace, executable string) (Script, error) {
	rv := Script{
		Namespace:  namespace,
		Executable: executable,
	}

	if err := isExecutable(rv.Path()); err != nil {
		return rv, fmt.Errorf("cannot find out executable %s: %w", rv.Path(), err)
	}

	if err := os.MkdirAll(rv.DataPath(), 0700); err != nil {
		return rv, fmt.Errorf("cannot create data path %s: %w", rv.DataPath(), err)
	}

	if err := os.MkdirAll(rv.CachePath(), 0700); err != nil {
		return rv, fmt.Errorf("cannot create cache path %s: %w", rv.CachePath(), err)
	}

	if err := os.MkdirAll(rv.StatePath(), 0700); err != nil {
		return rv, fmt.Errorf("cannot create state path %s: %w", rv.StatePath(), err)
	}

	if err := os.MkdirAll(rv.RuntimePath(), 0700); err != nil {
		return rv, fmt.Errorf("cannot create runtime path %s: %w", rv.RuntimePath(), err)
	}

	if err := readConfig(&rv); err != nil {
		return rv, err
	}

	if err := ensureTempDir(&rv); err != nil {
		return rv, fmt.Errorf("cannot create temporary dir: %w", err)
	}

	return rv, nil
}
