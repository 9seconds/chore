package cli

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Timeout time.Duration

func (t Timeout) Value() time.Duration {
	return time.Duration(t)
}

func (t *Timeout) UnmarshalText(b []byte) error {
	text := string(b)

	if value, err := strconv.ParseUint(text, 10, 64); err == nil {
		*t = Timeout(time.Duration(value) * time.Second)

		return nil
	}

	value, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf("cannot parse duration: %w", err)
	}

	if value < 0 {
		return errors.New("duration should be positive")
	}

	*t = Timeout(value)

	return nil
}
