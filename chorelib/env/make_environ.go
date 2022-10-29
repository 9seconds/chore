package env

import (
	"context"
	"fmt"
	"sync"

	"github.com/9seconds/chore/chorelib/script"
)

func MakeEnviron(ctx context.Context, script script.Script, args []string) ([]string, error) {
	parsedArgs, err := parseArgs(script.Config.Parameters, args)
	if err != nil {
		return nil, fmt.Errorf("cannot parse script arguments: %w", err)
	}

	environ := []string{
		makeEnvValue(EnvNamespace, script.Namespace),
		makeEnvValue(EnvCaller, script.Executable),
		makeEnvValue(EnvCallerPath, script.Path()),
		makeEnvValue(EnvPersistentDir, script.PersistentDir()),
		makeEnvValue(EnvTempDir, script.TempDir()),
	}

	valueStream := make(chan string, 1)
	wg := &sync.WaitGroup{}

	wg.Add(5)

	go func() {
		wg.Wait()
		close(valueStream)
	}()

	go func() {
		generateTime(ctx, valueStream)
		wg.Done()
	}()

	go func() {
		generateRunId(ctx, valueStream)
		wg.Done()
	}()

	go func() {
		generateCorrelateId(ctx, valueStream)
		wg.Done()
	}()

	go func() {
		generateMachineId(ctx, valueStream)
		wg.Done()
	}()

	go func() {
		generateArgs(ctx, valueStream, script, parsedArgs)
		wg.Done()
	}()

	if script.Config.Network {
		wg.Add(2)

		go func() {
			generateNetworkFromIPInfo(ctx, valueStream)
			wg.Done()
		}()

		go func() {
			generateNetworkIPv6(ctx, valueStream)
			wg.Done()
		}()
	}

	for value := range valueStream {
		environ = append(environ, value)
	}

	return environ, nil
}
