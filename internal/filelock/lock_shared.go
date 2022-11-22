package filelock

import (
	"context"
	"fmt"
	"os"

	"github.com/gofrs/flock"
)

type lockShared struct {
	lock   *flock.Flock
	locked bool
}

func (l *lockShared) Lock(ctx context.Context) error {
	if l.locked {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ok, err := l.lock.TryRLockContext(ctx, getDelay())

	switch {
	case err != nil:
		return fmt.Errorf("cannot acquire lock: %w", err)
	case !ok:
		return ErrLockCannotAcquire
	}

	l.locked = true

	return nil
}

func (l *lockShared) Unlock() error {
	if !l.locked {
		return nil
	}

	l.locked = false

	return l.lock.Unlock()
}

func newLockShared(path string) (*lockShared, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot stat path: %w", err)
	}

	return &lockShared{
		lock: flock.New(path),
	}, nil
}
