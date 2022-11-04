package config

import (
	"context"
	"fmt"
)

type Parameter interface {
	fmt.Stringer

	Type() string
	Required() bool
	Validate(context.Context, string) error
}
