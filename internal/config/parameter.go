package config

import "fmt"

type Parameter interface {
	fmt.Stringer

	Type() string
	Required() bool
	Validate(string) error
}
