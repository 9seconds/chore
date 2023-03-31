package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/9seconds/chore/internal/cli"
	"github.com/gosimple/slug"
)

var version = "dev"

func main() {
	slug.Lowercase = false
	slug.MaxLength = 100

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
		cli.NewEdit(),
		cli.NewRun(),
		cli.NewRemove(),
		cli.NewRename(),
		cli.NewShow(),
		cli.NewVault(),
		cli.NewGC(),
		cli.NewUpdate())

	root.SetIn(os.Stdin)
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	if err := root.ExecuteContext(ctx); err != nil {
		cancel()
		log.Fatal(err)
	}

	cancel()
}
