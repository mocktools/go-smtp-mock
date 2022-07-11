package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerRset(t *testing.T) {
	t.Run("returns new handlerRset", func(t *testing.T) {
		session, message, configuration := new(session), new(message), new(configuration)
		handler := newHandlerRset(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerRsetRun(t *testing.T) {
	t.Run("when successful RSET request", func(t *testing.T) {
		request := "RSET"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgRsetReceived
		message.helo = true
		handler := newHandlerRset(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", receivedMessage, configuration.responseDelayRset).Once().Return(nil)
		handler.run(request)

		assert.True(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, receivedMessage, message.rsetResponse)
	})

	t.Run("when failure RSET request, invalid command sequence", func(t *testing.T) {
		request := "RSET"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		errorMessage := configuration.msgInvalidCmdRsetSequence
		handler, err := newHandlerRset(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayRset).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, errorMessage, message.rsetResponse)
	})

	t.Run("when failure RSET request, invalid command", func(t *testing.T) {
		request := "RSET "
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		errorMessage := configuration.msgInvalidCmdRsetArg
		message.helo = true
		handler, err := newHandlerRset(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayRset).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, errorMessage, message.rsetResponse)
	})
}

func TestHandlerRsetClearMessage(t *testing.T) {
	t.Run("when not multiple message receiving condition erases all message data except HELO/EHLO command context", func(t *testing.T) {
		notEmptyMessage := createNotEmptyMessage()
		handler := newHandlerRset(new(session), notEmptyMessage, new(configuration))
		clearedMessage := &message{
			heloRequest:  notEmptyMessage.heloRequest,
			heloResponse: notEmptyMessage.heloResponse,
			helo:         notEmptyMessage.helo,
		}
		handler.clearMessage()

		assert.Same(t, notEmptyMessage, handler.message)
		assert.Equal(t, clearedMessage, handler.message)

		handler.message.rsetRequest = "42"
		handler.clearMessage()
		assert.Equal(t, clearedMessage, handler.message)
	})

	t.Run("when multiple message receiving condition does not updated message", func(t *testing.T) {
		configuration, message := &configuration{multipleMessageReceiving: true}, createNotEmptyMessage()
		handler := newHandlerRset(new(session), message, configuration)
		handler.clearMessage()

		assert.Equal(t, message, handler.message)
	})
}

func TestHandlerRsetWriteResult(t *testing.T) {
	request, response := "request context", "response context"
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when successful request received", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRset(session, message, configuration)
		session.On("writeResponse", response, configuration.responseDelayRset).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, response, message.rsetResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(message), errors.New(response)
		handler := newHandlerRset(session, message, configuration)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response, configuration.responseDelayRset).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, response, message.rsetResponse)
	})
}

func TestHandlerRsetIsInvalidCmdSequence(t *testing.T) {
	request, configuration, session := "some request", createConfiguration(), &sessionMock{}

	t.Run("when helo previous command was failure ", func(t *testing.T) {
		message, errorMessage := new(message), configuration.msgInvalidCmdRsetSequence
		handler, err := newHandlerRset(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayRset).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, errorMessage, message.rsetResponse)
	})

	t.Run("when helo previous command was successful ", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRset(session, message, configuration)
		message.helo = true

		assert.False(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.rset)
		assert.Empty(t, message.rsetRequest)
		assert.Empty(t, message.rsetResponse)
	})
}

func TestHandlerRsetIsInvalidCmdArg(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid RSET command", func(t *testing.T) {
		request, message, errorMessage := "RSET ", new(message), configuration.msgInvalidCmdRsetArg
		handler, err := newHandlerRset(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayRset).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, errorMessage, message.rsetResponse)
	})

	t.Run("when request includes valid RSET command", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRset(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("RSET"))
		assert.False(t, message.rset)
		assert.Empty(t, message.rsetRequest)
		assert.Empty(t, message.rsetResponse)
	})
}

func TestHandlerRsetIsInvalidRequest(t *testing.T) {
	configuration := createConfiguration()

	t.Run("when request includes invalid RSET command sequence, the previous command is not successful", func(t *testing.T) {
		request := "RSET"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdRsetSequence
		handler, err := newHandlerRset(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayRset).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, errorMessage, message.rsetResponse)
	})

	t.Run("when request includes invalid RSET command", func(t *testing.T) {
		request := "RSET "
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdRsetArg
		message.helo = true
		handler, err := newHandlerRset(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayRset).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.rset)
		assert.Equal(t, request, message.rsetRequest)
		assert.Equal(t, errorMessage, message.rsetResponse)
	})

	t.Run("when valid RSET request", func(t *testing.T) {
		request := "RSET"
		session, message := new(sessionMock), new(message)
		message.helo = true
		handler := newHandlerRset(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
		assert.False(t, message.rset)
		assert.Empty(t, message.rsetRequest)
		assert.Empty(t, message.rsetResponse)
	})
}
