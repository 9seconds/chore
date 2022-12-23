package argparse

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

const (
	FlagUnknown = ""
	FlagTrue    = "t"
	FlagFalse   = "f"
)

type ParsedArgs struct {
	Parameters map[string]string
	Flags      map[string]bool
	Positional []string
}

func (p ParsedArgs) Options() []string {
	options := make([]string, 0, len(p.Parameters)+len(p.Flags))

	for name, value := range p.Parameters {
		options = append(options, name+string(SeparatorKeyword)+value)
	}

	for name, required := range p.Flags {
		if required {
			options = append(options, string(PrefixFlagPositive)+name)
		} else {
			options = append(options, string(PrefixFlagNegative)+name)
		}
	}

	return options
}

// todo CHORE_RECURSION

func (p ParsedArgs) GetFlagValue(name string) string {
	value, ok := p.Flags[name]

	switch {
	case !ok:
		return FlagUnknown
	case value:
		return FlagTrue
	}

	return FlagFalse
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
		mixer.Write([]byte(p.GetFlagValue(v)))
		mixer.Write([]byte{0x02})
	}

	for _, v := range p.Positional {
		mixer.Write([]byte(v))
		mixer.Write([]byte{0x00})
	}

	return mixer.Sum(nil)
}
