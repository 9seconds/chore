//go:build darwin

package commands

import (
	"os"
	"syscall"
)

const (
	SignalInterrupt = syscall.SIGTERM
	SignalKill      = syscall.SIGKILL
)

var SignalsToRelay = []os.Signal{
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
	syscall.SIGPROF,
	syscall.SIGSYS,
	syscall.SIGTRAP,
	syscall.SIGURG,
	syscall.SIGVTALRM,
	syscall.SIGXCPU,
	syscall.SIGXFSZ,
	syscall.SIGIOT,
	syscall.SIGIO,
	syscall.SIGWINCH,
}
