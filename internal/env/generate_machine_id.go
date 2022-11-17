package env

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/denisbrodbeck/machineid"
)

func GenerateMachineID(ctx context.Context, results chan<- string, waiters *sync.WaitGroup) {
	if _, ok := os.LookupEnv(EnvMachineID); ok {
		return
	}

	waiters.Add(1)

	go func() {
		defer waiters.Done()

		if _, ok := os.LookupEnv(EnvMachineID); ok {
			return
		}

		value, err := machineid.ProtectedID("chore")
		if err != nil {
			log.Printf("cannot obtain machine id: %v", err)

			return
		}

		sendValue(ctx, results, EnvMachineID, value)
	}()
}
