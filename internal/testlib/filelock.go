package testlib

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type FileLockMock struct {
	mock.Mock
}

func (m *FileLockMock) Lock(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *FileLockMock) Unlock() error {
	return m.Called().Error(0)
}
