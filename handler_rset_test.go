package smtpmock

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
		request, session, message, configuration := "RSET", new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgRsetReceived
		handler := newHandlerRset(session, message, configuration)
		session.On("writeResponse", receivedMessage, configuration.responseDelayRset).Once().Return(nil)
		handler.run(request)

		assert.True(t, message.rset)
	})

	t.Run("when failure RSET request", func(t *testing.T) {
		request, session, message, configuration := "RSETbroken", new(sessionMock), new(message), createConfiguration()
		handler := newHandlerRset(session, message, configuration)
		handler.run(request)

		assert.False(t, message.rset)
	})
}
