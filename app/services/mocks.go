package services

import (
	"net"

	"github.com/stretchr/testify/mock"
)

type WSMock struct {
	mock.Mock
}

func (m *WSMock) ReadJSON(v interface{}) error {
	return m.Called(v).Error(0)
}

func (m *WSMock) WriteJSON(v interface{}) error {
	return m.Called(v).Error(0)
}

func (m *WSMock) RemoteAddr() net.Addr {
	args := m.Called()

	var addr net.Addr
	if item := args.Get(0); item != nil {
		addr = item.(net.Addr)
	}

	return addr
}

func (m *WSMock) Close() error {
	return m.Called().Error(0)
}
