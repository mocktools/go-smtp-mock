package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerQuit(t *testing.T) {
	t.Run("returns new handlerQuit", func(t *testing.T) {
		session, message, configuration := new(session), new(message), new(configuration)
		handler := newHandlerQuit(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerQuitRun(t *testing.T) {
	t.Run("when successful QUIT request", func(t *testing.T) {
		request, session, message, configuration := "QUIT", new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgQuitCmd
		handler := newHandlerQuit(session, message, configuration)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run(request)

		assert.True(t, message.quitSent)
	})

	t.Run("when failure QUIT request", func(t *testing.T) {
		request, session, message, configuration := "QUIT ", new(sessionMock), new(message), createConfiguration()
		handler := newHandlerQuit(session, message, configuration)
		handler.run(request)

		assert.False(t, message.quitSent)
	})
}
