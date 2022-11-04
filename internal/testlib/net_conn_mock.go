package testlib

import (
	"net"
	"time"

	"github.com/stretchr/testify/mock"
)

type NetConnMock struct {
	mock.Mock
}

func (m *NetConnMock) Read(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *NetConnMock) Write(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *NetConnMock) Close() error {
	return m.Called().Error(0)
}

func (m *NetConnMock) LocalAddr() net.Addr {
	return m.Called().Get(0).(net.Addr)
}

func (m *NetConnMock) RemoteAddr() net.Addr {
	return m.Called().Get(0).(net.Addr)
}

func (m *NetConnMock) SetDeadline(t time.Time) error {
	return m.Called(t).Error(0)
}

func (m *NetConnMock) SetReadDeadline(t time.Time) error {
	return m.Called(t).Error(0)
}

func (m *NetConnMock) SetWriteDeadline(t time.Time) error {
	return m.Called(t).Error(0)
}
