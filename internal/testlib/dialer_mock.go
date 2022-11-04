package testlib

import (
	"context"
	"net"

	"github.com/stretchr/testify/mock"
)

type DialerMock struct {
	mock.Mock
}

func (m *DialerMock) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	args := m.Called(ctx, network, address)
	return args.Get(0).(net.Conn), args.Error(1)
}
