package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/9seconds/chore/internal/filelock"
)

var lockRegexp = regexp.MustCompile(`([sx]:)?(.*)`)

type Lock struct {
	path string
	mode filelock.LockType
}

func (l Lock) LockMode() filelock.LockType {
	switch {
	case l.path == "":
		return filelock.LockTypeNo
	case l.mode == filelock.LockTypeNo:
		return filelock.LockTypeExclusive
	}

	return l.mode
}

func (l Lock) Path(defaultPath string) string {
	switch l.path {
	case "", MagicValue:
		return defaultPath
	}

	return l.path
}

func (l *Lock) UnmarshalText(b []byte) error {
	matches := lockRegexp.FindSubmatch(b)
	mode := string(matches[1])
	path := string(matches[2])

	switch mode {
	case "":
		l.mode = filelock.LockTypeNo
	case "s:":
		l.mode = filelock.LockTypeShared
	case "x:":
		l.mode = filelock.LockTypeExclusive
	default:
		return fmt.Errorf("unknown lock mode: %s", mode)
	}

	if path == MagicValue {
		l.path = MagicValue

		return nil
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("cannot get absolute path: %w", err)
	}

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("cannot stat path: %w", err)
	}

	l.path = path

	return nil
}
