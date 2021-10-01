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
	t.Run("when read request error", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler, err := newHandlerRcptto(session, message, configuration), errors.New("some read error")
		session.On("readRequest").Once().Return(EmptyString, err)
		handler.run()

		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid SMTP command", func(t *testing.T) {
		request := "RCPTTO user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmd
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid RCPTTO command sequence", func(t *testing.T) {
		request := "MAIL FROM: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmdRcpttoSequence
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid RCPTTO command argument", func(t *testing.T) {
		request := "RCPT TO: user@example"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmdRcpttoArg
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario enabled, request includes blacklisted RCPTTO email", func(t *testing.T) {
		email := "user@example.com"
		request := "RCPT TO: " + email
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast, configuration.blacklistedRcpttoEmails = true, []string{email}
		errorMessage := configuration.msgRcpttoBlacklistedEmail
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario enabled, request includes not registered RCPTTO email", func(t *testing.T) {
		email := "user@example.com"
		request := "RCPT TO: " + email
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast, configuration.notRegisteredEmails = true, []string{email}
		errorMessage := configuration.msgRcpttoNotRegisteredEmail
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario enabled, successful RCPTTO request", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler := newHandlerRcptto(session, message, configuration)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", DefaultReceivedMsg).Once().Return(nil)
		handler.run()

		assert.True(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, DefaultReceivedMsg, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario disabled, read request error during loop session", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, err := newHandlerRcptto(session, message, configuration), errors.New("some read error")
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(EmptyString, err)
		handler.run()

		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario disabled, no read request errors, 3 failured 1 successful RCPTTO requests", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, validRcpttoRequest := newHandlerRcptto(session, message, configuration), "RCPT TO: user@domain.com"

		errorMsgInvalidCmd := configuration.msgInvalidCmd
		session.On("readRequest").Once().Return("RCPTTO user@example.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmd)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmd).Once().Return(nil)

		errorMsgInvalidCmdRcpttoSequence := configuration.msgInvalidCmdRcpttoSequence
		session.On("readRequest").Once().Return("MAIL FROM: user@domain.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdRcpttoSequence)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdRcpttoSequence).Once().Return(nil)

		errorMsgInvalidCmdRcpttoArg := configuration.msgInvalidCmdRcpttoArg
		session.On("readRequest").Once().Return("RCPT TO: user@domain", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdRcpttoArg)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdRcpttoArg).Once().Return(nil)

		receivedMessage := configuration.msgRcpttoReceived
		session.On("clearError").Times(4).Return(nil)
		session.On("readRequest").Once().Return(validRcpttoRequest, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.rcptto)
		assert.Equal(t, validRcpttoRequest, message.rcpttoRequest)
		assert.Equal(t, receivedMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario disabled, 1 failured blacklisted RCPTTO email request, 1 successful request", func(t *testing.T) {
		email := "user@example.com"
		request, anotherRequest := "RCPT TO: "+email, "RCPT TO: user@another.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedRcpttoEmails = []string{email}
		handler := newHandlerRcptto(session, message, configuration)

		errorMsgRcpttoBlacklistedEmail := configuration.msgRcpttoBlacklistedEmail
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", errors.New(errorMsgRcpttoBlacklistedEmail)).Once().Return(nil)
		session.On("writeResponse", errorMsgRcpttoBlacklistedEmail).Once().Return(nil)

		receivedMessage := configuration.msgRcpttoReceived
		session.On("clearError").Times(2).Return(nil)
		session.On("readRequest").Once().Return(anotherRequest, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.rcptto)
		assert.Equal(t, anotherRequest, message.rcpttoRequest)
		assert.Equal(t, receivedMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario disabled, 1 failured not registered RCPTTO email request, 1 successful request", func(t *testing.T) {
		email := "user@example.com"
		request, anotherRequest := "RCPT TO: "+email, "RCPT TO: user@another.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.notRegisteredEmails = []string{email}
		handler := newHandlerRcptto(session, message, configuration)

		errorMsgRcpttoNotRegisteredEmail := configuration.msgRcpttoNotRegisteredEmail
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", errors.New(errorMsgRcpttoNotRegisteredEmail)).Once().Return(nil)
		session.On("writeResponse", errorMsgRcpttoNotRegisteredEmail).Once().Return(nil)

		receivedMessage := configuration.msgRcpttoReceived
		session.On("clearError").Times(2).Return(nil)
		session.On("readRequest").Once().Return(anotherRequest, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.rcptto)
		assert.Equal(t, anotherRequest, message.rcpttoRequest)
		assert.Equal(t, receivedMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario disabled, no read request errors, 5 failured RCPTTO requests, 1 successful request", func(t *testing.T) {
		blacklistedEmail, notRegisteredEmail := "blacklisted@example.com", "not_existent@example.com"
		requestBlacklistedRcpttoEmail, requestWithNotRegiteredRcpttoEmail := "RCPT TO: "+blacklistedEmail, "RCPT TO: "+notRegisteredEmail
		successfulRequest := "RCPT TO: successful@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedRcpttoEmails, configuration.notRegisteredEmails = []string{blacklistedEmail}, []string{notRegisteredEmail}
		handler := newHandlerRcptto(session, message, configuration)

		errorMsgInvalidCmd := configuration.msgInvalidCmd
		session.On("readRequest").Once().Return("RCPTTO user@example.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmd)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmd).Once().Return(nil)

		errorMsgInvalidCmdRcpttoSequence := configuration.msgInvalidCmdRcpttoSequence
		session.On("readRequest").Once().Return("MAIL FROM: user@domain.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdRcpttoSequence)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdRcpttoSequence).Once().Return(nil)

		errorMsgInvalidCmdRcpttoArg := configuration.msgInvalidCmdRcpttoArg
		session.On("readRequest").Once().Return("RCPT TO: user@domain", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdRcpttoArg)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdRcpttoArg).Once().Return(nil)

		errorMsgRcpttoBlacklistedEmail := configuration.msgRcpttoBlacklistedEmail
		session.On("readRequest").Once().Return(requestBlacklistedRcpttoEmail, nil)
		session.On("addError", errors.New(errorMsgRcpttoBlacklistedEmail)).Once().Return(nil)
		session.On("writeResponse", errorMsgRcpttoBlacklistedEmail).Once().Return(nil)

		errorMsgRcpttoNotRegisteredEmail := configuration.msgRcpttoNotRegisteredEmail
		session.On("readRequest").Once().Return(requestWithNotRegiteredRcpttoEmail, nil)
		session.On("addError", errors.New(errorMsgRcpttoNotRegisteredEmail)).Once().Return(nil)
		session.On("writeResponse", errorMsgRcpttoNotRegisteredEmail).Once().Return(nil)

		receivedMessage := configuration.msgRcpttoReceived
		session.On("clearError").Times(6).Return(nil)
		session.On("readRequest").Once().Return(successfulRequest, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.rcptto)
		assert.Equal(t, successfulRequest, message.rcpttoRequest)
		assert.Equal(t, receivedMessage, message.rcpttoResponse)
	})

	t.Run("when fail fast scenario disabled, successful RCPTTO request", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgRcpttoReceived
		handler := newHandlerRcptto(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, receivedMessage, message.rcpttoResponse)
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

func TestHandlerRcpttoIsInvalidCmd(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request, message, errorMessage := "RCPTTO", new(message), configuration.msgInvalidCmd
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmd(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when request includes valid SMTP command", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("RCPT TO:"))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})
}

func TestHandlerRcpttoIsInvalidCmdSequence(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command RCPTTO sequence", func(t *testing.T) {
		request, message, errorMessage := "MAIL FROM:", new(message), configuration.msgInvalidCmdRcpttoSequence
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

	t.Run("when request includes valid command MAILFROM sequence", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("RCPT TO:"))
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

	t.Run("when request includes valid command RCPTTO argument with <> sign", func(t *testing.T) {
		message := new(message)
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("RCPT TO: <user@example.com>"))
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

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request := "RCPTTO user@example.com"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmd
		handler, err := newHandlerRcptto(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.rcptto)
		assert.Equal(t, request, message.rcpttoRequest)
		assert.Equal(t, errorMessage, message.rcpttoResponse)
	})

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
		handler := newHandlerRcptto(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
		assert.False(t, message.rcptto)
		assert.Empty(t, message.rcpttoRequest)
		assert.Empty(t, message.rcpttoResponse)
	})
}
