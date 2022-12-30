package script

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/access"
	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/env"
)

type Script struct {
	Namespace  string
	Executable string

	config    config.Config
	closeOnce sync.Once
	tmpDir    string
}

func (s *Script) String() string {
	return s.Namespace + "/" + s.Executable
}

func (s *Script) Config() *config.Config {
	return &s.config
}

func (s *Script) NamespacePath() string {
	return filepath.Join(env.RootPathConfig(), s.Namespace)
}

func (s *Script) Path() string {
	return filepath.Join(s.NamespacePath(), s.Executable)
}

func (s *Script) ConfigPath() string {
	return s.Path() + ".hjson"
}

func (s *Script) DataPath() string {
	return filepath.Join(env.RootPathData(), s.Namespace, s.Executable)
}

func (s *Script) CachePath() string {
	return filepath.Join(env.RootPathCache(), s.Namespace, s.Executable)
}

func (s *Script) StatePath() string {
	return filepath.Join(env.RootPathState(), s.Namespace, s.Executable)
}

func (s *Script) RuntimePath() string {
	return filepath.Join(env.RootPathRuntime(), s.Namespace, s.Executable)
}

func (s *Script) TempPath() string {
	return s.tmpDir
}

func (s *Script) Environ(ctx context.Context, args argparse.ParsedArgs) []string {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	environ := []string{
		env.MakeValue(env.EnvNamespace, s.Namespace),
		env.MakeValue(env.EnvCaller, s.Executable),
		env.MakeValue(env.EnvPathCaller, s.Path()),
		env.MakeValue(env.EnvPathData, s.DataPath()),
		env.MakeValue(env.EnvPathCache, s.CachePath()),
		env.MakeValue(env.EnvPathState, s.StatePath()),
		env.MakeValue(env.EnvPathRuntime, s.RuntimePath()),
		env.MakeValue(env.EnvPathTemp, s.TempPath()),
	}

	for k, v := range args.Parameters {
		environ = append(
			environ,
			env.MakeValue(env.EnvParameterPrefix+strings.ToUpper(k), v))
	}

	for k, v := range args.Flags {
		environ = append(
			environ,
			env.MakeValue(env.EnvFlagPrefix+strings.ToUpper(k), string(v)))
	}

	waiterGroup := &sync.WaitGroup{}
	values := make(chan string, 1)

	env.GenerateRecursion(ctx, values, waiterGroup, s.Namespace, s.Executable, args)
	env.GenerateTime(ctx, values, waiterGroup)
	env.GenerateMachineID(ctx, values, waiterGroup)
	env.GenerateIds(ctx, values, waiterGroup, s.Path(), args)
	env.GenerateOS(ctx, values, waiterGroup)
	env.GenerateUser(ctx, values, waiterGroup)
	env.GenerateHostname(ctx, values, waiterGroup)
	env.GenerateGit(ctx, values, waiterGroup, s.config.Git)
	env.GenerateNetwork(ctx, values, waiterGroup, s.config.Network)
	env.GenerateNetworkIPv6(ctx, values, waiterGroup, s.config.Network)

	go func() {
		waiterGroup.Wait()
		close(values)
	}()

	for value := range values {
		environ = append(environ, value)
	}

	return environ
}

func (s *Script) Valid() error {
	path := s.Path()

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat path: %w", err)
	}

	if stat.IsDir() {
		return fmt.Errorf("path is directory: %w", err)
	}

	if stat.Size() == 0 {
		return errors.New("script is empty")
	}

	if err := access.Access(path, false, false, true); err != nil {
		return fmt.Errorf("cannot find out executable %s: %w", s.Path(), err)
	}

	return nil
}

func (s *Script) Init() error {
	if err := s.Valid(); err != nil {
		return fmt.Errorf("cannot find out executable %s: %w", s.Path(), err)
	}

	if err := EnsureDir(s.DataPath()); err != nil {
		return fmt.Errorf("cannot create data path %s: %w", s.DataPath(), err)
	}

	if err := EnsureDir(s.CachePath()); err != nil {
		return fmt.Errorf("cannot create cache path %s: %w", s.CachePath(), err)
	}

	if err := EnsureDir(s.StatePath()); err != nil {
		return fmt.Errorf("cannot create state path %s: %w", s.StatePath(), err)
	}

	if err := EnsureDir(s.RuntimePath()); err != nil {
		return fmt.Errorf("cannot create runtime path %s: %w", s.RuntimePath(), err)
	}

	dir, err := os.MkdirTemp("", env.ChoreDir+"-")
	if err != nil {
		return fmt.Errorf("cannot initialize tmp dir: %w", err)
	}

	s.tmpDir = dir

	file, err := os.Open(s.ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			// that'script fine, this means that optional config is just absent
			return nil
		}

		return fmt.Errorf("cannot read script config %script: %w", s.ConfigPath(), err)
	}

	defer file.Close()

	conf, err := config.Parse(file)
	if err != nil {
		return fmt.Errorf("cannot parse config file %script: %w", s.ConfigPath(), err)
	}

	s.config = conf

	return nil
}

func (s *Script) Cleanup() {
	s.closeOnce.Do(func() {
		if s.tmpDir != "" {
			os.RemoveAll(s.tmpDir)
		}

		s.tmpDir = ""
	})
}
