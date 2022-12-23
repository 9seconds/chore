package env

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/alessio/shellescape"
)

func GenerateRecursion(
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

		cli := append(
			[]string{executable, "run", namespace, script},
			args.Options()...)

		sendValue(ctx, results, EnvRecursion, shellescape.QuoteCommand(cli))
	}()
}
