package env

import (
	"context"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/9seconds/chore/internal/env/osversion"
)

func GenerateOS(ctx context.Context, results chan<- string, waiters *sync.WaitGroup) {
	waiters.Add(1)

	go func() {
		defer waiters.Done()

		sendValue(ctx, results, OSType, runtime.GOOS)
		sendValue(ctx, results, OSArch, runtime.GOARCH)

		if _, ok := os.LookupEnv(OSID); ok {
			return
		}

		version, err := osversion.Get()
		if err != nil {
			log.Printf("cannot get os version: %v", err)

			return
		}

		sendValue(ctx, results, OSID, version.ID)
		sendValue(ctx, results, OSVersion, version.Version)
		sendValue(ctx, results, OSCodename, version.Codename)
		sendValue(ctx, results, OSVersionMajor, strconv.FormatUint(version.Major, 10))
		sendValue(ctx, results, OSVersionMinor, strconv.FormatUint(version.Minor, 10))
	}()
}
