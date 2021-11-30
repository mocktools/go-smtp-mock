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
	t.Run("when successful HELO request", func(t *testing.T) {
		request := "HELO example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		receivedMessage := configuration.msgHeloReceived
		handler := newHandlerHelo(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", receivedMessage).Once().Return(nil)
		handler.run(request)

		assert.True(t, message.helo)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, receivedMessage, message.heloResponse)
	})

	t.Run("when failure HELO request, invalid command argument", func(t *testing.T) {
		request := "HELO"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		errorMessage := configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.helo)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when failure HELO request, blacklisted HELO domain", func(t *testing.T) {
		domainName := "example.com"
		request := "HELO " + domainName
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		errorMessage := configuration.msgHeloBlacklistedDomain
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.helo)
		assert.True(t, message.isCleared())
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})
}

func TestHandlerHeloClearMessage(t *testing.T) {
	t.Run("erases all handler message data, changes cleared status to true", func(t *testing.T) {
		notEmptyMessage := createNotEmptyMessage()
		handler, clearedMessage := newHandlerHelo(new(session), notEmptyMessage, new(configuration)), &message{cleared: true}
		handler.clearMessage()

		assert.Same(t, notEmptyMessage, handler.message)
		assert.Equal(t, clearedMessage, handler.message)

		handler.message.heloRequest = "42"
		handler.clearMessage()

		assert.Equal(t, clearedMessage, handler.message)
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
