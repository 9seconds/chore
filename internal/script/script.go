package script

import (
	"context"
	"fmt"
	"sync"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/binutils"
	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script/config"
	"github.com/gosimple/slug"
)

type Script struct {
	Namespace  string
	Executable string
	ID         string
	Config     config.Config

	tmpDir         string
	ensureDirMutex sync.Mutex
}

func (s *Script) String() string {
	return s.Namespace + "/" + s.Executable
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
		env.MakeValue(env.EnvIDRun, s.ID),
		env.MakeValue(env.EnvPathCaller, s.Path()),
		env.MakeValue(env.EnvPathData, s.DataPath()),
		env.MakeValue(env.EnvPathCache, s.CachePath()),
		env.MakeValue(env.EnvPathState, s.StatePath()),
		env.MakeValue(env.EnvPathTemp, s.TempPath()),
		env.MakeValue(
			env.EnvSlug,
			slug.Make(fmt.Sprintf(
				"%s-%s-%s-%s",
				s.Namespace, s.Executable, s.ID, args.ToSlugString()))),
	}

	for name := range args.Parameters {
		environ = append(
			environ,
			env.MakeValue(
				env.ParameterName(name),
				args.GetParameter(name)))
		environ = append(
			environ,
			env.MakeValue(
				env.ParameterNameList(name),
				args.GetParameterList(name)))
	}

	for k, v := range args.Flags {
		environ = append(environ, env.MakeValue(env.FlagName(k), string(v)))
	}

	waiterGroup := &sync.WaitGroup{}
	values := make(chan string, 1)

	env.GenerateSelf(ctx, values, waiterGroup, s.Namespace, s.Executable, args)
	env.GenerateTime(ctx, values, waiterGroup)
	env.GenerateMachineID(ctx, values, waiterGroup)
	env.GenerateIds(ctx, values, waiterGroup, s.Path(), args)
	env.GenerateOS(ctx, values, waiterGroup)
	env.GenerateUser(ctx, values, waiterGroup)
	env.GenerateHostname(ctx, values, waiterGroup)
	env.GenerateGit(ctx, values, waiterGroup, s.Config.Git)
	env.GenerateNetwork(ctx, values, waiterGroup, s.Config.Network)
	env.GenerateNetworkIPv6(ctx, values, waiterGroup, s.Config.Network)

	go func() {
		waiterGroup.Wait()
		close(values)
	}()

	for value := range values {
		environ = append(environ, value)
	}

	return environ
}

func (s *Script) EnsureDirs() error {
	s.ensureDirMutex.Lock()

	if s.tmpDir == "" {
		tmpdir, err := paths.TempDir()
		if err != nil {
			return fmt.Errorf("cannot create temp directory: %w", err)
		}

		s.tmpDir = tmpdir
	}

	s.ensureDirMutex.Unlock()

	return paths.EnsureRoots(s.Namespace, s.Executable)
}

func New(namespace, executable string) (*Script, error) {
	scr := &Script{
		Namespace:  namespace,
		Executable: executable,
		ID:         binutils.NewID(),
	}

	if err := ValidateScript(scr.Path()); err != nil {
		return nil, fmt.Errorf("invalid script: %w", err)
	}

	conf, err := ValidateConfig(scr.ConfigPath())
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	scr.Config = conf

	return scr, nil
}
