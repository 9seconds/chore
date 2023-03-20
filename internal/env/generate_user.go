package env

import (
	"context"
	"log"
	"os"
	"os/user"
	"sync"
)

func GenerateUser(ctx context.Context, results chan<- string, waiters *sync.WaitGroup) {
	if _, ok := os.LookupEnv(UserName); ok {
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

		sendValue(ctx, results, UserUID, user.Uid)
		sendValue(ctx, results, UserGID, user.Gid)
		sendValue(ctx, results, UserName, user.Username)
	}()
}
