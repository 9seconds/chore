package commands

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	StopGracefulPeriod = 5 * time.Second
	CheckProcessEvery  = 50 * time.Millisecond
)

func osSidecarSignals(ctx context.Context, waiters *sync.WaitGroup, cmd *exec.Cmd) {
	defer waiters.Done()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)

	defer signal.Stop(sigChan)

	signal.Notify(sigChan, SignalsToRelay...)

	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigChan:
			if !osIsProcessAlive(cmd.Process) {
				return
			}

			log.Printf("send %s to %d", sig, cmd.Process.Pid)

			if err := osSendSignal(cmd.Process, sig); err != nil {
				log.Printf("cannot send %v to process %d: %v", sig, cmd.Process.Pid, err)
			}
		}
	}
}

func osSidecarGracefulShutdown(ctx context.Context, waiters *sync.WaitGroup, cmd *exec.Cmd) {
	defer waiters.Done()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	<-ctx.Done()

	if !osIsProcessAlive(cmd.Process) {
		return
	}

	log.Printf("timeout. send interrupt signal")

	ticker := time.NewTicker(CheckProcessEvery)
	defer ticker.Stop()

	gracefulTimer := time.NewTimer(StopGracefulPeriod)
	defer gracefulTimer.Stop()

	if err := osSendSignal(cmd.Process, SignalInterrupt); err != nil {
		log.Printf("cannot send %v to process %d: %v", SignalInterrupt, cmd.Process.Pid, err)

		return
	}

	for {
		select {
		case <-ticker.C:
			if !osIsProcessAlive(cmd.Process) {
				return
			}
		case <-gracefulTimer.C:
			if !osIsProcessAlive(cmd.Process) {
				return
			}

			log.Printf("graceful period is over. send kill signal")

			if err := osSendSignal(cmd.Process, SignalKill); err != nil {
				log.Printf("cannot send %v to process %d: %v", SignalInterrupt, cmd.Process.Pid, err)
			}

			return
		}
	}
}

func osIsProcessAlive(proc *os.Process) bool {
	err := proc.Signal(syscall.Signal(0))

	if errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH) {
		return false
	}

	if err != nil {
		log.Printf("cannot correctly detect if process is still alive: %v", err)
	}

	return true
}

func osSendSignal(proc *os.Process, sig os.Signal) error {
	if err := proc.Signal(sig); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}

	return nil
}
