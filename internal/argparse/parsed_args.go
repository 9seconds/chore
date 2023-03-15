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
	"github.com/gosimple/slug"
)

const (
	SerializePrefixPositional   = '0'
	SerializePrefixFlagPositive = '1'
	SerializePrefixFlagNegative = '2'
	SerializePrefixParameter    = '3'
	SerializeParameterSeparator = '_'
)

type ParsedArgs struct {
	Parameters         map[string][]string
	Flags              map[string]string
	Positional         []string
	ExplicitPositional bool
}

func (p ParsedArgs) GetParameterList(key string) string {
	builder := strings.Builder{}

	for _, v := range p.Parameters[key] {
		builder.WriteString(v)
		builder.WriteRune('\n')
	}

	return builder.String()
}

func (p ParsedArgs) GetParameter(key string) string {
	if values := p.Parameters[key]; len(values) > 0 {
		return values[len(values)-1]
	}

	return ""
}

func (p ParsedArgs) ToSelfStringChunks() []string {
	chunks := make([]string, 0, len(p.Flags)+len(p.Parameters))

	for _, key := range binutils.SortedMapKeys(p.Flags) {
		prefix := PrefixFlagNegative
		if p.Flags[key] == FlagTrue {
			prefix = PrefixFlagPositive
		}

		chunks = append(chunks, fmt.Sprintf("%c%s", prefix, key))
	}

	for _, key := range binutils.SortedMapKeys(p.Parameters) {
		for _, value := range p.Parameters[key] {
			chunks = append(chunks, fmt.Sprintf("%s%c%s", key, SeparatorKeyword, value))
		}
	}

	return chunks
}

func (p ParsedArgs) ToSlugString() string {
	chunks := make([]string, 0, len(p.Positional)+len(p.Flags)+len(p.Parameters))

	for _, v := range p.Positional {
		chunks = append(chunks, fmt.Sprintf("%c%s", SerializePrefixPositional, v))
	}

	params := make([]string, 0, len(p.Parameters))

	for key := range p.Parameters {
		params = append(params, key)
	}

	// sort param keys by an overall length of whole k=v sum, groupped by key
	sort.Slice(params, func(one, another int) bool {
		getLength := func(key string) int {
			sum := len(p.Parameters) * (1 + len(key))

			for _, v := range p.Parameters[key] {
				sum += len(v)
			}

			return sum
		}

		return getLength(params[one]) < getLength(params[another])
	})

	for _, key := range params {
		var values []string

		values = append(values, p.Parameters[key]...)

		sort.Sort(sort.Reverse(sort.StringSlice(values)))

		for _, val := range values {
			chunks = append(
				chunks,
				fmt.Sprintf(
					"%c%s%c%s",
					SerializePrefixParameter,
					key,
					SerializeParameterSeparator,
					val))
		}
	}

	for _, key := range binutils.SortedMapKeys(p.Flags) {
		prefix := SerializePrefixFlagPositive

		if p.Flags[key] != FlagTrue {
			prefix = SerializePrefixFlagNegative
		}

		chunks = append(chunks, fmt.Sprintf("%c%s", prefix, key))
	}

	return slug.Make(strings.Join(chunks, " "))
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

	for name, values := range p.Parameters {
		for _, value := range values {
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

	binutils.MixLength(mixer, len(p.Parameters)) //nolint: errcheck

	for _, key := range binutils.SortedMapKeys(p.Parameters) {
		binutils.MixString(mixer, key)                    //nolint: errcheck
		binutils.MixStringSlice(mixer, p.Parameters[key]) //nolint: errcheck
	}

	for _, key := range binutils.SortedMapKeys(p.Flags) {
		binutils.MixString(mixer, key)          //nolint: errcheck
		binutils.MixString(mixer, p.Flags[key]) //nolint: errcheck
	}

	binutils.MixStringSlice(mixer, p.Positional) //nolint: errcheck

	return binutils.ToString(mixer.Sum(nil))
}
