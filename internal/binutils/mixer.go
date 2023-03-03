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
	MixTypeStringMap   byte = 0x03
)

func MixString(writer io.Writer, str string) error {
	if _, err := writer.Write([]byte{MixTypeString}); err != nil {
		return fmt.Errorf("cannot mix string type: %w", err)
	}

	if err := binary.Write(writer, binary.LittleEndian, uint64(len(str))); err != nil {
		return fmt.Errorf("cannot mix string length: %w", err)
	}

	if _, err := io.WriteString(writer, str); err != nil {
		return fmt.Errorf("cannot mix string: %w", err)
	}

	return nil
}

func MixStringSlice(writer io.Writer, strings []string) error {
	if _, err := writer.Write([]byte{MixTypeStringSlice}); err != nil {
		return fmt.Errorf("cannot mix string slice type: %w", err)
	}

	if err := binary.Write(writer, binary.LittleEndian, uint64(len(strings))); err != nil {
		return fmt.Errorf("cannot mix string length: %w", err)
	}

	for idx := range strings {
		if err := MixString(writer, strings[idx]); err != nil {
			return fmt.Errorf("cannot mix %d string: %w", idx, err)
		}
	}

	return nil
}

func MixStringsMap(writer io.Writer, data map[string]string) error {
	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	if _, err := writer.Write([]byte{MixTypeStringMap}); err != nil {
		return fmt.Errorf("cannot mix string map type: %w", err)
	}

	if err := binary.Write(writer, binary.LittleEndian, uint64(len(data))); err != nil {
		return fmt.Errorf("cannot mix string length: %w", err)
	}

	for _, key := range keys {
		if err := MixString(writer, key); err != nil {
			return fmt.Errorf("cannot mix key %s: %w", key, err)
		}

		if err := MixString(writer, data[key]); err != nil {
			return fmt.Errorf("cannot mix value of %s: %w", key, err)
		}
	}

	return nil
}
