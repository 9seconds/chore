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

		sendValue(ctx, results, EnvOSType, runtime.GOOS)
		sendValue(ctx, results, EnvOSArch, runtime.GOARCH)

		if _, ok := os.LookupEnv(EnvOSID); ok {
			return
		}

		version, err := osversion.Get()
		if err != nil {
			log.Printf("cannot get os version: %v", err)

			return
		}

		sendValue(ctx, results, EnvOSID, version.ID)
		sendValue(ctx, results, EnvOSVersion, version.Version)
		sendValue(ctx, results, EnvOSCodename, version.Codename)
		sendValue(ctx, results, EnvOSVersionMajor, strconv.FormatUint(version.Major, 10))
		sendValue(ctx, results, EnvOSVersionMinor, strconv.FormatUint(version.Minor, 10))
	}()
}
