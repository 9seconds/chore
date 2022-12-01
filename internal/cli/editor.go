package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var errCannotFindOutEditor = errors.New("cannot find out editor")

type Editor string

func (e Editor) Value() (string, error) {
	if e != "" {
		return string(e), nil
	}

	if value := os.Getenv("VISUAL"); value != "" {
		return value, nil
	}

	if value := os.Getenv("EDITOR"); value != "" {
		return value, nil
	}

	if path, err := exec.LookPath("sensible-editor"); err == nil {
		return path, nil
	}

	if path, err := exec.LookPath("editor"); err == nil {
		return path, nil
	}

	if path, err := exec.LookPath("nano"); err == nil {
		return path, nil
	}

	return "", errCannotFindOutEditor
}

func (e *Editor) UnmarshalText(b []byte) error {
	path, err := exec.LookPath(string(b))
	if err != nil {
		return fmt.Errorf("cannot detect given editor: %w", err)
	}

	*e = Editor(path)

	return nil
}
