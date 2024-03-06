package smtpmock

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeNow(t *testing.T) {
	t.Run("wrappes time.Now() in function", func(t *testing.T) {
		assert.Equal(t, time.Now().Day(), timeNow().Day())
	})
}

func TestTimeSleep(t *testing.T) {
	t.Run("wrappes time.Sleep() in function, returns delay", func(t *testing.T) {
		delay := 0

		assert.Equal(t, delay, timeSleep(delay))
	})
}

func TestSessionIsErrorFound(t *testing.T) {
	t.Run("when error exists", func(t *testing.T) {
		session := &session{err: errors.New("some error messsage")}

		assert.True(t, session.isErrorFound())
	})

	t.Run("when error not exists", func(t *testing.T) {
		assert.False(t, new(session).isErrorFound())
	})
}

func TestSessionAddError(t *testing.T) {
	t.Run("assigns error to session.err", func(t *testing.T) {
		err, session := errors.New("some error messsage"), new(session)
		session.addError(err)

		assert.Error(t, err, session.err)
	})
}

func TestSessionClearError(t *testing.T) {
	t.Run("clears session.err", func(t *testing.T) {
		session := &session{err: errors.New("some error messsage")}
		session.clearError()

		assert.NoError(t, session.err)
	})
}

func TestSessionSetTimeout(t *testing.T) {
	timeStub, timeout := time.Now(), 42
	timeNow = func() time.Time { return timeStub }

	t.Run("sets connection deadline for session", func(*testing.T) {
		connection := netConnectionMock{}
		connection.On("SetDeadline", timeNow().Add(time.Duration(timeout)*time.Second)).Once().Return(nil)
		session := &session{connection: connection}
		session.setTimeout(timeout)
	})

	t.Run("when connection error", func(t *testing.T) {
		errorMessage, connection, logger := "some connection error", netConnectionMock{}, new(loggerMock)
		err := errors.New(errorMessage)
		connection.On("SetDeadline", timeNow().Add(time.Duration(timeout)*time.Second)).Once().Return(err)
		logger.On("error", errorMessage).Once().Return(nil)
		session := &session{connection: connection, logger: logger}
		session.setTimeout(timeout)

		assert.Error(t, session.err)
		assert.Same(t, session.err, err)
	})
}

func TestSessionDiscardBufin(t *testing.T) {
	t.Run("discardes the bufin remnants", func(*testing.T) {
		bufin := new(bufioReaderMock)
		session := &session{bufin: bufin}
		bufin.On("Buffered").Once().Return(42)
		bufin.On("Discard", 42).Once().Return(42, nil)
		session.discardBufin()
	})

	t.Run("discardes the bufin remnants with error", func(t *testing.T) {
		errorMessage, bufin, logger := "bufin discard error", new(bufioReaderMock), new(loggerMock)
		session, err := &session{bufin: bufin, logger: logger}, errors.New(errorMessage)
		bufin.On("Buffered").Once().Return(42)
		bufin.On("Discard", 42).Once().Return(42, err)
		logger.On("error", errorMessage).Once().Return(nil)
		session.discardBufin()

		assert.Error(t, session.err)
		assert.Same(t, session.err, err)
	})
}

func TestNewSession(t *testing.T) {
	t.Run("creates new SMTP session", func(t *testing.T) {
		connectionAddress := "127.0.0.1:25"
		connection, address, logger := netConnectionMock{}, netAddressMock{}, new(loggerMock)
		address.On("String").Once().Return(connectionAddress)
		connection.On("RemoteAddr").Once().Return(address)
		session := newSession(connection, logger)

		assert.Equal(t, connection, session.connection)
		assert.Equal(t, connectionAddress, session.address)
		assert.Equal(t, bufio.NewReader(connection), session.bufin)
		assert.Equal(t, bufio.NewWriter(connection), session.bufout)
		assert.Equal(t, logger, session.logger)
	})
}

func TestSessionReadRequest(t *testing.T) {
	t.Run("extracts trimmed string from bufin without error", func(t *testing.T) {
		capturedStringContext := "Some string context"
		stringContext := capturedStringContext + "\r\n other string"
		binaryData := strings.NewReader(stringContext)
		bufin, logger := bufio.NewReader(binaryData), new(loggerMock)
		session := &session{bufin: bufin, logger: logger}
		logger.On("infoActivity", sessionRequestMsg+capturedStringContext).Once().Return(nil)
		request, err := session.readRequest()

		assert.Equal(t, capturedStringContext, request)
		assert.NoError(t, err)
		assert.NoError(t, session.err)
	})

	t.Run("extracts string from bufin with error", func(t *testing.T) {
		var delim uint8 = '\n'
		errorMessage, bufin, logger := "read error", new(bufioReaderMock), new(loggerMock)
		err := errors.New(errorMessage)
		bufin.On("ReadString", delim).Once().Return(emptyString, err)
		logger.On("error", errorMessage).Once().Return(nil)
		session := &session{bufin: bufin, logger: logger}
		request, err := session.readRequest()

		assert.Equal(t, emptyString, request)
		assert.Error(t, err)
		assert.Same(t, session.err, err)
	})
}

