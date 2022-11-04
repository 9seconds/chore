package env

import (
	"context"
	"os"
	"sync"

	"github.com/9seconds/chore/internal/argparse"
)

const idLength = 32

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

		checksum := EncodeBytes(args.Checksum())
		isolatedID := chainValues(checksum, scriptID)
		chainedIsolatedID := chainValues(isolatedID, os.Getenv(EnvIDChainIsolated))

		sendValue(ctx, results, EnvIDUnique, generateRandomString(idLength))
		sendValue(ctx, results, EnvIDIsolated, isolatedID)
		sendValue(ctx, results, EnvIDChainIsolated, chainedIsolatedID)

		if _, ok := os.LookupEnv(EnvIDChainUnique); !ok {
			sendValue(ctx, results, EnvIDChainUnique, generateRandomString(idLength))
		}
	}()
}
