package env

import "context"

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
