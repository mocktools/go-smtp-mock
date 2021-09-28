package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerHelo(t *testing.T) {
	t.Run("returns new handlerHelo", func(t *testing.T) {
		session, message, configuration := new(session), new(message), new(configuration)
		handler := newHandlerHelo(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerHeloRun(t *testing.T) {
	t.Run("when read request error", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler, err := newHandlerHelo(session, message, configuration), errors.New("some read error")
		session.On("readRequest").Once().Return(EmptyString, err)
		handler.run()

		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid SMTP command", func(t *testing.T) {
		request := "HI example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmd
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid HELO command sequence", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmdHeloSequence
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid HELO command argument", func(t *testing.T) {
		request := "HELO user@example"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		errorMessage := configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, request includes blacklisted HELO domain", func(t *testing.T) {
		domainName := "example.com"
		request := "HELO " + domainName
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast, configuration.blacklistedHeloDomains = true, []string{domainName}
		errorMessage := configuration.msgHeloBlacklistedDomain
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, successful HELO request", func(t *testing.T) {
		request := "HELO example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		receivedMessage := configuration.msgHeloReceived
		handler := newHandlerHelo(session, message, configuration)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, receivedMessage, message.heloResponse)
	})

	t.Run("when fail fast scenario disabled, read request error during loop session", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, err := newHandlerHelo(session, message, configuration), errors.New("some read error")
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(EmptyString, err)
		handler.run()

		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})

	t.Run("when fail fast scenario disabled, no read request errors, 3 failured 1 successful HELO requests", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, validHeloRequest := newHandlerHelo(session, message, configuration), "HELO domain.com"

		errorMsgInvalidCmd := configuration.msgInvalidCmd
		session.On("readRequest").Once().Return("HI example.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmd)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmd).Once().Return(nil)

		errorMsgInvalidCmdMailfromSequence := configuration.msgInvalidCmdHeloSequence
		session.On("readRequest").Once().Return("MAIL FROM: user@domain.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdMailfromSequence)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdMailfromSequence).Once().Return(nil)

		errorMsgInvalidCmdMailfromArg := configuration.msgInvalidCmdHeloArg
		session.On("readRequest").Once().Return("HELO user@domain", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdMailfromArg)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdMailfromArg).Once().Return(nil)

		receivedMessage := configuration.msgHeloReceived
		session.On("clearError").Times(4).Return(nil)
		session.On("readRequest").Once().Return(validHeloRequest, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.helo)
		assert.Equal(t, validHeloRequest, message.heloRequest)
		assert.Equal(t, receivedMessage, message.heloResponse)
	})

	t.Run("when fail fast scenario disabled, 1 failured blacklisted HELO domain request, 1 successful request", func(t *testing.T) {
		domainName := "example.com"
		request, anotherRequest := "EHLO "+domainName, "HELO another.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		handler := newHandlerHelo(session, message, configuration)

		errorMsgHeloBlacklistedDomain := configuration.msgHeloBlacklistedDomain
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", errors.New(errorMsgHeloBlacklistedDomain)).Once().Return(nil)
		session.On("writeResponse", errorMsgHeloBlacklistedDomain).Once().Return(nil)

		receivedMessage := configuration.msgHeloReceived
		session.On("clearError").Times(2).Return(nil)
		session.On("readRequest").Once().Return(anotherRequest, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.helo)
		assert.Equal(t, anotherRequest, message.heloRequest)
		assert.Equal(t, receivedMessage, message.heloResponse)
	})

	t.Run("when fail fast scenario disabled, no read request errors, 4 failured, 1 successful HELO requests", func(t *testing.T) {
		domainName := "example.com"
		requestWithBlacklistedHeloDomain, anotherRequest := "EHLO "+domainName, "HELO another.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		handler := newHandlerHelo(session, message, configuration)

		errorMsgInvalidCmd := configuration.msgInvalidCmd
		session.On("readRequest").Once().Return("HI example.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmd)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmd).Once().Return(nil)

		errorMsgInvalidCmdMailfromSequence := configuration.msgInvalidCmdHeloSequence
		session.On("readRequest").Once().Return("MAIL FROM: user@domain.com", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdMailfromSequence)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdMailfromSequence).Once().Return(nil)

		errorMsgInvalidCmdMailfromArg := configuration.msgInvalidCmdHeloArg
		session.On("readRequest").Once().Return("HELO user@domain", nil)
		session.On("addError", errors.New(errorMsgInvalidCmdMailfromArg)).Once().Return(nil)
		session.On("writeResponse", errorMsgInvalidCmdMailfromArg).Once().Return(nil)

		errorMsgHeloBlacklistedDomain := configuration.msgHeloBlacklistedDomain
		session.On("clearError").Times(4).Return(nil)
		session.On("readRequest").Once().Return(requestWithBlacklistedHeloDomain, nil)
		session.On("addError", errors.New(errorMsgHeloBlacklistedDomain)).Once().Return(nil)
		session.On("writeResponse", errorMsgHeloBlacklistedDomain).Once().Return(nil)

		receivedMessage := configuration.msgHeloReceived
		session.On("clearError").Times(5).Return(nil)
		session.On("readRequest").Once().Return(anotherRequest, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.helo)
		assert.Equal(t, anotherRequest, message.heloRequest)
		assert.Equal(t, receivedMessage, message.heloResponse)
	})

	t.Run("when fail fast scenario disabled, successful HELO request", func(t *testing.T) {
		request := "HELO example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgHeloReceived
		handler := newHandlerHelo(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run()

		assert.True(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, receivedMessage, message.heloResponse)
	})
}

func TestHandlerHeloWriteResult(t *testing.T) {
	request, response := "request context", "response context"
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when successful request received", func(t *testing.T) {
		message := new(message)
		handler := newHandlerHelo(session, message, configuration)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, response, message.heloResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(message), errors.New(response)
		handler := newHandlerHelo(session, message, configuration)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, response, message.heloResponse)
	})
}

