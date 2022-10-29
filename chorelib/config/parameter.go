package config

import (
	"fmt"
)

type Parameter interface {
	fmt.Stringer

	Type() string
	Required() bool
	Validate(string) error
}

type parameterBase struct {
	required bool
}

func (p parameterBase) Required() bool {
	return p.required
}

func (p parameterBase) String() string {
	if p.required {
		return "required=true"
	}
	return "required=false"
}
