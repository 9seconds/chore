package env

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"

	"github.com/9seconds/chore/chorelib/config"
	"github.com/9seconds/chore/chorelib/script"
)

func generateArgs(ctx context.Context, result chan<- string, script script.Script, args map[string][]string) {
	argNames := make([]string, 0, len(args))

	for name := range args {
		argNames = append(argNames, name)
	}

	sort.Strings(argNames)
	hasher := sha256.New()

	hasher.Write([]byte(script.Path()))
	hasher.Write([]byte{0x00})
	binary.Write(hasher, binary.LittleEndian, uint32(len(argNames)))

	for _, name := range argNames {
		value := strings.Join(args[name], ScriptArgListSeparator)

		sendEnvValue(ctx, result, EnvArgPrefix+strings.ToUpper(name), value)
		hasher.Write([]byte(value))
		hasher.Write([]byte{0x00})
	}

	sendEnvValue(ctx, result, EnvCacheId, encodeBytes(hasher.Sum(nil)))
}

func parseArgs(confParameters map[string]config.Parameter, args []string) (map[string][]string, error) {
	values := make(map[string][]string)

	for _, arg := range args {
		name, value, found := strings.Cut(arg, ArgKeywordSeparator)
		if !found {
			return nil, fmt.Errorf("cannot find %s separator in argument %s", ArgKeywordSeparator, arg)
		}

		name = strings.ToLower(name)

		spec, ok := confParameters[name]
		if !ok {
			return nil, fmt.Errorf("unknown parameter %s", name)
		}

		if err := spec.Validate(value); err != nil {
			return nil, fmt.Errorf("incorrect value %s for parameter %s: %w", name, value, err)
		}

		values[name] = append(values[name], value)
	}

	for name, param := range confParameters {
		if _, ok := values[name]; !ok && param.Required() {
			return nil, fmt.Errorf("value for %s was not specified", name)
		}
	}

	return values, nil
}
