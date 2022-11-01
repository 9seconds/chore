package argparse

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
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

func (p ParsedArgs) Checksum(seed []byte) []byte {
	argNames := make([]string, 0, len(p.Keywords))

	for k := range p.Keywords {
		argNames = append(argNames, k)
	}

	sort.Strings(argNames)

	mixer := sha256.New()
	binary.Write(mixer, binary.LittleEndian, uint64(len(seed)))
	mixer.Write(seed)

	binary.Write(mixer, binary.LittleEndian, uint64(len(argNames)))

	for _, v := range argNames {
		mixer.Write([]byte(v))
		mixer.Write([]byte{0x01})
		mixer.Write([]byte(p.Keywords[v]))
		mixer.Write([]byte{0x00})
	}

	binary.Write(mixer, binary.LittleEndian, uint64(len(p.Positional)))

	for _, v := range p.Positional {
		mixer.Write([]byte(v))
		mixer.Write([]byte{0x00})
	}

	return mixer.Sum(nil)
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
