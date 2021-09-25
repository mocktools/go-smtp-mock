package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerWriteResult(t *testing.T) {
	request, response := "request context", "response context"
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when successful request received", func(t *testing.T) {
		message := new(message)
		handler := &handler{session, message, configuration}
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, response, message.heloResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(message), errors.New(response)
		handler := &handler{session, message, configuration}
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, response, message.heloResponse)
	})
}

func TestHandlerIsInvalidCmd(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request, message := "HI", new(message)
		handler, err := &handler{session, message, configuration}, errors.New(DefaultInvalidCmdMsg)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdMsg).Once().Return(nil)

		assert.True(t, handler.isInvalidCmd(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdMsg, message.heloResponse)
	})

	t.Run("when request includes valid SMTP command", func(t *testing.T) {
		message := new(message)
		handler := &handler{session, message, configuration}

		assert.False(t, handler.isInvalidCmd("RCPT TO:"))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})
}
