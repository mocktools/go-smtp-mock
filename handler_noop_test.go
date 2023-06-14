package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerNoop(t *testing.T) {
	t.Run("returns new handleNoop", func(t *testing.T) {
		session, message, configuration := new(session), new(Message), new(configuration)
		handler := newHandlerNoop(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerNoopRun(t *testing.T) {
	t.Run("when successful NOOP request", func(t *testing.T) {
		request, session, message, configuration := "NOOP", new(sessionMock), new(Message), createConfiguration()
		receivedMessage := configuration.msgNoopReceived
		handler := newHandlerNoop(session, message, configuration)
		session.On("writeResponse", receivedMessage, configuration.responseDelayQuit).Once().Return(nil)
		handler.run(request)

		assert.True(t, message.noop)
	})

	t.Run("when failure NOOP request", func(t *testing.T) {
		request, session, message, configuration := "NOOP ", new(sessionMock), new(Message), createConfiguration()
		handler := newHandlerNoop(session, message, configuration)
		handler.run(request)

		assert.False(t, message.noop)
	})
}

func TestHandlerNoopIsInvalidRequest(t *testing.T) {
	handler := newHandlerNoop(new(session), new(Message), new(configuration))

	t.Run("when request includes invalid NOOP command", func(t *testing.T) {
		request := "NOOP "

		assert.True(t, handler.isInvalidRequest(request))
	})

	t.Run("when request includes valid NOOP command", func(t *testing.T) {
		request := "NOOP"

		assert.False(t, handler.isInvalidRequest(request))
	})
}
