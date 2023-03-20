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

		sendValue(ctx, results, Bin, executable)

		cli := []string{executable, "run", namespace, script}
		cli = append(cli, args.ToSelfStringChunks()...)

		sendValue(ctx, results, Self, shellescape.QuoteCommand(cli))
	}()
}