func TestSessionReadBytes(t *testing.T) {
	t.Run("extracts line in bytes from bufin without error", func(t *testing.T) {
		str := "stringContext\n"
		bufin, logger := bufio.NewReader(strings.NewReader(str)), new(loggerMock)
		session := &session{bufin: bufin, logger: logger}
		logger.On("infoActivity", sessionRequestMsg+sessionBinaryDataMsg).Once().Return(nil)
		request, err := session.readBytes()

		assert.Equal(t, []uint8(str), request)
		assert.NoError(t, err)
		assert.NoError(t, session.err)
	})

	t.Run("extracts line in bytes from bufin with error", func(t *testing.T) {
		var delim uint8 = '\n'
		errorMessage, bufin, logger := "read error", new(bufioReaderMock), new(loggerMock)
		err := errors.New(errorMessage)
		bufin.On("ReadBytes", delim).Once().Return([]byte{}, err)
		logger.On("error", errorMessage).Once().Return(nil)
		session := &session{bufin: bufin, logger: logger}
		request, err := session.readBytes()

		assert.Equal(t, []byte{}, request)
		assert.Error(t, err)
		assert.Same(t, session.err, err)
	})
}

func TestSessionResponseDelay(t *testing.T) {
	t.Run("when default session response delay", func(t *testing.T) {
		assert.Equal(t, defaultSessionResponseDelay, new(session).responseDelay(0))
	})

	t.Run("when custom session response delay", func(t *testing.T) {
		timeSleep = func(delay int) int { return delay }
		delay, logger := 42, new(loggerMock)
		logger.On("infoActivity", fmt.Sprintf("%s: %d sec", sessionResponseDelayMsg, delay)).Once().Return(nil)
		session := &session{logger: logger}

		assert.Equal(t, delay, session.responseDelay(delay))
	})
}

func TestSessionWriteResponse(t *testing.T) {
	t.Run("writes server response to bufout without response delay and error", func(t *testing.T) {
		response := "some response"
		binaryData := bytes.NewBufferString("")
		bufout, logger := bufio.NewWriter(binaryData), new(loggerMock)
		logger.On("infoActivity", sessionResponseMsg+response).Once().Return(nil)
		session := &session{bufout: bufout, logger: logger}
		session.writeResponse(response, defaultSessionResponseDelay)

		assert.Equal(t, response+"\r\n", binaryData.String())
		assert.NoError(t, session.err)
	})

	t.Run("writes server response to bufout with response delay and without error", func(t *testing.T) {
		timeSleep = func(delay int) int { return delay }
		response, delay := "some response", 42
		binaryData := bytes.NewBufferString("")
		bufout, logger := bufio.NewWriter(binaryData), new(loggerMock)
		logger.On("infoActivity", sessionResponseMsg+response).Once().Return(nil)
		logger.On("infoActivity", fmt.Sprintf("%s: %d sec", sessionResponseDelayMsg, delay)).Once().Return(nil)
		session := &session{bufout: bufout, logger: logger}
		session.writeResponse(response, delay)

		assert.Equal(t, response+"\r\n", binaryData.String())
		assert.NoError(t, session.err)
	})

	t.Run("writes server response to bufout with error", func(t *testing.T) {
		response, errorMessage, bufout, logger := "some response", "write error", new(bufioWriterMock), new(loggerMock)
		err := errors.New(errorMessage)
		bufout.On("WriteString", response+"\r\n").Once().Return(0, err)
		bufout.On("Flush").Once().Return(err)
		logger.On("warning", errorMessage).Once().Return(nil)
		logger.On("infoActivity", sessionResponseMsg+response).Once().Return(nil)
		session := &session{bufout: bufout, logger: logger}
		session.writeResponse(response, defaultSessionResponseDelay)

		assert.NoError(t, session.err)
	})
}

func TestSessionFinish(t *testing.T) {
	t.Run("closes session connection without error", func(t *testing.T) {
		connection, logger := netConnectionMock{}, new(loggerMock)
		connection.On("Close").Once().Return(nil)
		logger.On("infoActivity", sessionEndMsg).Once().Return(nil)
		session := &session{connection: connection, logger: logger}
		session.finish()

		assert.NoError(t, session.err)
	})

	t.Run("closes session connection with error", func(t *testing.T) {
		errorMessage := "connection error"
		connection, logger, err := netConnectionMock{}, new(loggerMock), errors.New(errorMessage)
		connection.On("Close").Once().Return(err)
		logger.On("warning", errorMessage).Once().Return(nil)
		logger.On("infoActivity", sessionEndMsg).Once().Return(nil)
		session := &session{connection: connection, logger: logger}
		session.finish()

		assert.NoError(t, session.err)
	})
}
