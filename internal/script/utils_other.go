//go:build !(aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris)

package script

func isExecutable(path string) error {
	return nil
}
