package argparse

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/9seconds/chore/internal/config"
)

const (
	PrefixFlagPositive = '+'
	PrefixFlagNegative = '-'
	PrefixLiteral      = ':'
	SeparatorKeyword   = '='
)

func Parse(
	ctx context.Context,
	args []string,
	flags map[string]bool,
	parameters map[string]config.Parameter,
) (ParsedArgs, error) {
	parsed, err := parseArgs(args, flags, parameters)
	if err != nil {
		return parsed, err
	}

	return parsed, validateArgs(ctx, parsed, flags, parameters)
}

func parseArgs( //nolint: cyclop
	args []string,
	flags map[string]bool,
	parameters map[string]config.Parameter,
) (ParsedArgs, error) {
	parsed := ParsedArgs{
		Parameters: make(map[string]string),
		Flags:      make(map[string]bool),
		Positional: []string{},
	}

	positionalTime := false

	for idx, arg := range args {
		if !utf8.ValidString(arg) {
			return parsed, fmt.Errorf("argument %d is not valid UTF-8 string", idx+1)
		}

		rune, _ := utf8.DecodeRuneInString(arg)

		switch {
		case rune == PrefixLiteral:
			positionalTime = true

			parsed.Positional = append(parsed.Positional, arg[1:])
		case rune == PrefixFlagPositive, rune == PrefixFlagNegative:
			flagName := arg[1:]

			if positionalTime {
				return parsed, fmt.Errorf("unexpected flag %s while processing positionals", flagName)
			}

			name := config.NormalizeName(flagName)

			if _, ok := flags[name]; !ok {
				return parsed, fmt.Errorf("unknown flag %s", flagName)
			}

			parsed.Flags[name] = rune == PrefixFlagPositive
		case strings.ContainsRune(arg, SeparatorKeyword):
			if positionalTime {
				return parsed, fmt.Errorf("unexpected parameter %s while processing positionals", arg)
			}

			indexRune := strings.IndexRune(arg, SeparatorKeyword)
			name, value := arg[:indexRune], arg[indexRune+1:]
			name = config.NormalizeName(name)

			if _, ok := parameters[name]; !ok {
				return parsed, fmt.Errorf("unknown parameter %s", name)
			}

			parsed.Parameters[name] = value
		default:
			positionalTime = true

			parsed.Positional = append(parsed.Positional, arg)
		}
	}

	return parsed, nil
}

func validateArgs( //nolint: cyclop
	ctx context.Context,
	args ParsedArgs,
	flags map[string]bool,
	parameters map[string]config.Parameter,
) error {
	for name, required := range flags {
		if _, ok := args.Flags[name]; !ok && required {
			return fmt.Errorf("flag '%s' is required but value is not provided", name)
		}
	}

	for name, param := range parameters {
		if _, ok := args.Parameters[name]; !ok && param.Required() {
			return fmt.Errorf("parameter '%s' is required but value is not provided", name)
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	waiters := &sync.WaitGroup{}
	errChan := make(chan error)

	waiters.Add(len(args.Parameters))

	for name, value := range args.Parameters {
		go func(name, value string) {
			defer waiters.Done()

			if err := parameters[name].Validate(ctx, value); err != nil {
				select {
				case <-ctx.Done():
				case errChan <- fmt.Errorf("invalid value for parameter %s: %w", name, err):
				}
			}
		}(name, value)
	}

	go func() {
		waiters.Wait()
		close(errChan)
	}()

	return <-errChan
}
