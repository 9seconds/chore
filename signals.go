//go:build !unix

package main

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
)

func makeMainContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt)
}

func relaySignals(ctx context.Context, cmd *exec.Cmd) {

}
