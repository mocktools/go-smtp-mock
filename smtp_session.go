package smtpmock

import (
	"bufio"
	"net"
)

// smtpSession interfaces

type bufin interface {
	ReadString(byte) (string, error)
}

type bufout interface {
	WriteString(string) (int, error)
	Flush() error
}

// SMTP client-server session
type smtpSession struct {
	connection                             net.Conn
	address                                string
	bufin                                  bufin
	bufout                                 bufout
	readError, writeError, connectionError error
}

// smtpSession builder
func newSmtpSession(connection net.Conn) *smtpSession {
	return &smtpSession{
		connection: connection,
		address:    connection.RemoteAddr().String(),
		bufin:      bufio.NewReader(connection),
		bufout:     bufio.NewWriter(connection),
	}
}

// smtpSession methods

// Reades client request from the session. When error case happened writes it to readError
func (session *smtpSession) readRequest() string {
	request, err := session.bufin.ReadString('\n')
	if err != nil {
		session.readError = err
	}

	return request
}

// Writes server response to the client session. When error case happened writes it to writeError
func (session *smtpSession) writeResponse(response string) {
	bufout := session.bufout
	_, err := bufout.WriteString(response + "\r\n")
	if err != nil {
		session.writeError = err
	}
	bufout.Flush()
}

// Finishes SMTP session. When error case happened writes it to connectionError
func (session *smtpSession) finish() {
	err := session.connection.Close()
	if err != nil {
		session.connectionError = err
	}
}
