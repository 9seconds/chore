package config

import "fmt"

type Flag interface {
	Required() bool
	Description() string
	String() string
}

type baseFlag struct {
	required    bool
	description string
}

func (b baseFlag) Required() bool {
	return b.required
}

func (b baseFlag) Description() string {
	return b.description
}

func (b baseFlag) String() string {
	return fmt.Sprintf("%q (required=%t)", b.description, b.required)
}

func NewFlag(description string, required bool) Flag {
	return baseFlag{
		required:    required,
		description: description,
	}
}
