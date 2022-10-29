package env

import (
	"context"
	"log"
	"os"

	"github.com/denisbrodbeck/machineid"
)

func generateMachineId(ctx context.Context, result chan<- string) {
	if _, ok := os.LookupEnv(EnvMachineId); !ok {
		value, err := machineid.ProtectedID("chore")
		if err != nil {
			log.Printf("cannot obtain machine id: %v", err)

			return
		}

		sendEnvValue(ctx, result, EnvMachineId, value)
	}
}
