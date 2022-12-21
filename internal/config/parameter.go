package config

import "context"

type Parameter interface {
	String() string
	Type() string
	Required() bool
	Validate(context.Context, string) error
}
