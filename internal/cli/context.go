package cli

import (
	"context"
	"time"

	"github.com/9seconds/chore/internal/filelock"
)

type Context struct {
	ctx    context.Context
	cancel context.CancelFunc
	lock   filelock.Lock
}

func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c Context) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c Context) Err() error {
	return c.ctx.Err()
}

func (c Context) Value(key any) any {
	return c.ctx.Value(key)
}

func (c Context) Close() error {
	c.cancel()

	if c.lock != nil {
		return c.lock.Unlock()
	}

	return nil
}

func (c Context) WithTimeout(timeout time.Duration) Context {
	c.ctx, c.cancel = context.WithTimeout(c.ctx, timeout)

	return c
}

func (c Context) WithLock(lock filelock.Lock) Context {
	if c.lock != nil {
		panic("lock is already set")
	}

	c.lock = lock

	return c
}

func (c Context) Start() error {
	return c.lock.Lock(c)
}

func NewContext(ctx context.Context) Context {
	ctx, cancel := context.WithCancel(ctx)

	return Context{
		ctx:    ctx,
		cancel: cancel,
	}
}
