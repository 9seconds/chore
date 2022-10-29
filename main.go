package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/9seconds/chore/chorelib"
	"github.com/alecthomas/kong"
)

func main() {
	cliCtx := kong.Parse(
		&CLI,
		kong.Name("chore"),
		kong.Description("Execution environment for a small helper scripts."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:   true,
			Summary:   true,
			Tree:      true,
			FlagsLast: true,
		}),
		kong.Vars{
			"version": version,
		})

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if !CLI.Debug {
		log.SetOutput(io.Discard)
	}

	appCtx, cancel := makeMainContext()
	defer cancel()

	go func() {
		<-appCtx.Done()
		log.Println("application context is closed")
	}()

	if err := os.MkdirAll(chorelib.Home, 0750); err != nil {
		log.Fatalf("cannot create home directory %s: %s", chorelib.Home, err.Error())
	}

	err := cliCtx.Run(Context{appCtx})

	if errors.Is(err, scriptExitError{}) {
		os.Exit(err.(scriptExitError).code)
	}

	cliCtx.FatalIfErrorf(err)
}
