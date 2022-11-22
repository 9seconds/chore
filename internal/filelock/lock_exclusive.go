package filelock

import (
	"context"
	"fmt"
	"os"

	"github.com/gofrs/flock"
)

type lockExclusive struct {
	lock   *flock.Flock
	locked bool
}

func (l *lockExclusive) Lock(ctx context.Context) error {
	if l.locked {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ok, err := l.lock.TryLockContext(ctx, getDelay())

	switch {
	case err != nil:
		return fmt.Errorf("cannot acquire lock: %w", err)
	case !ok:
		return ErrLockCannotAcquire
	}

	l.locked = true

	return nil
}

func (l *lockExclusive) Unlock() error {
	if !l.locked {
		return nil
	}

	l.locked = false

	return l.lock.Unlock()
}

func newLockExclusive(path string) (*lockExclusive, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot stat path: %w", err)
	}

	return &lockExclusive{
		lock: flock.New(path),
	}, nil
}
