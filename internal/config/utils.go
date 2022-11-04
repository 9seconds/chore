package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parseCSV(value string) []string {
	splitted := strings.Split(value, ",")
	results := make([]string, 0, len(splitted))

	for _, v := range splitted {
		v = strings.TrimSpace(v)

		if v != "" {
			results = append(results, v)
		}
	}

	return results
}

func parseBool(spec map[string]string, name string) (bool, error) {
	if value, ok := spec[name]; ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return false, fmt.Errorf("cannot parse %s: %w", name, err)
		}

		return parsed, nil
	}

	return false, nil
}

func parseRegexp(spec map[string]string, name string) (*regexp.Regexp, error) {
	if value, ok := spec[name]; ok {
		parsed, err := regexp.Compile(value)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %s: %w", name, err)
		}

		return parsed, nil
	}

	return nil, nil
}
