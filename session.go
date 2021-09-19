package smtpmock

import (
	"bufio"
	"net"
	"strings"
)

// SMTP client-server session interface
type sessionInterface interface {
	readRequest() (string, error)
	writeResponse(string)
	addError(error)
	clearError()
}

// session interfaces

type bufin interface {
	ReadString(byte) (string, error)
}

type bufout interface {
	WriteString(string) (int, error)
	Flush() error
}

// SMTP client-server session
type session struct {
	connection net.Conn
	address    string
	bufin      bufin
	bufout     bufout
	err        error
	logger     logger
}

// SMTP session builder. Creates new session
func newSession(connection net.Conn, logger logger) *session {
	return &session{
		connection: connection,
		address:    connection.RemoteAddr().String(),
		bufin:      bufio.NewReader(connection),
		bufout:     bufio.NewWriter(connection),
		logger:     logger,
	}
}

// SMTP session methods

// Returns true if session error exists, otherwise returns false
func (session *session) isErrorFound() bool {
	return session.err != nil
}

// session.err setter
func (session *session) addError(err error) {
	session.err = err
}

// Sets session.err = nil
func (session *session) clearError() {
	session.err = nil
}

// Reades client request from the session, returns trimmed string.
// When error case happened writes it to session.err and triggers logger with error level
func (session *session) readRequest() (string, error) {
	request, err := session.bufin.ReadString('\n')
	if err == nil {
		return strings.TrimSpace(request), err
	}

	session.err = err
	session.logger.error(err.Error())
	return EmptyString, err
}

// Writes server response to the client session. When error case happened triggers
// logger with warning level
func (session *session) writeResponse(response string) {
	bufout := session.bufout
	if _, err := bufout.WriteString(response + "\r\n"); err != nil {
		session.logger.warning(err.Error())
	}
	bufout.Flush()
}

// Finishes SMTP session. When error case happened triggers logger with warning level
func (session *session) finish() {
	if err := session.connection.Close(); err != nil {
		session.logger.warning(err.Error())
	}
}
