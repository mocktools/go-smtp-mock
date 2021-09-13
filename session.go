package smtpmock

import (
	"bufio"
	"net"
)

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
	connection                             net.Conn
	address                                string
	bufin                                  bufin
	bufout                                 bufout
	readError, writeError, connectionError error
}

// SMTP session builder
func newSession(connection net.Conn) *session {
	return &session{
		connection: connection,
		address:    connection.RemoteAddr().String(),
		bufin:      bufio.NewReader(connection),
		bufout:     bufio.NewWriter(connection),
	}
}

// SMTP session methods

// Reades client request from the session. When error case happened writes it to readError
func (session *session) readRequest() string {
	request, err := session.bufin.ReadString('\n')
	if err != nil {
		session.readError = err
	}

	return request
}

// Writes server response to the client session. When error case happened writes it to writeError
func (session *session) writeResponse(response string) {
	bufout := session.bufout
	_, err := bufout.WriteString(response + "\r\n")
	if err != nil {
		session.writeError = err
	}
	bufout.Flush()
}

// Finishes SMTP session. When error case happened writes it to connectionError
func (session *session) finish() {
	err := session.connection.Close()
	if err != nil {
		session.connectionError = err
	}
}
