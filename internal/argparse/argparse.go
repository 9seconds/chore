package argparse

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/9seconds/chore/internal/script/config"
)

const (
	PrefixFlagPositive = '+'
	PrefixFlagNegative = '-'
	PrefixLiteral      = ':'
	SeparatorKeyword   = '='

	FlagTrue  = "t"
	FlagFalse = "f"

	DefaultListDelimiter = ":"
	PositionalDelimiter  = "--"
)

var ListDelimiter = string(os.PathListSeparator)

func Parse(args []string, listDelimiter string) (ParsedArgs, error) { //nolint: cyclop
	parsed := ParsedArgs{
		Parameters:    make(map[string]string),
		Flags:         make(map[string]string),
		Positional:    []string{},
		ListDelimiter: listDelimiter,
	}

	for idx, arg := range args {
		if !utf8.ValidString(arg) {
			return parsed, fmt.Errorf("argument %d is not valid UTF-8 string", idx+1)
		}

		rune, _ := utf8.DecodeRuneInString(arg)

		switch {
		case arg == PositionalDelimiter && !parsed.ExplicitPositional:
			parsed.ExplicitPositional = true
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

			parsed.Parameters[name] = strings.TrimPrefix(value, listDelimiter)
		default:
			parsed.Positional = append(parsed.Positional, arg)
		}
	}

	return parsed, nil
}
