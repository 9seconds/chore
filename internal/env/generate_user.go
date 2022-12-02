package env

import (
	"context"
	"log"
	"os"
	"os/user"
	"sync"
)

func GenerateUser(ctx context.Context, results chan<- string, waiters *sync.WaitGroup) {
	if _, ok := os.LookupEnv(EnvUserName); ok {
		return
	}

	waiters.Add(1)

	go func() {
		defer waiters.Done()

		user, err := user.Current()
		if err != nil {
			log.Printf("cannot get current user: %v", err)

			return
		}

		sendValue(ctx, results, EnvUserUID, user.Uid)
		sendValue(ctx, results, EnvUserGID, user.Gid)
		sendValue(ctx, results, EnvUserName, user.Username)
	}()
}
