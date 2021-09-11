package smtpmock

import (
	"bufio"
	"net"
	"strings"
	"time"
)

// Returns time.Time with current time. Allows to stub time.Now()
var timeNow func() time.Time = func() time.Time { return time.Now() }

// SMTP client-server session interface
type sessionInterface interface {
	setTimeout(int)
	readRequest() (string, error)
	writeResponse(string)
	addError(error)
	clearError()
	discardBufin()
	readBytes() ([]byte, error)
	isErrorFound() bool
	finish()
}

// session interfaces

type bufin interface {
	ReadString(byte) (string, error)
	Buffered() int
	Discard(int) (int, error)
	ReadBytes(byte) ([]byte, error)
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

// Sets session timeout from now to the specified duration in seconds
func (session *session) setTimeout(timeout int) {
	err := session.connection.SetDeadline(
		timeNow().Add(time.Duration(timeout) * time.Second),
	)

	if err != nil {
		session.err = err
		session.logger.error(err.Error())
	}
}

// Discardes the bufin remnants
func (session *session) discardBufin() {
	bufin := session.bufin
	_, err := bufin.Discard(bufin.Buffered())

	if err != nil {
		session.err = err
		session.logger.error(err.Error())
	}
}

// Reades client request from the session, returns trimmed string.
// When error case happened writes it to session.err and triggers logger with error level
func (session *session) readRequest() (string, error) {
	request, err := session.bufin.ReadString('\n')
	if err == nil {
		trimmedRequest := strings.TrimSpace(request)
		session.logger.infoActivity(SessionRequestMsg + trimmedRequest)
		return trimmedRequest, err
	}

	session.err = err
	session.logger.error(err.Error())
	return EmptyString, err
}

// Reades client request from the session, returns bytes.
// When error case happened writes it to session.err and triggers logger with error level
func (session *session) readBytes() ([]byte, error) {
	var request []byte
	request, err := session.bufin.ReadBytes('\n')
	if err == nil {
		session.logger.infoActivity(SessionRequestMsg + SessionBinaryDataMsg)
		return request, err
	}

	session.err = err
	session.logger.error(err.Error())
	return request, err
}

// Writes server response to the client session. When error case happened triggers
// logger with warning level
func (session *session) writeResponse(response string) {
	bufout := session.bufout
	if _, err := bufout.WriteString(response + "\r\n"); err != nil {
		session.logger.warning(err.Error())
	}
	bufout.Flush()
	session.logger.infoActivity(SessionResponseMsg + response)
}

// Finishes SMTP session. When error case happened triggers logger with warning level
func (session *session) finish() {
	if err := session.connection.Close(); err != nil {
		session.logger.warning(err.Error())
	}

	session.logger.infoActivity(SessionEndMsg)
}
