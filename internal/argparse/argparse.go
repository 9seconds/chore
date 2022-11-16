package argparse

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/config"
	"github.com/alessio/shellescape"
)

const (
	SeparatorKeyword    = "="
	SeparatorPositional = "--"
)

type validatedValue struct {
	name  string
	value string
}

func Parse(ctx context.Context, parameters map[string]config.Parameter, args []string) (ParsedArgs, error) { //nolint: cyclop
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	keywords := make(map[string][]string)
	rValue := ParsedArgs{
		Keywords: make(map[string]string),
	}

	waiters := &sync.WaitGroup{}
	errChan := make(chan error)
	resChan := make(chan validatedValue)

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

		validateValue(ctx, waiters, spec, name, value, resChan, errChan)
	}

	go func() {
		waiters.Wait()
		cancel()
	}()

parametersLoop:
	for {
		select {
		case <-ctx.Done():
			break parametersLoop
		case val := <-resChan:
			keywords[val.name] = append(keywords[val.name], val.value)
		case err := <-errChan:
			return rValue, err
		}
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

func validateValue(ctx context.Context,
	waiters *sync.WaitGroup,
	spec config.Parameter,
	name, value string,
	resChan chan<- validatedValue, errChan chan<- error,
) {
	waiters.Add(1)

	go func() {
		defer waiters.Done()

		err := spec.Validate(ctx, value)

		if err != nil {
			resChan = nil
			err = fmt.Errorf("incorrect value %s for parameter %s: %w", name, value, err)
		} else {
			errChan = nil
		}

		select {
		case <-ctx.Done():
		case resChan <- validatedValue{name: name, value: value}:
		case errChan <- err:
		}
	}()
}
