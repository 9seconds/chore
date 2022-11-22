package filelock

import "context"

type lockNoop struct{}

func (l lockNoop) Lock(_ context.Context) error { return nil }
func (l lockNoop) Unlock() error                { return nil }

func newNoopLock() lockNoop {
	return lockNoop{}
}
