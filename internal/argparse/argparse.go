package argparse

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/config"
)

const (
	PrefixFlagPositive = "+"
	PrefixFlagNegative = "-"
	PrefixLiteral      = ":"
	SeparatorKeyword   = "="
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

func parseArgs(
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

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, PrefixLiteral):
			positionalTime = true

			parsed.Positional = append(parsed.Positional, strings.TrimPrefix(arg, PrefixLiteral))
		case strings.HasPrefix(arg, PrefixFlagPositive), strings.HasPrefix(arg, PrefixFlagNegative):
			if positionalTime {
				return parsed, fmt.Errorf("unexpected flag %s while processing positionals", arg)
			}

			name := normalizeArgName(arg[1:])

			if _, ok := flags[name]; !ok {
				return parsed, fmt.Errorf("unknown flag %s", name)
			}

			parsed.Flags[name] = strings.HasPrefix(arg, PrefixFlagPositive)
		case strings.Contains(arg, SeparatorKeyword):
			if positionalTime {
				return parsed, fmt.Errorf("unexpected parameter %s while processing positionals", arg)
			}

			name, value, _ := strings.Cut(arg, SeparatorKeyword)
			name = normalizeArgName(name)

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
					return
				default:
				}

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

func normalizeArgName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ToLower(name)

	return name
}
