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

// bufio.Reader mock
type bufioReaderMock struct {
	mock.Mock
}

func (buf bufioReaderMock) ReadString(delim byte) (string, error) {
	args := buf.Called(delim)
	return args.String(0), args.Error(1)
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

// logger mock
type loggerMock struct {
	mock.Mock
}

func (logger *loggerMock) infoActivity(message string) {
	logger.Called(message)
}

func (logger *loggerMock) info(message string) {
	logger.Called(message)
}

func (logger *loggerMock) warning(message string) {
	logger.Called(message)
}

func (logger *loggerMock) error(message string) {
	logger.Called(message)
}

// session mock
type sessionMock struct {
	mock.Mock
}

func (session *sessionMock) readRequest() (string, error) {
	args := session.Called()
	return args.String(0), args.Error(1)
}

func (session *sessionMock) writeResponse(response string) {
	session.Called(response)
}

func (session *sessionMock) addError(err error) {
	session.Called(err)
}

func (session *sessionMock) clearError() {
	session.Called()
}
