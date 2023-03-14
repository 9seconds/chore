package config_test

import (
	"regexp"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	allPasswords := make(map[string]bool)
	digitsRegexp := regexp.MustCompile(`\d`)
	symbolsRegexp := regexp.MustCompile(`[^0-9A-Za-z]`)

	for i := 0; i < 5000; i++ {
		pass := config.GeneratePassword()

		assert.Len(t, pass, config.PasswordLength)
		assert.Len(t, digitsRegexp.FindAllString(pass, -1), config.PasswordNumDigits)
		assert.Len(t, symbolsRegexp.FindAllString(pass, -1), config.PasswordNumSymbols)

		allPasswords[pass] = true
	}

	assert.Len(t, allPasswords, 5000)
}
