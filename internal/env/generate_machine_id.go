package env

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/denisbrodbeck/machineid"
)

func GenerateMachineId(ctx context.Context, results chan<- string, wg *sync.WaitGroup) {
	if _, ok := os.LookupEnv(EnvMachineId); ok {
		return
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		value, err := machineid.ProtectedID("chore")
		if err != nil {
			log.Printf("cannot obtain machine id: %v", err)

			return
		}

		sendValue(ctx, results, EnvMachineId, value)
	}()
}
