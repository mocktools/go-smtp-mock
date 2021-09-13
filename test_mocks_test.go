package smtpmock

import (
	"net"
	"time"

	"github.com/stretchr/testify/mock"
)

// Testing mocks

// net.Addr mock
type netAddressMock struct {
	mock.Mock
}

func (addr netAddressMock) Network() string {
	args := addr.Called()
	return args.String(0)
}

func (addr netAddressMock) String() string {
	args := addr.Called()
	return args.String(0)
}

// net.Conn mock
type netConnectionMock struct {
	mock.Mock
}

func (connection netConnectionMock) LocalAddr() net.Addr {
	args := connection.Called()
	return args.Get(0).(net.Addr)
}

func (connection netConnectionMock) RemoteAddr() net.Addr {
	args := connection.Called()
	return args.Get(0).(net.Addr)
}

func (connection netConnectionMock) Read(b []byte) (n int, err error) {
	args := connection.Called(b)
	return args.Get(0).(int), args.Error(1)
}

func (connection netConnectionMock) Write(b []byte) (n int, err error) {
	args := connection.Called(b)
	return args.Get(0).(int), args.Error(1)
}

func (connection netConnectionMock) SetDeadline(t time.Time) error {
	args := connection.Called(t)
	return args.Error(0)
}

func (connection netConnectionMock) SetReadDeadline(t time.Time) error {
	args := connection.Called(t)
	return args.Error(0)
}

func (connection netConnectionMock) SetWriteDeadline(t time.Time) error {
	args := connection.Called(t)
	return args.Error(0)
}

func (connection netConnectionMock) Close() error {
	args := connection.Called()
	return args.Error(0)
}

// bufio.Writer mock
type bufioWriterMock struct {
	mock.Mock
}

func (buf bufioWriterMock) WriteString(s string) (int, error) {
	args := buf.Called(s)
	return args.Int(0), args.Error(1)
}

func (buf bufioWriterMock) Flush() error {
	args := buf.Called()
	return args.Error(0)
}
