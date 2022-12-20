package argparse

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/config"
	"github.com/alessio/shellescape"
)

const (
	SeparatorKeyword = ":"
)

type validatedValue struct {
	index int
	name  string
	value string
}

func Parse(ctx context.Context, parameters map[string]config.Parameter, args []string) (ParsedArgs, error) { //nolint: cyclop
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	keywords := make(map[string]map[int]string)
	rValue := ParsedArgs{
		Keywords: make(map[string]string),
	}

	waiters := &sync.WaitGroup{}
	errChan := make(chan error)
	resChan := make(chan validatedValue)
	keywordStage := true

	for idx, arg := range args {
		name, value, found := strings.Cut(arg, SeparatorKeyword)

		if !found {
			keywordStage = false

			rValue.Positional = append(rValue.Positional, arg)

			continue
		}

		if !keywordStage {
			return rValue, fmt.Errorf("unexpected keyword parameter %s", arg)
		}

		name = strings.ToLower(name)
		name = strings.ReplaceAll(name, "-", "_")

		spec, ok := parameters[name]
		if !ok {
			return rValue, fmt.Errorf("unknown parameter %s", name)
		}

		validateValue(ctx, waiters, spec, idx, name, value, resChan, errChan)
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
			if keywords[val.name] == nil {
				keywords[val.name] = make(map[int]string)
			}

			keywords[val.name][val.index] = val.value
		case err := <-errChan:
			return rValue, err
		}
	}

	for name, param := range parameters {
		if _, ok := keywords[name]; !ok && param.Required() {
			return rValue, fmt.Errorf("absent value for parameter %s", name)
		}
	}

	for name, values := range keywords {
		orders := make([]int, 0, len(values))
		kwValues := make([]string, 0, len(values))

		for idx := range values {
			orders = append(orders, idx)
		}

		sort.Ints(orders)

		for _, idx := range orders {
			kwValues = append(kwValues, values[idx])
		}

		rValue.Keywords[name] = shellescape.QuoteCommand(kwValues)
	}

	return rValue, nil
}

func validateValue(ctx context.Context,
	waiters *sync.WaitGroup,
	spec config.Parameter,
	index int,
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
		case resChan <- validatedValue{index: index, name: name, value: value}:
		case errChan <- err:
		}
	}()
}
