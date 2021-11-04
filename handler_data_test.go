package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerData(t *testing.T) {
	t.Run("returns new handlerData", func(t *testing.T) {
		session, message, configuration := new(session), new(message), new(configuration)
		handler := newHandlerData(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerDataRun(t *testing.T) {
	t.Run("when read request error", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler, err := newHandlerData(session, message, configuration), errors.New("some read error")
		session.On("readRequest").Once().Return(EmptyString, err)
		handler.run()

		assert.False(t, message.data)
		assert.Empty(t, message.dataRequest)
		assert.Empty(t, message.dataResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid SMTP command", func(t *testing.T) {
		request := "DATE"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmd
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid DATA command sequence", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmdDataSequence
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when fail fast scenario enabled, successful DATA request", func(t *testing.T) {
		request := "DATA"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler := newHandlerData(session, message, configuration)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", DefaultReadyForReceiveMsg).Once().Return(nil)
		handler.run()

		assert.True(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, DefaultReadyForReceiveMsg, message.dataResponse)
	})

	t.Run("when fail fast scenario disabled, read request error during loop session", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, err := newHandlerData(session, message, configuration), errors.New("some read error")
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(EmptyString, err)
		handler.run()

		assert.False(t, message.data)
		assert.Empty(t, message.dataRequest)
		assert.Empty(t, message.dataResponse)
	})

	t.Run("when fail fast scenario disabled, successful DATA request", func(t *testing.T) {
		request := "DATA"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgDataReceived
		handler := newHandlerData(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, receivedMessage, message.dataResponse)
	})
}

func TestHandlerDataWriteResult(t *testing.T) {
	request, response := "request context", "response context"
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when successful request received", func(t *testing.T) {
		message := new(message)
		handler := newHandlerData(session, message, configuration)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, response, message.dataResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(message), errors.New(response)
		handler := newHandlerData(session, message, configuration)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, response, message.dataResponse)
	})
}

func TestHandlerDataIsInvalidCmd(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request, message, errorMessage := "DATE", new(message), configuration.msgInvalidCmd
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmd(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when request includes valid SMTP command", func(t *testing.T) {
		message := new(message)
		handler := newHandlerData(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("DATA"))
		assert.False(t, message.data)
		assert.Empty(t, message.dataRequest)
		assert.Empty(t, message.dataResponse)
	})
}

func TestHandlerDataIsInvalidCmdSequence(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command DATA sequence", func(t *testing.T) {
		request, message, errorMessage := "RCPT TO:", new(message), configuration.msgInvalidCmdDataSequence
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when request includes valid command DATA sequence", func(t *testing.T) {
		message := new(message)
		handler := newHandlerData(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("DATA"))
		assert.False(t, message.data)
		assert.Empty(t, message.dataRequest)
		assert.Empty(t, message.dataResponse)
	})
}

func TestHandlerDataIsInvalidRequest(t *testing.T) {
	configuration := createConfiguration()

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request := "DATE"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmd
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when request includes invalid DATA command sequence", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdDataSequence
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when valid DATA request", func(t *testing.T) {
		request := "DATA"
		session, message := new(sessionMock), new(message)
		handler := newHandlerData(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
		assert.False(t, message.data)
		assert.Empty(t, message.dataRequest)
		assert.Empty(t, message.dataResponse)
	})
}
