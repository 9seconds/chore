package config

import (
	"github.com/sethvargo/go-password/password"
)

const (
	PasswordLength     = 16
	PasswordNumDigits  = 3
	PasswordNumSymbols = 3
)

func GeneratePassword() string {
	return password.MustGenerate(PasswordLength, PasswordNumDigits, PasswordNumSymbols, false, true)
}
