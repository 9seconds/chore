//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package script

import "golang.org/x/sys/unix"

func isExecutable(path string) error {
	return unix.Access(path, unix.X_OK)
}
