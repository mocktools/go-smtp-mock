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
	connection net.Conn
	address    string
	bufin      bufin
	bufout     bufout
	err        error
}

// SMTP session builder. Creates new session
func newSession(connection net.Conn) *session {
	return &session{
		connection: connection,
		address:    connection.RemoteAddr().String(),
		bufin:      bufio.NewReader(connection),
		bufout:     bufio.NewWriter(connection),
	}
}

// SMTP session methods

// Returns true if session error exists, otherwise returns false
func (session *session) isErrorFound() bool {
	return session.err != nil
}

// Reades client request from the session. When error case happened writes it to session.err
func (session *session) readRequest() (string, error) {
	request, err := session.bufin.ReadString('\n')
	if err != nil {
		session.err = err
	}

	return request, err
}

// Writes server response to the client session. When error case happened writes it to session.err
func (session *session) writeResponse(response string) {
	bufout := session.bufout
	_, err := bufout.WriteString(response + "\r\n")
	if err != nil {
		session.err = err
	}
	bufout.Flush()
}

// Finishes SMTP session. When error case happened writes it to session.err
func (session *session) finish() {
	err := session.connection.Close()
	if err != nil {
		session.err = err
	}
}
