package env_test

import (
	"strings"
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

func TestEnviron(t *testing.T) {
	t.Setenv(env.ParameterName("X"), "1")
	t.Setenv(env.FlagName("Y"), "1")
	t.Setenv(env.ParameterNameList("Z"), "1")

	for _, value := range env.Environ() {
		assert.False(t, strings.HasPrefix(value, env.EnvFlagPrefix))
		assert.False(t, strings.HasPrefix(value, env.EnvParameterPrefix))
		assert.False(t, strings.HasPrefix(value, env.EnvParameterPrefixList))
	}
}
