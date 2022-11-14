//go:build !(aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos)

package access

import "os"

func Access(path string, readable, writable, executable bool) error {
	_, err := os.Stat(path)

	return err
}
