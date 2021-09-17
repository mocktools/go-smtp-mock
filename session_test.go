package smtpmock

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionIsErrorFound(t *testing.T) {
	t.Run("when error exists", func(t *testing.T) {
		session := &session{err: errors.New("some error messsage")}

		assert.True(t, session.isErrorFound())
	})

	t.Run("when error not exists", func(t *testing.T) {
		assert.False(t, new(session).isErrorFound())
	})
}

func TestNewSession(t *testing.T) {
	t.Run("creates new SMTP session", func(t *testing.T) {
		connectionAddress := "127.0.0.1:25"
		connection, address := netConnectionMock{}, netAddressMock{}
		address.On("String").Once().Return(connectionAddress)
		connection.On("RemoteAddr").Once().Return(address)
		session := newSession(connection)

		assert.Equal(t, connection, session.connection)
		assert.Equal(t, connectionAddress, session.address)
		assert.Equal(t, bufio.NewReader(connection), session.bufin)
		assert.Equal(t, bufio.NewWriter(connection), session.bufout)
	})
}

func TestSessionReadRequest(t *testing.T) {
	t.Run("extracts string from bufin without error", func(t *testing.T) {
		capturedStringContext := "Some string context\n"
		stringContext := capturedStringContext + "other string"
		binaryData := strings.NewReader(stringContext)
		bufin := bufio.NewReader(binaryData)
		session := &session{bufin: bufin}
		request, err := session.readRequest()

		assert.Equal(t, capturedStringContext, request)
		assert.NoError(t, err)
		assert.NoError(t, session.err)
	})

	t.Run("extracts string from bufin with error", func(t *testing.T) {
		capturedStringContext := "Some string context"
		binaryData := bytes.NewBufferString(capturedStringContext)
		bufin := bufio.NewReader(binaryData)
		session := &session{bufin: bufin}
		request, err := session.readRequest()

		assert.Equal(t, capturedStringContext, request)
		assert.Error(t, err)
		assert.Same(t, session.err, err)
	})
}

func TestSessionWriteResponse(t *testing.T) {
	t.Run("writes server response to bufout without error", func(t *testing.T) {
		response := "some response"
		binaryData := bytes.NewBufferString("")
		bufout := bufio.NewWriter(binaryData)
		session := &session{bufout: bufout}
		session.writeResponse(response)

		assert.Equal(t, response+"\r\n", binaryData.String())
		assert.NoError(t, session.err)
	})

	t.Run("writes server response to bufout with error", func(t *testing.T) {
		response, errorMessage := "some response", "write error"
		err := errors.New(errorMessage)
		bufout := new(bufioWriterMock)
		session := &session{bufout: bufout}
		bufout.On("WriteString", response+"\r\n").Once().Return(0, err)
		bufout.On("Flush").Once().Return(err)
		session.writeResponse(response)

		assert.EqualError(t, session.err, errorMessage)
	})
}

func TestSessionFinish(t *testing.T) {
	t.Run("closes session connection without error", func(t *testing.T) {
		connection := netConnectionMock{}
		connection.On("Close").Once().Return(nil)
		session := &session{connection: connection}
		session.finish()

		assert.NoError(t, session.err)
	})

	t.Run("closes session connection with error", func(t *testing.T) {
		errorMessage := "connection error"
		connection, err := netConnectionMock{}, errors.New(errorMessage)
		connection.On("Close").Once().Return(err)
		session := &session{connection: connection}
		session.finish()

		assert.EqualError(t, session.err, errorMessage)
	})
}
