package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerMailfrom(t *testing.T) {
	t.Run("returns new handlerMailfrom", func(t *testing.T) {
		session, message, configuration := new(session), new(message), new(configuration)
		handler := newHandlerMailfrom(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerMailfromRun(t *testing.T) {
	t.Run("when successful MAILFROM request", func(t *testing.T) {
		request := "MAIL FROM: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgMailfromReceived
		message.helo = true
		handler := newHandlerMailfrom(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run(request)

		assert.True(t, message.mailfrom)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, receivedMessage, message.mailfromResponse)
	})

	t.Run("when failure MAILFROM request, invalid command sequence", func(t *testing.T) {
		request := "MAIL FROM: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		errorMessage := configuration.msgInvalidCmdMailfromSequence
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.mailfrom)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when failure MAILFROM request, invalid command argument", func(t *testing.T) {
		request := "MAIL FROM"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		errorMessage := configuration.msgInvalidCmdMailfromArg
		message.helo = true
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.mailfrom)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when failure MAILFROM request, request includes blacklisted MAILFROM email", func(t *testing.T) {
		email := "user@example.com"
		request := "MAIL FROM: " + email
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		message.helo, configuration.blacklistedMailfromEmails = true, []string{email}
		errorMessage := configuration.msgMailfromBlacklistedEmail
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.mailfrom)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})
}

func TestHandlerMailfromClearMessage(t *testing.T) {
	t.Run("erases all handler message data from MAILFROM command, changes cleared status to true", func(t *testing.T) {
		notEmptyMessage := createNotEmptyMessage()
		handler := newHandlerMailfrom(new(session), notEmptyMessage, new(configuration))
		clearedMessage := &message{
			heloRequest:  notEmptyMessage.heloRequest,
			heloResponse: notEmptyMessage.heloResponse,
			helo:         notEmptyMessage.helo,
			cleared:      true,
		}
		handler.clearMessage()

		assert.Same(t, notEmptyMessage, handler.message)
		assert.Equal(t, clearedMessage, handler.message)

		handler.message.mailfromRequest = "42"
		handler.clearMessage()
		assert.Equal(t, clearedMessage, handler.message)
	})
}

func TestHandlerMailfromWriteResult(t *testing.T) {
	request, response := "request context", "response context"
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when successful request received", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, response, message.mailfromResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(message), errors.New(response)
		handler := newHandlerMailfrom(session, message, configuration)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, response, message.mailfromResponse)
	})
}

func TestHandlerMailfromIsInvalidCmdSequence(t *testing.T) {
	configuration, session, request := createConfiguration(), &sessionMock{}, "MAIL FROM: <user@domain.com>"

	t.Run("when request includes invalid command MAILFROM sequence, the previous command is not successful ", func(t *testing.T) {
		message, errorMessage := new(message), configuration.msgInvalidCmdMailfromSequence
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when request includes valid command MAILFROM sequence, the previous command is successful ", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)
		message.helo = true

		assert.False(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})
}

func TestHandlerMaifromIsInvalidCmdArg(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command MAILFROM argument", func(t *testing.T) {
		request, message, errorMessage := "MAIL FROM: email@invalid", new(message), configuration.msgInvalidCmdMailfromArg
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when request includes valid command MAILFROM argument without <> sign", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("MAIL FROM: user@example.com"))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})

	t.Run("when request includes valid command MAILFROM argument without <> sign without space", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("MAIL FROM:user@example.com"))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})

	t.Run("when request includes valid command MAILFROM argument with <> sign", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("MAIL FROM: <user@example.com>"))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})

	t.Run("when request includes valid command MAILFROM argument with <> sign without space", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("MAIL FROM:<user@example.com>"))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})
}

func TestHandlerMailfromMailfromEmail(t *testing.T) {
	handler := new(handlerMailfrom)

	t.Run("when request includes valid email address without <> sign", func(t *testing.T) {
		validEmail := "user@example.com"

		assert.Equal(t, validEmail, handler.mailfromEmail("MAIL FROM: "+validEmail))
	})

	t.Run("when request includes valid email address with <> sign", func(t *testing.T) {
		validEmail := "user@example.com"

		assert.Equal(t, validEmail, handler.mailfromEmail("MAIL FROM: "+"<"+validEmail+">"))
	})

	t.Run("when request includes invalid email address", func(t *testing.T) {
		invalidEmail := "user@invalid"

		assert.Equal(t, EmptyString, handler.mailfromEmail("MAIL FROM: "+invalidEmail))
	})
}

func TestHandlerHeloIsBlacklistedEmail(t *testing.T) {
	email := "user@example.com"
	request := "MAIL FROM: " + email

	t.Run("when request includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedMailfromEmails = []string{email}
		errorMessage := configuration.msgMailfromBlacklistedEmail
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isBlacklistedEmail(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when request not includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isBlacklistedEmail(request))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})
}

func TestHandlerMailfromIsInvalidRequest(t *testing.T) {
	configuration := createConfiguration()

	t.Run("when request includes invalid MAILFROM command sequence, the previous command is not successful", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdMailfromSequence
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when request includes invalid MAILFROM command argument", func(t *testing.T) {
		request := "MAIL FROM: user@example"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdMailfromArg
		message.helo = true
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when request includes blacklisted MAILFROM email", func(t *testing.T) {
		configuration, blacklistedEmail := createConfiguration(), "user@example.com"
		request := "MAIL FROM: " + blacklistedEmail
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgHeloBlacklistedDomain
		configuration.blacklistedMailfromEmails, message.helo = []string{blacklistedEmail}, true
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when valid MAILFROM request", func(t *testing.T) {
		request := "MAIL FROM: user@example.com"
		session, message := new(sessionMock), new(message)
		message.helo = true
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})
}