func TestHandlerHeloIsInvalidCmd(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request, message, errorMessage := "HI", new(message), configuration.msgInvalidCmd
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmd(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes valid SMTP command", func(t *testing.T) {
		message := new(message)
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("RCPT TO:"))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})
}

func TestHandlerHeloIsInvalidCmdSequence(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command HELO sequence", func(t *testing.T) {
		request, message, errorMessage := "MAIL FROM:", new(message), configuration.msgInvalidCmdHeloSequence
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes valid command HELO sequence", func(t *testing.T) {
		message := new(message)
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidCmd("EHLO"))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})
}

func TestHandlerHeloIsInvalidCmdArg(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command HELO argument", func(t *testing.T) {
		request, message, errorMessage := "HELO name.zone42", new(message), configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes valid command HELO argument", func(t *testing.T) {
		message := new(message)
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("HELO example.com"))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})
}

func TestHandlerHeloHeloDomain(t *testing.T) {
	handler := new(handlerHelo)

	t.Run("when request includes valid domain name", func(t *testing.T) {
		validDomainName := "example.com"

		assert.Equal(t, validDomainName, handler.heloDomain("HELO "+validDomainName))
	})

	t.Run("when request not includes valid domain name", func(t *testing.T) {
		invalidDomainName := "name.42"

		assert.Equal(t, EmptyString, handler.heloDomain("HELO "+invalidDomainName))
	})
}

func TestHandlerHeloIsBlacklistedDomain(t *testing.T) {
	domainName := "example.com"
	request := "EHLO " + domainName

	t.Run("when request includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		errorMessage := configuration.msgQuitCmd
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isBlacklistedDomain(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request not includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isBlacklistedDomain(request))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})
}

func TestHandlerHeloIsInvalidRequest(t *testing.T) {
	configuration := createConfiguration()

	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request := "HI example.com"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmd
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes invalid HELO command sequence", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdHeloSequence
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes invalid HELO command argument", func(t *testing.T) {
		request := "HELO user@example"
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes blacklisted HELO domain", func(t *testing.T) {
		configuration, blacklistedDomain := createConfiguration(), "example.com"
		request := "HELO " + blacklistedDomain
		session, message, errorMessage := new(sessionMock), new(message), configuration.msgHeloBlacklistedDomain
		configuration.blacklistedHeloDomains = []string{blacklistedDomain}
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when valid HELO request", func(t *testing.T) {
		request := "HELO example.com"
		session, message := new(sessionMock), new(message)
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})
}
