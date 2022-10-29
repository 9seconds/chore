package env

import (
	"context"
)

const runIDLength = 32

func generateRunId(ctx context.Context, result chan<- string) {
	sendEnvValue(
		ctx,
		result,
		EnvRunId,
		generateRandomString(runIDLength))
}
