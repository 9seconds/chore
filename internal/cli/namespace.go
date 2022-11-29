package cli

import (
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/env"
)

type Namespace string

func (n Namespace) Value() string {
	return string(n)
}

func (n *Namespace) UnmarshalText(b []byte) error {
	text := string(b)

	if text != MagicValue {
		*n = Namespace(text)

		return nil
	}

	text, ok := os.LookupEnv(env.EnvNamespace)
	if !ok {
		return fmt.Errorf("namespace is magic but no value for %s is provided", env.EnvNamespace)
	}

	*n = Namespace(text)

	return nil
}
