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
	wg *sync.WaitGroup,
	scriptId string,
	args argparse.ParsedArgs) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		checksum := EncodeBytes(args.Checksum())
		isolatedId := chainValues(checksum, scriptId)
		chainedIsolatedId := chainValues(isolatedId, os.Getenv(EnvIdChainIsolated))

		sendValue(ctx, results, EnvIdUnique, generateRandomString(idLength))
		sendValue(ctx, results, EnvIdIsolated, isolatedId)
		sendValue(ctx, results, EnvIdChainIsolated, chainedIsolatedId)

		if _, ok := os.LookupEnv(EnvIdChainUnique); !ok {
			sendValue(ctx, results, EnvIdChainUnique, generateRandomString(idLength))
		}
	}()
}
