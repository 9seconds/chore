package env

import (
	"context"
	"os"
	"sync"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/ids"
)

const idLength = 32

func GenerateIds(
	ctx context.Context,
	results chan<- string,
	waiters *sync.WaitGroup,
	scriptID, runID string,
	args argparse.ParsedArgs,
) {
	waiters.Add(1)

	go func() {
		defer waiters.Done()

		chainRun := os.Getenv(EnvIDChainRun)
		if chainRun == "" {
			chainRun = ids.New()
		}

		checksum := ids.Encode(args.Checksum())
		isolatedID := ids.Chain(scriptID, checksum)
		chainedIsolatedID := ids.Chain(os.Getenv(EnvIDChainIsolated), scriptID, checksum)

		sendValue(ctx, results, EnvIDRun, runID)
		sendValue(ctx, results, EnvIDChainRun, chainRun)
		sendValue(ctx, results, EnvIDIsolated, isolatedID)
		sendValue(ctx, results, EnvIDChainIsolated, chainedIsolatedID)
		sendValue(ctx, results, EnvIDChainIsolated, chainedIsolatedID)
	}()
}
