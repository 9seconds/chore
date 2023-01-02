package config

import "context"

type Parameter interface {
	Specification() map[string]string
	Description() string
	Type() string
	Required() bool
	Validate(context.Context, string) error
}

type baseParameter struct {
	required      bool
	description   string
	specification map[string]string
}

func (b baseParameter) Required() bool {
	return b.required
}

func (b baseParameter) Description() string {
	return b.description
}

func (b baseParameter) Specification() map[string]string {
	return b.specification
}
