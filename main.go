package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/9seconds/chore/internal/cli"
)

var version = "dev"

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()
		cancel()
		log.Println("application context is closed")
	}()

	root := cli.NewRoot(version)

	root.InitDefaultCompletionCmd()
	root.InitDefaultHelpFlag()
	root.InitDefaultVersionFlag()
	root.InitDefaultHelpCmd()

	root.AddCommand(
		cli.NewRun(),
		cli.NewShow(),
		cli.NewEditConfig(),
		cli.NewEditScript(),
		cli.NewGC())

	if err := root.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}

	cancel()
}
