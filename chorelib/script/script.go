package script

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/9seconds/chore/chorelib"
	"github.com/9seconds/chore/chorelib/config"
	"golang.org/x/sys/unix"
)

type Script struct {
	Namespace  string
	Executable string
	Config     config.Config

	tmpDir string
}

func (s Script) String() string {
	return s.Executable
}

func (s Script) NamespacePath() string {
	return filepath.Join(chorelib.Home, s.Namespace)
}

func (s Script) Path() string {
	return filepath.Join(chorelib.Home, s.Namespace, s.Executable)
}

func (s Script) ConfigPath() string {
	return filepath.Join(chorelib.Home, s.Namespace, s.Executable+".json")
}

func (s Script) PersistentDir() string {
	return filepath.Join(chorelib.PersistentDirs, s.Namespace, s.Executable)
}

func (s Script) TempDir() string {
	return s.tmpDir
}

func (s Script) IsValid() error {
	stat, err := os.Stat(s.Path())

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return fmt.Errorf("script %s does not exist", s.Path())
	case err != nil:
		return fmt.Errorf("cannot stat script %s: %w", s.Path(), err)
	case stat.IsDir():
		return fmt.Errorf("script %s is a directory", s.Path())
	case unix.Access(s.Path(), unix.X_OK) != nil:
		return fmt.Errorf("script %s is not executable: %v", s.Path(), stat.Mode())
	}

	fp, err := os.Open(s.Path())
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", s.Path(), err)
	}

	fp.Close()

	return nil
}

func (s *Script) Init() error {
	if err := s.ensureTmpDir(); err != nil {
		return err
	}

	if err := s.ensurePersistentDir(); err != nil {
		return err
	}

	if err := s.readConfig(); err != nil {
		return err
	}

	return nil
}

func (s *Script) ensureTmpDir() error {
	tmpDir, err := os.MkdirTemp(
		"",
		fmt.Sprintf("chore-%s-%s-", s.Namespace, s.Executable))
	if err != nil {
		return fmt.Errorf("cannot create tmp dir: %s", err)
	}

	s.tmpDir = tmpDir

	return nil
}

func (s *Script) ensurePersistentDir() error {
	if err := os.MkdirAll(s.PersistentDir(), 0750); err != nil {
		return fmt.Errorf(
			"cannot ensure presence of persistent dir %s: %w",
			s.PersistentDir(),
			err)
	}

	return nil
}

func (s *Script) readConfig() error {
	conf := config.Config{}
	file, err := os.Open(s.ConfigPath())

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil
	case err != nil:
		return fmt.Errorf(
			"cannot read config file %s: %w",
			s.ConfigPath(),
			err)
	}

	defer file.Close()

	conf, err = config.Parse(file)
	if err != nil {
		return fmt.Errorf("cannot parse config: %w", err)
	}

	s.Config = conf

	return nil
}

func (s *Script) Cleanup() error {
	if s.tmpDir != "" {
		err := os.RemoveAll(s.tmpDir)
		s.tmpDir = ""

		return err
	}

	return nil
}
