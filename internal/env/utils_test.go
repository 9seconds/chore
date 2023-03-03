package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/assert"
)

func TestMakeValue(t *testing.T) {
	testNames := [][]string{
		{"", "value", "=value"},
		{"x", "y", "x=y"},
		{"x", "y z", "x=y z"},
		{"x", "", "x="},
	}

	for _, values := range testNames {
		name := values[0]
		value := values[1]
		expected := values[2]

		t.Run(name+"->"+value, func(t *testing.T) {
			assert.Equal(t, expected, env.MakeValue(name, value))
		})
	}
}
