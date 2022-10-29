//go:build unix

package script

import "golang.org/x/sys/unix"

func isExecutable(path string) error {
	return unix.Access(path, unix.X_OK)
}
