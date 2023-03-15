package argparse

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/anmitsu/go-shlex"
)

const (
	PrefixFlagPositive = '+'
	PrefixFlagNegative = '-'
	PrefixLiteral      = ':'
	SeparatorKeyword   = '='

	FlagTrue  = "t"
	FlagFalse = "f"

	PositionalDelimiter = "--"
)

func Parse(args []string) (ParsedArgs, error) { //nolint: cyclop
	parsed := ParsedArgs{
		Parameters: make(map[string][]string),
		Flags:      make(map[string]string),
		Positional: []string{},
	}

	for idx, arg := range args {
		if !utf8.ValidString(arg) {
			return parsed, fmt.Errorf("argument %d is not valid UTF-8 string", idx+1)
		}

		rune, _ := utf8.DecodeRuneInString(arg)

		switch {
		case arg == PositionalDelimiter && !parsed.IsPositionalTime():
			parsed.ExplicitPositional = true
		case parsed.ExplicitPositional:
			parsed.Positional = append(parsed.Positional, arg)
		case rune == PrefixLiteral:
			parsed.Positional = append(parsed.Positional, arg[1:])
		case rune == PrefixFlagPositive, rune == PrefixFlagNegative:
			flagName := arg[1:]

			if parsed.IsPositionalTime() {
				return parsed, fmt.Errorf("unexpected flag %s while processing positionals", flagName)
			}

			name := config.NormalizeName(flagName)

			if name == "" {
				return parsed, fmt.Errorf("incorrect flag %s", arg)
			}

			if rune == PrefixFlagPositive {
				parsed.Flags[name] = FlagTrue
			} else {
				parsed.Flags[name] = FlagFalse
			}
		case strings.ContainsRune(arg, SeparatorKeyword):
			if parsed.IsPositionalTime() {
				return parsed, fmt.Errorf("unexpected parameter %s while processing positionals", arg)
			}

			indexRune := strings.IndexRune(arg, SeparatorKeyword)
			name, value := arg[:indexRune], arg[indexRune+1:]
			name = config.NormalizeName(name)

			if name == "" {
				return parsed, fmt.Errorf("incorrect parameter %s", arg)
			}

			values, err := shlex.Split(value, true)
			if err != nil {
				return parsed, fmt.Errorf("cannot split parameter %s: %w", arg, err)
			}

			parsed.Parameters[name] = append(parsed.Parameters[name], values...)
		default:
			parsed.Positional = append(parsed.Positional, arg)
		}
	}

	for _, values := range parsed.Parameters {
		sort.Strings(values)
	}

	return parsed, nil
}
