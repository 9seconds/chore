package script

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/access"
	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/env"
	"github.com/adrg/xdg"
)

const defaultDirPermission = 0o700

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

func (s Script) Environ(ctx context.Context, args argparse.ParsedArgs) []string {
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

	for k, v := range args.Keywords {
		environ = append(
			environ,
			env.MakeValue(env.EnvArgPrefix+strings.ToUpper(k), v))
	}

	waiterGroup := &sync.WaitGroup{}
	values := make(chan string, 1)

	env.GenerateTime(ctx, values, waiterGroup)
	env.GenerateMachineID(ctx, values, waiterGroup)
	env.GenerateIds(ctx, values, waiterGroup, s.Path(), args)
	env.GenerateOS(ctx, values, waiterGroup)

	if s.Config.Network {
		env.GenerateNetwork(ctx, values, waiterGroup)
		env.GenerateNetworkIPv6(ctx, values, waiterGroup)
	}

	go func() {
		waiterGroup.Wait()
		close(values)
	}()

	for value := range values {
		environ = append(environ, value)
	}

	return environ
}

func New(namespace, executable string) (Script, error) {
	rValue := Script{
		Namespace:  namespace,
		Executable: executable,
	}

	if err := access.Access(rValue.Path(), false, false, true); err != nil {
		return rValue, fmt.Errorf("cannot find out executable %s: %w", rValue.Path(), err)
	}

	if err := os.MkdirAll(rValue.DataPath(), defaultDirPermission); err != nil {
		return rValue, fmt.Errorf("cannot create data path %s: %w", rValue.DataPath(), err)
	}

	if err := os.MkdirAll(rValue.CachePath(), defaultDirPermission); err != nil {
		return rValue, fmt.Errorf("cannot create cache path %s: %w", rValue.CachePath(), err)
	}

	if err := os.MkdirAll(rValue.StatePath(), defaultDirPermission); err != nil {
		return rValue, fmt.Errorf("cannot create state path %s: %w", rValue.StatePath(), err)
	}

	if err := os.MkdirAll(rValue.RuntimePath(), defaultDirPermission); err != nil {
		return rValue, fmt.Errorf("cannot create runtime path %s: %w", rValue.RuntimePath(), err)
	}

	if err := readConfig(&rValue); err != nil {
		return rValue, err
	}

	if err := ensureTempDir(&rValue); err != nil {
		return rValue, fmt.Errorf("cannot create temporary dir: %w", err)
	}

	return rValue, nil
}
