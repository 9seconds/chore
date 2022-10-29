//go:build !unix

package script

func isExecutable(path string) error {
	return nil
}
