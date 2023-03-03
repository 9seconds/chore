package env

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/alessio/shellescape"
)

func GenerateSelf(
	ctx context.Context,
	results chan<- string,
	waiters *sync.WaitGroup,
	namespace, script string,
	args argparse.ParsedArgs,
) {
	waiters.Add(1)

	go func() {
		defer waiters.Done()

		executable, err := os.Executable()
		if err != nil {
			log.Printf("cannot find out current executable: %v", err)

			return
		}

		cli := []string{executable, "run", namespace, script}

		for name, value := range args.Parameters {
			cli = append(cli, name+string(argparse.SeparatorKeyword)+value)
		}

		for name, value := range args.Flags {
			if value == argparse.FlagTrue {
				cli = append(cli, string(argparse.PrefixFlagPositive)+name)
			} else {
				cli = append(cli, string(argparse.PrefixFlagNegative)+name)
			}
		}

		sendValue(ctx, results, EnvSelf, shellescape.QuoteCommand(cli))
	}()
}
