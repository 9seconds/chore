package argparse

import (
	"fmt"
	"strings"

	"github.com/9seconds/chore/internal/config"
	"github.com/alessio/shellescape"
)

const (
	SeparatorKeyword    = "="
	SeparatorPositional = "--"
)

type ParsedArgs struct {
	Keywords   map[string]string
	Positional []string
}

func Parse(parameters map[string]config.Parameter, args []string) (ParsedArgs, error) {
	keywords := make(map[string][]string)
	rv := ParsedArgs{
		Keywords: make(map[string]string),
	}

	for idx, arg := range args {
		if arg == SeparatorPositional {
			rv.Positional = make([]string, len(args)-idx-1)
			copy(rv.Positional, args[idx+1:])

			break
		}

		name, value, found := strings.Cut(arg, SeparatorKeyword)
		if !found {
			return rv, fmt.Errorf("cannot find %s separator in argument %s", SeparatorKeyword, arg)
		}

		name = strings.ToLower(name)

		spec, ok := parameters[name]
		if !ok {
			return rv, fmt.Errorf("unknown parameter %s", name)
		}

		if err := spec.Validate(value); err != nil {
			return rv, fmt.Errorf("incorrect value %s for parameter %s: %w", name, value, err)
		}

		keywords[name] = append(keywords[name], value)
	}

	for name, param := range parameters {
		if _, ok := keywords[name]; !ok && param.Required() {
			return rv, fmt.Errorf("absent value for parameter %s", name)
		}
	}

	for k, v := range keywords {
		rv.Keywords[k] = shellescape.QuoteCommand(v)
	}

	return rv, nil
}
