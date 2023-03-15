package edit

import (
	"errors"
	"os"
	"os/exec"
)

var ErrCannotFindOutEditor = errors.New("cannot find out editor")

type FlagEditor struct {
	Value string
}

func (f *FlagEditor) Get() (string, error) {
	if value := f.String(); value != "" {
		return value, nil
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

	return "", ErrCannotFindOutEditor
}

func (f *FlagEditor) Type() string {
	return "executable"
}

func (f *FlagEditor) String() string {
	return f.Value
}

func (f *FlagEditor) Set(value string) error {
	var err error

	if value != "" {
		value, err = exec.LookPath(value)
		if err != nil {
			return err
		}
	}

	f.Value = value

	return nil
}
