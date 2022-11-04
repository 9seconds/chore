package argparse

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

type ParsedArgs struct {
	Keywords   map[string]string
	Positional []string
}

func (p ParsedArgs) Checksum() []byte {
	argNames := make([]string, 0, len(p.Keywords))

	for k := range p.Keywords {
		argNames = append(argNames, k)
	}

	sort.Strings(argNames)

	mixer := sha256.New()
	binary.Write(mixer, binary.LittleEndian, uint64(len(argNames))) //nolint: errcheck

	for _, v := range argNames {
		mixer.Write([]byte(v))
		mixer.Write([]byte{0x01})
		mixer.Write([]byte(p.Keywords[v]))
		mixer.Write([]byte{0x00})
	}

	binary.Write(mixer, binary.LittleEndian, uint64(len(p.Positional))) //nolint: errcheck

	for _, v := range p.Positional {
		mixer.Write([]byte(v))
		mixer.Write([]byte{0x00})
	}

	return mixer.Sum(nil)
}
