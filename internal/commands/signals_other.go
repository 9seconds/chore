//go:build !(aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris)

package commands

import "os"

const (
	SignalInterrupt = os.Interrupt
	SignalKill      = os.Kill
)

var SignalsToRelay = []os.Signal{}
