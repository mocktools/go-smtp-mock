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
	t.Run("when read request error", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New("some read error")
		session.On("readRequest").Once().Return(EmptyString, err)
		handler.run()

		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid SMTP command", func(t *testing.T) {
		request := "MAILFROM user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmd
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid MAILFROM command sequence", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmdMailfromSequence
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid MAILFROM command argument", func(t *testing.T) {
		request := "MAIL FROM: user@example"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmdMailfromArg
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when fail fast scenario enabled, request includes blacklisted MAILFROM email", func(t *testing.T) {
		email := "user@example.com"
		request := "MAIL FROM: " + email
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast, configuration.blacklistedMailfromEmails = true, []string{email}
		errorMessage := configuration.msgMailfromBlacklistedEmail
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when fail fast scenario enabled, successful MAILFROM request", func(t *testing.T) {
		request := "MAIL FROM: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler := newHandlerMailfrom(session, message, configuration)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", DefaultReceivedMsg).Once().Return(nil)
		handler.run()

		assert.True(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, DefaultReceivedMsg, message.mailfromResponse)
	})

	t.Run("when fail fast scenario disabled, read request error during loop session", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New("some read error")
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(EmptyString, err)
		handler.run()

		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})

	t.Run("when fail fast scenario disabled, no read request errors, 3 failured 1 successful MAILFROM requests", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, validMailfromRequest := newHandlerMailfrom(session, message, configuration), "MAIL FROM: user@domain.com"

		errorMsgInvalidCmd := configuration.msgInvalidCmd
		session.On("readRequest").Once().Return("MAILFROM user@example.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmd)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmd).Once().Return(nil)

		errorMsgInvalidCmdMailfromSequence := configuration.msgInvalidCmdMailfromSequence
		session.On("readRequest").Once().Return("EHLO domain.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdMailfromSequence)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdMailfromSequence).Once().Return(nil)

		errorMsgInvalidCmdMailfromArg := configuration.msgInvalidCmdMailfromArg
		session.On("readRequest").Once().Return("MAIL FROM: user@domain", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdMailfromArg)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdMailfromArg).Once().Return(nil)

		receivedMessage := configuration.msgMailfromReceived
		session.On("clearError").Times(4).Return(nil)
		session.On("readRequest").Once().Return(validMailfromRequest, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.mailfrom)
		assert.Equal(t, validMailfromRequest, message.mailfromRequest)
		assert.Equal(t, receivedMessage, message.mailfromResponse)
	})

	t.Run("when fail fast scenario disabled, request includes blacklisted MAILFROM email", func(t *testing.T) {
		email := "user@example.com"
		request := "MAIL FROM: " + email
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedMailfromEmails = []string{email}
		errorMessage := configuration.msgMailfromBlacklistedEmail
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when fail fast scenario disabled, no read request errors, 4 failured MAILFROM requests", func(t *testing.T) {
		email := "user@example.com"
		requestWithBlacklistedMailfromEmail := "MAIL FROM: " + email
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedMailfromEmails = []string{email}
		handler := newHandlerMailfrom(session, message, configuration)

		errorMsgInvalidCmd := configuration.msgInvalidCmd
		session.On("readRequest").Once().Return("MAILFROM user@example.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmd)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmd).Once().Return(nil)

		errorMsgInvalidCmdMailfromSequence := configuration.msgInvalidCmdMailfromSequence
		session.On("readRequest").Once().Return("EHLO domain.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdMailfromSequence)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdMailfromSequence).Once().Return(nil)

		errorMsgInvalidCmdMailfromArg := configuration.msgInvalidCmdMailfromArg
		session.On("readRequest").Once().Return("MAIL FROM: user@domain", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdMailfromArg)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdMailfromArg).Once().Return(nil)

		errorMsgMailfromBlacklistedEmail := configuration.msgMailfromBlacklistedEmail
		session.On("clearError").Times(4).Return(nil)
		session.On("readRequest").Once().Return(requestWithBlacklistedMailfromEmail, nil)
		session.On("addError", errors.New(errorMsgMailfromBlacklistedEmail)).Once().Return(nil)
		session.On("writeResponse", errorMsgMailfromBlacklistedEmail).Once().Return(nil)
		handler.run()

		assert.False(t, message.mailfrom)
		assert.Equal(t, requestWithBlacklistedMailfromEmail, message.mailfromRequest)
		assert.Equal(t, DefaultQuitMsg, message.mailfromResponse)
	})

	t.Run("when fail fast scenario disabled, successful MAILFROM request", func(t *testing.T) {
		request := "MAIL FROM: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgMailfromReceived
		handler := newHandlerMailfrom(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, receivedMessage, message.mailfromResponse)
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

func TestHandlerMailfromIsInvalidCmd(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request, message, errorMessage := "MAILFROM", new(message), configuration.msgInvalidCmd
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmd(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when request includes valid SMTP command", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("RCPT TO:"))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})
}

func TestHandlerMailfromIsInvalidCmdSequence(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command MAILFROM sequence", func(t *testing.T) {
		request, message, errorMessage := "EHLO", new(message), configuration.msgInvalidCmdMailfromSequence
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when request includes valid command MAILFROM sequence", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("MAIL FROM:"))
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

	t.Run("when request includes valid command MAILFROM argument with <> sign", func(t *testing.T) {
		message := new(message)
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("MAIL FROM: <user@example.com>"))
		assert.False(t, message.mailfrom)
		assert.Empty(t, message.mailfromRequest)
		assert.Empty(t, message.mailfromResponse)
	})
}

func TestHandlerMailfromIsInvalidRequest(t *testing.T) {
	configuration := createConfiguration()

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request := "MAILFROM user@example.com"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmd
		handler, err := newHandlerMailfrom(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.mailfrom)
		assert.Equal(t, request, message.mailfromRequest)
		assert.Equal(t, errorMessage, message.mailfromResponse)
	})

	t.Run("when request includes invalid MAILFROM command sequence", func(t *testing.T) {
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
		handler := newHandlerMailfrom(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
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
