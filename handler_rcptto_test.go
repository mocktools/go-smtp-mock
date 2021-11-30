package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerRcptto(t *testing.T) {
	t.Run("returns new handlerRcptto", func(t *testing.T) {
		session, message, configuration := new(session), new(message), new(configuration)
		handler := newHandlerRcptto(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerRcpttoRun(t *testing.T) {
	t.Run("when successful RCPTTO request", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgRcpttoReceived
		message.helo, message.mailfrom = true, true
		handler := newHandlerRcptto(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run(request)

		assert.True(t, message.rcptto)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, receivedMessage, message.rcpttoResponse)
	})

	t.Run("when failure RCPTTO request, invalid command sequence", func(t *testing.T) {
		request := "MAIL FROM: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		errorMessage := configuration.msgInvalidCmdRcpttoSequence
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.rcptto)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when failure RCPTTO request, invalid command argument", func(t *testing.T) {
		request := "RCPT TO: user@example"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		message.helo, message.mailfrom = true, true
		errorMessage := configuration.msgInvalidCmdRcpttoArg
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.rcptto)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when failure RCPTTO request, request includes blacklisted RCPTTO email", func(t *testing.T) {
		email := "user@example.com"
		request := "RCPT TO: " + email
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast, configuration.blacklistedRcpttoEmails = true, []string{email}
		message.helo, message.mailfrom = true, true
		errorMessage := configuration.msgRcpttoBlacklistedEmail
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.rcptto)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when failure RCPTTO request, request includes not registered RCPTTO email", func(t *testing.T) {
		email := "user@example.com"
		request := "RCPT TO: " + email
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast, configuration.notRegisteredEmails = true, []string{email}
		message.helo, message.mailfrom = true, true
		errorMessage := configuration.msgRcpttoNotRegisteredEmail
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.rcptto)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})
}

func TestHandlerRcpttoClearMessage(t *testing.T) {
	t.Run("erases all handler message data from RCPTTO command, changes cleared status to true", func(t *testing.T) {
		notEmptyMessage := createNotEmptyMessage()
		handler := newHandlerRcptto(new(session), notEmptyMessage, new(configuration))
		clearedMessage := &message{
			heloRequest:      notEmptyMessage.heloRequest,
			heloResponse:     notEmptyMessage.heloResponse,
			helo:             notEmptyMessage.helo,
			mailfromRequest:  notEmptyMessage.mailfromRequest,
			mailfromResponse: notEmptyMessage.mailfromResponse,
			mailfrom:         notEmptyMessage.mailfrom,
			cleared:          true,
		}
		handler.clearMessage()

		assert.Same(t, notEmptyMessage, handler.message)
		assert.Equal(t, clearedMessage, handler.message)

		handler.message.rcpttoRequest = "42"
		handler.clearMessage()
		assert.Equal(t, clearedMessage, handler.message)
	})
}

func TestHandlerRcpttoWriteResult(t *testing.T) {
	request, response := "request context", "response context"
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when successful request received", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRcptto(session, message, configuration)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, response, message.rcpttoResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(message), errors.New(response)
		handler := newHandlerRcptto(session, message, configuration)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, response, message.rcpttoResponse)
	})
}

func TestHandlerRcpttoIsInvalidCmdSequence(t *testing.T) {
	request, configuration, session := "some request", createConfiguration(), &sessionMock{}

	t.Run("when none of the previous command was successful", func(t *testing.T) {
		message, errorMessage := new(message), configuration.msgInvalidCmdRcpttoSequence
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when mailfrom previous command was failure", func(t *testing.T) {
		message, errorMessage := new(message), configuration.msgInvalidCmdRcpttoSequence
		message.helo = true
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when all of the previous commands was successful", func(t *testing.T) {
		message := new(message)
		message.helo, message.mailfrom = true, true
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})
}

func TestHandlerRcpttoIsInvalidCmdArg(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command RCPTTO argument", func(t *testing.T) {
		request, message, errorMessage := "RCPT TO: email@invalid", new(message), configuration.msgInvalidCmdRcpttoArg
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when request includes valid command RCPTTO argument without <> sign", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("RCPT TO: user@example.com"))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})

	t.Run("when request includes valid command RCPTTO argument without <> sign without space", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("RCPT TO:user@example.com"))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})

	t.Run("when request includes valid command RCPTTO argument with <> sign", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("RCPT TO: <user@example.com>"))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})

	t.Run("when request includes valid command RCPTTO argument with <> sign withoyt space", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("RCPT TO:<user@example.com>"))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})
}

func TestHandlerRcpttoRcpttoEmail(t *testing.T) {
	handler := new(handlerRcptto)

	t.Run("when request includes valid email address without <> sign", func(t *testing.T) {
		validEmail := "user@example.com"

		assert.Equal(t, validEmail, handler.rcpttoEmail("RCPT TO: "+validEmail))
	})

	t.Run("when request includes valid email address with <> sign", func(t *testing.T) {
		validEmail := "user@example.com"

		assert.Equal(t, validEmail, handler.rcpttoEmail("RCPT TO: "+"<"+validEmail+">"))
	})

	t.Run("when request includes invalid email address", func(t *testing.T) {
		invalidEmail := "user@invalid"

		assert.Equal(t, EmptyString, handler.rcpttoEmail("RCPT TO: "+invalidEmail))
	})
}

func TestHandlerRcpttoIsBlacklistedEmail(t *testing.T) {
	email := "user@example.com"
	request := "RCPT TO: " + email

	t.Run("when request includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedRcpttoEmails = []string{email}
		errorMessage := configuration.msgRcpttoBlacklistedEmail
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isBlacklistedEmail(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when request not includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isBlacklistedEmail(request))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})
}

func TestHandlerRcpttoIsNotRegisteredEmail(t *testing.T) {
	email := "user@example.com"
	request := "RCPT TO: " + email

	t.Run("when request includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.notRegisteredEmails = []string{email}
		errorMessage := configuration.msgRcpttoNotRegisteredEmail
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isNotRegisteredEmail(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when request not includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isNotRegisteredEmail(request))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})
}

func TestHandlerRcpttoIsInvalidRequest(t *testing.T) {
	configuration := createConfiguration()

	t.Run("when request includes invalid RCPTTO command sequence", func(t *testing.T) {
		request := "MAIL FROM: user@example.com"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdRcpttoSequence
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when request includes invalid RCPTTO command argument", func(t *testing.T) {
		request := "RCPT TO: user@example"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdRcpttoArg
		message.helo, message.mailfrom = true, true
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when request includes blacklisted RCPTTO email", func(t *testing.T) {
		configuration, blacklistedEmail := createConfiguration(), "user@example.com"
		request := "RCPT TO: " + blacklistedEmail
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgRcpttoBlacklistedEmail
		configuration.blacklistedRcpttoEmails = []string{blacklistedEmail}
		message.helo, message.mailfrom = true, true
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when request includes not registered RCPTTO email", func(t *testing.T) {
		configuration, notRegisteredEmail := createConfiguration(), "user@example.com"
		request := "RCPT TO: " + notRegisteredEmail
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgRcpttoNotRegisteredEmail
		configuration.notRegisteredEmails = []string{notRegisteredEmail}
		message.helo, message.mailfrom = true, true
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when valid RCPTTO request", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message := new(sessionMock), new(message)
		message.helo, message.mailfrom = true, true
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})
}
