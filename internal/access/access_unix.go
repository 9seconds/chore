//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package access

import (
	"golang.org/x/sys/unix"
)

func Access(path string, readable, writable, executable bool) error {
	var flags uint32 = unix.F_OK

	if readable || writable || executable {
		flags = 0

		if readable {
			flags |= unix.R_OK
		}

		if writable {
			flags |= unix.W_OK
		}

		if executable {
			flags |= unix.X_OK
		}
	}

	return unix.Access(path, flags)
}
