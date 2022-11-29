package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

type Lock string

func (l Lock) Value(defaultPath string) string {
	if l == MagicValue {
		return defaultPath
	}

	return string(l)
}

func (l *Lock) UnmarshalText(b []byte) error {
	text := string(b)

	if text == MagicValue {
		*l = Lock(MagicValue)

		return nil
	}

	text, err := filepath.Abs(text)
	if err != nil {
		return fmt.Errorf("cannot get absolute path: %w", err)
	}

	if _, err := os.Stat(text); err != nil {
		return fmt.Errorf("cannot stat path: %w", err)
	}

	*l = Lock(text)

	return nil
}
