package env

import (
	"context"
	"os"
)

const correlateIdLength = 32

func generateCorrelateId(ctx context.Context, result chan<- string) {
	if _, ok := os.LookupEnv(EnvCorrelateId); !ok {
		sendEnvValue(
			ctx,
			result,
			EnvCorrelateId,
			generateRandomString(correlateIdLength))
	}
}
