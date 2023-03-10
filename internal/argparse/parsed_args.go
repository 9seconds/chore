package argparse

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/binutils"
	"github.com/9seconds/chore/internal/script/config"
	"github.com/alessio/shellescape"
)

const (
	SerializePrefixPositional   = "0"
	SerializePrefixFlagPositive = "1"
	SerializePrefixFlagNegative = "2"
	SerializePrefixParameter    = "3"
	SerializeParameterSeparator = "_"
)

type ParsedArgs struct {
	Parameters         map[string]string
	Flags              map[string]string
	Positional         []string
	ExplicitPositional bool
	ListDelimiter      string
}

func (p ParsedArgs) SerializedString() string {
	result := make([]string, 0, len(p.Positional)+len(p.Flags)+len(p.Parameters))
	params := make([]string, 0, len(p.Parameters))
	flags := make([]string, 0, len(p.Flags))

	for _, v := range p.Positional {
		result = append(result, SerializePrefixPositional+v)
	}

	for k, v := range p.Parameters {
		params = append(params, SerializePrefixParameter+k+SerializeParameterSeparator+v)
	}

	for key, value := range p.Flags {
		if value == FlagTrue {
			flags = append(flags, SerializePrefixFlagPositive+key)
		} else {
			flags = append(flags, SerializePrefixFlagNegative+key)
		}
	}

	sort.Strings(params)
	sort.Strings(flags)

	result = append(result, params...)
	result = append(result, flags...)

	return shellescape.QuoteCommand(result)
}

func (p ParsedArgs) Validate( //nolint: cyclop
	ctx context.Context,
	flags map[string]config.Flag,
	parameters map[string]config.Parameter,
) error {
	for name := range p.Flags {
		if _, ok := flags[name]; !ok {
			return fmt.Errorf("unknown flag %s", name)
		}
	}

	for name, flag := range flags {
		if _, ok := p.Flags[name]; !ok && flag.Required() {
			return fmt.Errorf("mandatory flag %s was not provided", name)
		}
	}

	for name := range p.Parameters {
		if _, ok := parameters[name]; !ok {
			return fmt.Errorf("unknown parameter %s", name)
		}
	}

	for name, parameter := range parameters {
		if _, ok := p.Parameters[name]; !ok && parameter.Required() {
			return fmt.Errorf("mandatory parameter %s was not provided", name)
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	waiters := &sync.WaitGroup{}
	errChan := make(chan error)

	for name, listValues := range p.Parameters {
		for _, value := range strings.Split(listValues, p.ListDelimiter) {
			waiters.Add(1)

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
	}

	go func() {
		waiters.Wait()
		close(errChan)
	}()

	return <-errChan
}

func (p ParsedArgs) IsPositionalTime() bool {
	return p.ExplicitPositional || len(p.Positional) > 0
}

func (p ParsedArgs) Checksum() string {
	mixer := sha256.New()

	binutils.MixStringsMap(mixer, p.Parameters)  //nolint: errcheck
	binutils.MixStringsMap(mixer, p.Flags)       //nolint: errcheck
	binutils.MixStringSlice(mixer, p.Positional) //nolint: errcheck

	return binutils.ToString(mixer.Sum(nil))
}
