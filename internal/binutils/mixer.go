package binutils

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

const (
	MixTypeString      byte = 0x01
	MixTypeStringSlice byte = 0x02
	MixTypeMap         byte = 0x03
)

func MixString(writer io.Writer, str string) error {
	if _, err := writer.Write([]byte{MixTypeString}); err != nil {
		return fmt.Errorf("cannot write header: %w", err)
	}

	if err := MixLength(writer, len(str)); err != nil {
		return fmt.Errorf("cannot mix length: %w", err)
	}

	if _, err := io.WriteString(writer, str); err != nil {
		return fmt.Errorf("cannot mix string: %w", err)
	}

	return nil
}

func MixStringSlice(writer io.Writer, data []string) error {
	if _, err := writer.Write([]byte{MixTypeStringSlice}); err != nil {
		return fmt.Errorf("cannot write header: %w", err)
	}

	if err := MixLength(writer, len(data)); err != nil {
		return fmt.Errorf("cannot mix length: %w", err)
	}

	for _, v := range data {
		if err := MixString(writer, v); err != nil {
			return fmt.Errorf("cannot mix %s: %w", v, err)
		}
	}

	return nil
}

func MixLength(writer io.Writer, length int) error {
	return binary.Write(writer, binary.LittleEndian, uint64(length))
}

func SortedMapKeys[T ~string, V any](data map[T]V) []string {
	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, string(k))
	}

	sort.Strings(keys)

	return keys
}
