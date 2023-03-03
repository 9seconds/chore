package script

import (
	"context"
	"fmt"
	"sync"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/paths"
)

type Script struct {
	Namespace  string
	Executable string

	config config.Config
	tmpDir string
}

func (s *Script) String() string {
	return s.Namespace + "/" + s.Executable
}

func (s *Script) Config() *config.Config {
	return &s.config
}

func (s *Script) Path() string {
	return paths.ConfigNamespaceScript(s.Namespace, s.Executable)
}

func (s *Script) ConfigPath() string {
	return paths.ConfigNamespaceScriptConfig(s.Namespace, s.Executable)
}

func (s *Script) DataPath() string {
	return paths.DataNamespaceScript(s.Namespace, s.Executable)
}

func (s *Script) CachePath() string {
	return paths.CacheNamespaceScript(s.Namespace, s.Executable)
}

func (s *Script) StatePath() string {
	return paths.StateNamespaceScript(s.Namespace, s.Executable)
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
		env.MakeValue(env.EnvPathTemp, s.TempPath()),
	}

	for k, v := range args.Parameters {
		environ = append(environ, env.MakeValue(env.ParameterName(k), v))
	}

	for k, v := range args.Flags {
		environ = append(environ, env.MakeValue(env.FlagName(k), string(v)))
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

func (s *Script) Init() error {
	if err := ValidateScript(s.Path()); err != nil {
		return fmt.Errorf("invalid script: %w", err)
	}

	if err := paths.EnsureRoots(s.Namespace, s.Executable); err != nil {
		return fmt.Errorf("cannot ensure script roots: %w", err)
	}

	dir, err := paths.TempDir()
	if err != nil {
		return fmt.Errorf("cannot initialize tmp dir: %w", err)
	}

	s.tmpDir = dir

	conf, err := ValidateConfig(s.ConfigPath())
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	s.config = conf

	return nil
}
