package env

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/Showmax/go-fqdn"
)

func GenerateHostname(ctx context.Context, results chan<- string, waiters *sync.WaitGroup) {
	waiters.Add(1)

	go func() {
		defer waiters.Done()

		if _, ok := os.LookupEnv(EnvHostname); !ok {
			if value, err := os.Hostname(); err == nil {
				sendValue(ctx, results, EnvHostname, value)
			} else {
				log.Printf("cannot get hostname: %v", err)
			}
		}

		if _, ok := os.LookupEnv(EnvHostnameFQDN); !ok {
			if value, err := fqdn.FqdnHostname(); err == nil {
				sendValue(ctx, results, EnvHostnameFQDN, value)
			} else {
				log.Printf("cannot get fqdn hostname: %v", err)
			}
		}
	}()
}
