package env

import (
	"context"
	"os"
	"regexp"
)

var environCleanupRegexp = regexp.MustCompile("^" + Prefix + "(?:P|PL|F)_.*$")

func Environ() []string {
	baseEnviron := os.Environ()
	processed := make([]string, 0, len(baseEnviron))

	for _, value := range baseEnviron {
		if !environCleanupRegexp.MatchString(value) {
			processed = append(processed, value)
		}
	}

	return processed
}

func MakeValue(name, value string) string {
	return name + "=" + value
}

func sendValue(ctx context.Context, results chan<- string, name, value string) {
	if value != "" {
		select {
		case <-ctx.Done():
		case results <- MakeValue(name, value):
		}
	}
}
