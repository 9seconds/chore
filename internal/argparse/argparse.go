package argparse

import (
	"context"
	"fmt"
	"strings"

	"github.com/9seconds/chore/internal/config"
	"github.com/alessio/shellescape"
)

const (
	SeparatorKeyword    = "="
	SeparatorPositional = "--"
)

}

func Parse(ctx context.Context, parameters map[string]config.Parameter, args []string) (ParsedArgs, error) {
	keywords := make(map[string][]string)
	rValue := ParsedArgs{
		Keywords: make(map[string]string),
	}

	for idx, arg := range args {
		if arg == SeparatorPositional {
			rValue.Positional = make([]string, len(args)-idx-1)
			copy(rValue.Positional, args[idx+1:])

			break
		}

		name, value, found := strings.Cut(arg, SeparatorKeyword)
		if !found {
			return rValue, fmt.Errorf("cannot find %s separator in argument %s", SeparatorKeyword, arg)
		}

		name = strings.ToLower(name)

		spec, ok := parameters[name]
		if !ok {
			return rValue, fmt.Errorf("unknown parameter %s", name)
		}

		if err := spec.Validate(ctx, value); err != nil {
			return rValue, fmt.Errorf("incorrect value %s for parameter %s: %w", name, value, err)
		}

		keywords[name] = append(keywords[name], value)
	}

	for name, param := range parameters {
		if _, ok := keywords[name]; !ok && param.Required() {
			return rValue, fmt.Errorf("absent value for parameter %s", name)
		}
	}

	for k, v := range keywords {
		rValue.Keywords[k] = shellescape.QuoteCommand(v)
	}

	return rValue, nil
}
