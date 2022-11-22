package filelock

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrLockCannotAcquire = errors.New("cannot acquire lock")

const (
	LockMinDelayMS = 10
	LockMaxDelayMS = int(time.Second / time.Millisecond)
)

type LockType byte

const (
	LockTypeNo LockType = iota
	LockTypeExclusive
	LockTypeShared
)

func (l LockType) String() string {
	switch l {
	case LockTypeNo:
		return "no"
	case LockTypeExclusive:
		return "exclusive"
	case LockTypeShared:
		return "shared"
	}

	return fmt.Sprintf("unknown: %x", byte(l))
}

type Lock interface {
	Lock(context.Context) error
	Unlock() error
}

func New(lockType LockType, path string) (Lock, error) {
	switch lockType {
	case LockTypeNo:
		return newNoopLock(), nil
	case LockTypeExclusive:
		return newLockExclusive(path)
	case LockTypeShared:
		return newLockShared(path)
	}

	return nil, fmt.Errorf("unknown lock type %v", lockType)
}

func getDelay() time.Duration {
	base := rand.Intn(LockMaxDelayMS-LockMinDelayMS) + LockMinDelayMS

	return time.Millisecond * time.Duration(base)
}
