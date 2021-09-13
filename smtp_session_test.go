package smtpmock

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSmtpSession(t *testing.T) {
	t.Run("creates new SMTP session", func(t *testing.T) {
		connectionAddress := "127.0.0.1:25"
		connection, address := netConnectionMock{}, netAddressMock{}
		address.On("String").Once().Return(connectionAddress)
		connection.On("RemoteAddr").Once().Return(address)
		smtpSession := newSmtpSession(connection)

		assert.Equal(t, connection, smtpSession.connection)
		assert.Equal(t, connectionAddress, smtpSession.address)
		assert.Equal(t, bufio.NewReader(connection), smtpSession.bufin)
		assert.Equal(t, bufio.NewWriter(connection), smtpSession.bufout)
	})
}

func TestSmtpSessionReadRequest(t *testing.T) {
	t.Run("extracts string from bufin without error", func(t *testing.T) {
		capturedStringContext := "Some string context\n"
		stringContext := capturedStringContext + "other string"
		binaryData := strings.NewReader(stringContext)
		bufin := bufio.NewReader(binaryData)
		smtpSession := &smtpSession{bufin: bufin}

		assert.Equal(t, capturedStringContext, smtpSession.readRequest())
		assert.NoError(t, smtpSession.readError)
	})

	t.Run("extracts string from bufin with error", func(t *testing.T) {
		capturedStringContext := "Some string context"
		binaryData := bytes.NewBufferString(capturedStringContext)
		bufin := bufio.NewReader(binaryData)
		smtpSession := &smtpSession{bufin: bufin}

		assert.Equal(t, capturedStringContext, smtpSession.readRequest())
		assert.Error(t, smtpSession.readError)
	})
}

func TestSmtpSessionWriteResponse(t *testing.T) {
	t.Run("writes server response to bufout without error", func(t *testing.T) {
		response := "some response"
		binaryData := bytes.NewBufferString("")
		bufout := bufio.NewWriter(binaryData)
		smtpSession := &smtpSession{bufout: bufout}
		smtpSession.writeResponse(response)

		assert.Equal(t, response+"\r\n", binaryData.String())
		assert.NoError(t, smtpSession.writeError)
	})

	t.Run("writes server response to bufout with error", func(t *testing.T) {
		response, errorMessage := "some response", "write error"
		err := errors.New(errorMessage)
		bufout := new(bufioWriterMock)
		smtpSession := &smtpSession{bufout: bufout}
		bufout.On("WriteString", response+"\r\n").Once().Return(0, err)
		bufout.On("Flush").Once().Return(err)
		smtpSession.writeResponse(response)

		assert.EqualError(t, smtpSession.writeError, errorMessage)
	})
}

func TestSmtpSessionFinish(t *testing.T) {
	t.Run("closes session connection without error", func(t *testing.T) {
		connection := netConnectionMock{}
		connection.On("Close").Once().Return(nil)
		smtpSession := &smtpSession{connection: connection}
		smtpSession.finish()

		assert.NoError(t, smtpSession.connectionError)
	})

	t.Run("closes session connection with error", func(t *testing.T) {
		errorMessage := "connection error"
		connection, err := netConnectionMock{}, errors.New(errorMessage)
		connection.On("Close").Once().Return(err)
		smtpSession := &smtpSession{connection: connection}
		smtpSession.finish()

		assert.EqualError(t, smtpSession.connectionError, errorMessage)
	})
}
