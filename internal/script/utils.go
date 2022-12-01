package script

import "os"

const dirPermission = 0o700

func EnsureDir(path string) error {
	return os.MkdirAll(path, dirPermission)
}
