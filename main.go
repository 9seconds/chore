package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/commands"
	"github.com/gosimple/slug"
)

func main() {
	defer commands.Exit(0)

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

	root := cli.NewRoot(getVersion())

	root.InitDefaultCompletionCmd()
	root.InitDefaultHelpFlag()
	root.InitDefaultVersionFlag()
	root.InitDefaultHelpCmd()

	root.AddCommand(
		cli.NewEditConfig(),
		cli.NewRun(),
		cli.NewShow(),
		cli.NewEditScriptConfig(),
		cli.NewEditScript(),
		cli.NewGC())

	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		commands.Exit(1)
	}

	cancel()
}

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("cannot read build info")
	}

	commit := ""
	date := ""
	isDirty := ""

	for _, setting := range info.Settings {
		switch {
		case setting.Key == "vcs.revision":
			commit = setting.Value
		case setting.Key == "vcs.time":
			date = setting.Value
		case setting.Key == "vcs.modified" && setting.Value == "true":
			isDirty = "[!] "
		}
	}

	if commit == "" {
		return "dev"
	}

	return fmt.Sprintf("%s%s (%s, %s)", isDirty, commit, date, info.GoVersion)
}
