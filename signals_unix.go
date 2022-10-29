//go:build unix

package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func makeMainContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
}

func relaySignals(ctx context.Context, cmd *exec.Cmd) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		syscall.SIGABRT,
		syscall.SIGALRM,
		syscall.SIGFPE,
		syscall.SIGHUP,
		syscall.SIGILL,
		syscall.SIGPIPE,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGCONT,
		syscall.SIGSTOP,
		syscall.SIGTSTP,
		syscall.SIGTTIN,
		syscall.SIGTTOU,
		syscall.SIGBUS,
		syscall.SIGPOLL,
		syscall.SIGPROF,
		syscall.SIGSYS,
		syscall.SIGTRAP,
		syscall.SIGURG,
		syscall.SIGVTALRM,
		syscall.SIGXCPU,
		syscall.SIGXFSZ,
		syscall.SIGIOT,
		syscall.SIGSTKFLT,
		syscall.SIGIO,
		syscall.SIGCLD,
		syscall.SIGPWR,
		syscall.SIGWINCH,
		syscall.SIGUNUSED)

	go func() {
		defer close(sigChan)
		defer signal.Stop(sigChan)

		for {
			select {
			case <-ctx.Done():
				return
			case sig := <-sigChan:
				if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
					return
				}

				log.Printf("send %v to %d", sig, cmd.Process.Pid)

				if err := cmd.Process.Signal(sig); err != nil {
					log.Printf("cannot send %v to process %d: %v", sig, cmd.Process.Pid, err)
				}
			}
		}
	}()
}
