package env

import (
	"context"
	"os"
	"sync"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/binutils"
)

func GenerateIds(
	ctx context.Context,
	results chan<- string,
	waiters *sync.WaitGroup,
	scriptID string,
	args argparse.ParsedArgs,
) {
	waiters.Add(1)

	go func() {
		defer waiters.Done()

		chainRun := os.Getenv(IDChainRun)
		if chainRun == "" {
			chainRun = binutils.NewID()
		}

		checksum := args.Checksum()
		isolatedID := binutils.Chain(scriptID, checksum)
		chainedIsolatedID := binutils.Chain(os.Getenv(IDChainIsolated), scriptID, checksum)

		sendValue(ctx, results, IDChainRun, chainRun)
		sendValue(ctx, results, IDIsolated, isolatedID)
		sendValue(ctx, results, IDChainIsolated, chainedIsolatedID)
	}()
}
