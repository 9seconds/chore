package testlib

import (
	"context"
	"net"

	"github.com/stretchr/testify/mock"
)

type DNSResolverMock struct {
	mock.Mock
}

func (m *DNSResolverMock) LookupAddr(ctx context.Context, addr string) ([]string, error) {
	args := m.Called(ctx, addr)
	return args.Get(0).([]string), args.Error(1)
}

func (m *DNSResolverMock) LookupMX(ctx context.Context, addr string) ([]*net.MX, error) {
	args := m.Called(ctx, addr)
	return args.Get(0).([]*net.MX), args.Error(1)
}
