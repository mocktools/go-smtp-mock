package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerData(t *testing.T) {
	t.Run("returns new handlerData", func(t *testing.T) {
		session, message, configuration := new(session), new(Message), new(configuration)
		handler := newHandlerData(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerDataRun(t *testing.T) {
	t.Run("when successful DATA request", func(t *testing.T) {
		request, session, message, configuration := "DATA", new(sessionMock), new(Message), createConfiguration()
		handlerMessage, receivedMessage := &handlerMessageMock{}, configuration.msgDataReceived
		message.helo, message.mailfrom, message.rcptto = true, true, true
		handler, responseDelay := newHandlerData(session, message, configuration), configuration.responseDelayData
		handler.handlerMessage = handlerMessage
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", defaultReadyForReceiveMsg, responseDelay).Once().Return(nil)
		session.On("writeResponse", receivedMessage, responseDelay).Once().Return(nil)
		handlerMessage.On("run").Once().Return(nil)
		handler.run(request)

		assert.True(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, receivedMessage, message.dataResponse)
	})

	t.Run("when failure DATA request, invalid command sequence", func(t *testing.T) {
		request := "DATA"
		session, message, configuration := new(sessionMock), new(Message), createConfiguration()
		errorMessage := configuration.msgInvalidCmdDataSequence
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayData).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when failure DATA request, invalid command", func(t *testing.T) {
		request := "DATA:"
		session, message, configuration := new(sessionMock), new(Message), createConfiguration()
		message.helo, message.mailfrom, message.rcptto = true, true, true
		errorMessage := configuration.msgInvalidCmd
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayData).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})
}

func TestHandlerDataClearMessage(t *testing.T) {
	t.Run("erases all handler message data from DATA command", func(t *testing.T) {
		notEmptyMessage := createNotEmptyMessage()
		handler := newHandlerData(new(session), notEmptyMessage, new(configuration))
		clearedMessage := &Message{
			heloRequest:           notEmptyMessage.heloRequest,
			heloResponse:          notEmptyMessage.heloResponse,
			helo:                  notEmptyMessage.helo,
			mailfromRequest:       notEmptyMessage.mailfromRequest,
			mailfromResponse:      notEmptyMessage.mailfromResponse,
			mailfrom:              notEmptyMessage.mailfrom,
			rcpttoRequestResponse: notEmptyMessage.rcpttoRequestResponse,
			rcptto:                notEmptyMessage.rcptto,
		}
		handler.clearMessage()

		assert.Same(t, notEmptyMessage, handler.message)
		assert.Equal(t, clearedMessage, handler.message)

		handler.message.dataRequest = "42"
		handler.clearMessage()
		assert.Equal(t, clearedMessage, handler.message)
	})
}

func TestHandlerDataProcessIncomingMessage(t *testing.T) {
	t.Run("when successful request received", func(*testing.T) {
		handlerMessage := &handlerMessageMock{}
		handler := &handlerData{handlerMessage: handlerMessage}
		handler.handlerMessage = handlerMessage
		handlerMessage.On("run").Once().Return(nil)
		handler.processIncomingMessage()
	})
}

func TestHandlerDataWriteResult(t *testing.T) {
	request, response := "request context", "response context"
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when successful request received", func(t *testing.T) {
		message := new(Message)
		handler := newHandlerData(session, message, configuration)
		session.On("writeResponse", response, configuration.responseDelayData).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, response, message.dataResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(Message), errors.New(response)
		handler := newHandlerData(session, message, configuration)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response, configuration.responseDelayData).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, response, message.dataResponse)
	})
}

func TestHandlerDataIsInvalidCmdSequence(t *testing.T) {
	request, configuration, session := "some request", createConfiguration(), &sessionMock{}

	t.Run("when none of the previous command was successful", func(t *testing.T) {
		message, errorMessage := new(Message), configuration.msgInvalidCmdDataSequence
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayData).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when rcptto previous command was failure", func(t *testing.T) {
		message, errorMessage := new(Message), configuration.msgInvalidCmdDataSequence
		message.helo, message.mailfrom = true, true
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayData).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when mailfrom, rcptto previous commands were failure", func(t *testing.T) {
		message, errorMessage := new(Message), configuration.msgInvalidCmdDataSequence
		message.helo = true
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayData).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when all of the previous commands was successful", func(t *testing.T) {
		message := new(Message)
		message.helo, message.mailfrom, message.rcptto = true, true, true
		handler := newHandlerData(session, message, configuration)

		assert.False(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.data)
		assert.Empty(t, message.dataRequest)
		assert.Empty(t, message.dataResponse)
	})
}

func TestHandlerDataIsInvalidCmd(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command DATA", func(t *testing.T) {
		request, message, errorMessage := "DATA ", new(Message), configuration.msgInvalidCmdDataSequence
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayData).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when request includes valid command DATA", func(t *testing.T) {
		message := new(Message)
		handler := newHandlerData(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("DATA"))
		assert.False(t, message.data)
		assert.Empty(t, message.dataRequest)
		assert.Empty(t, message.dataResponse)
	})
}

func TestHandlerDataIsInvalidRequest(t *testing.T) {
	request, configuration, session := "DATA", createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid DATA command sequence", func(t *testing.T) {
		message, errorMessage := new(Message), configuration.msgInvalidCmdDataSequence
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayData).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when request includes invalid command DATA", func(t *testing.T) {
		request, message, errorMessage := "DATA:", new(Message), configuration.msgInvalidCmd
		message.helo, message.mailfrom, message.rcptto = true, true, true
		handler, err := newHandlerData(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayData).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.data)
		assert.Equal(t, request, message.dataRequest)
		assert.Equal(t, errorMessage, message.dataResponse)
	})

	t.Run("when valid DATA request", func(t *testing.T) {
		message := new(Message)
		message.helo, message.mailfrom, message.rcptto = true, true, true
		handler := newHandlerData(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
		assert.False(t, message.data)
		assert.Empty(t, message.dataRequest)
		assert.Empty(t, message.dataResponse)
	})
}
