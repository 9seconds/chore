package argparse

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"sync"

	"github.com/9seconds/chore/internal/config"
)

type ParsedArgs struct {
	Parameters map[string]string
	Flags      map[string]FlagValue
	Positional []string
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

	waiters.Add(len(p.Parameters))

	go func() {
		waiters.Wait()
		close(errChan)
	}()

	for name, value := range p.Parameters {
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

	return <-errChan
}

func (p ParsedArgs) IsPositionalTime() bool {
	return len(p.Positional) > 0
}

func (p ParsedArgs) OptionsAsCli() []string {
	options := make([]string, 0, len(p.Parameters)+len(p.Flags))

	for name, value := range p.Parameters {
		options = append(options, name+string(SeparatorKeyword)+value)
	}

	for name, value := range p.Flags {
		if value == FlagTrue {
			options = append(options, string(PrefixFlagPositive)+name)
		} else {
			options = append(options, string(PrefixFlagNegative)+name)
		}
	}

	return options
}

func (p ParsedArgs) Checksum() []byte {
	parameterNames := make([]string, 0, len(p.Parameters))
	flagNames := make([]string, 0, len(p.Flags))

	for k := range p.Parameters {
		parameterNames = append(parameterNames, k)
	}

	for k := range p.Flags {
		flagNames = append(flagNames, k)
	}

	sort.Strings(parameterNames)
	sort.Strings(flagNames)

	mixer := sha256.New()

	binary.Write(mixer, binary.LittleEndian, uint64(len(p.Parameters))) //nolint: errcheck
	binary.Write(mixer, binary.LittleEndian, uint64(len(p.Flags)))      //nolint: errcheck
	binary.Write(mixer, binary.LittleEndian, uint64(len(p.Positional))) //nolint: errcheck

	for _, v := range parameterNames {
		mixer.Write([]byte(v))
		mixer.Write([]byte{0x00})
		mixer.Write([]byte(p.Parameters[v]))
		mixer.Write([]byte{0x01})
	}

	for _, v := range flagNames {
		mixer.Write([]byte(v))
		mixer.Write([]byte{0x00})
		mixer.Write([]byte{byte(p.Flags[v])})
		mixer.Write([]byte{0x02})
	}

	for _, v := range p.Positional {
		mixer.Write([]byte(v))
		mixer.Write([]byte{0x00})
	}

	return mixer.Sum(nil)
}
