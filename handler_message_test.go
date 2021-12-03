package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerMessage(t *testing.T) {
	t.Run("returns new handlerMessage", func(t *testing.T) {
		session, message, configuration := new(session), new(message), new(configuration)
		handler := newHandlerMessage(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerMessageRun(t *testing.T) {
	t.Run("when read request error", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler, err := newHandlerMessage(session, message, configuration), errors.New("some read error")
		session.On("readBytes").Once().Return([]byte{}, err)
		handler.run()

		assert.False(t, message.msg)
		assert.Empty(t, message.msgRequest)
		assert.Empty(t, message.msgResponse)
	})

	t.Run("when message size limit reached", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), newConfiguration(ConfigurationAttr{msqSizeLimit: 1})
		errorMessage := configuration.msgMsgSizeIsTooBig
		handler, err := newHandlerMessage(session, message, configuration), errors.New(errorMessage)
		session.On("readBytes").Once().Return([]uint8("some message"), nil)
		session.On("readBytes").Once().Return([]uint8(".\r\n"), nil)
		session.On("discardBufin").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.msg)
		assert.Equal(t, emptyString, message.msgRequest)
		assert.Equal(t, errorMessage, message.msgResponse)
	})

	t.Run("when message received", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, msgContext := newHandlerMessage(session, message, configuration), "some message"
		session.On("readBytes").Once().Return([]uint8("."+msgContext), nil)
		session.On("readBytes").Once().Return([]uint8(".\r\n"), nil)
		session.On("writeResponse", defaultReceivedMsg).Once().Return(nil)
		handler.run()

		assert.True(t, message.msg)
		assert.Equal(t, msgContext, message.msgRequest)
		assert.Equal(t, defaultReceivedMsg, message.msgResponse)
	})
}

func TestHandlerMessageWriteResult(t *testing.T) {
	request, response := "request context", "response context"
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when successful request received", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMessage(session, message, configuration)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.msg)
		assert.Equal(t, request, message.msgRequest)
		assert.Equal(t, response, message.msgResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(message), errors.New(response)
		handler := newHandlerMessage(session, message, configuration)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.msg)
		assert.Equal(t, request, message.msgRequest)
		assert.Equal(t, response, message.msgResponse)
	})
}
