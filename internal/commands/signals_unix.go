//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package commands

import (
	"os"
	"syscall"
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
	// syscall.SIGCLD,
	syscall.SIGPWR,
	syscall.SIGWINCH,
	syscall.SIGUNUSED,
}
