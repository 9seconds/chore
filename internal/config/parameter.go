package config

import "context"

type Parameter interface {
	String() string
	Description() string
	Type() string
	Required() bool
	Validate(context.Context, string) error
}

type baseParameter struct {
	required    bool
	description string
}

func (b baseParameter) Required() bool {
	return b.required
}

func (b baseParameter) Description() string {
	return b.description
}
