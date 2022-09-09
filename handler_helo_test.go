package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerHelo(t *testing.T) {
	t.Run("returns new handlerHelo", func(t *testing.T) {
		session, message, configuration := new(session), new(Message), new(configuration)
		handler := newHandlerHelo(session, message, configuration)

		assert.Same(t, session, handler.session)
		assert.Same(t, message, handler.message)
		assert.Same(t, configuration, handler.configuration)
	})
}

func TestHandlerHeloRun(t *testing.T) {
	t.Run("when successful HELO request", func(t *testing.T) {
		request := "HELO example.com"
		session, message, configuration := new(sessionMock), new(Message), createConfiguration()
		receivedMessage := configuration.msgHeloReceived
		handler := newHandlerHelo(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", receivedMessage, configuration.responseDelayHelo).Once().Return(nil)
		handler.run(request)

		assert.True(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, receivedMessage, message.heloResponse)
	})

	t.Run("when failure HELO request, invalid command argument", func(t *testing.T) {
		request := "HELO"
		session, message, configuration := new(sessionMock), new(Message), createConfiguration()
		errorMessage := configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when failure HELO request, blacklisted HELO domain", func(t *testing.T) {
		domainName := "example.com"
		request := "HELO " + domainName
		session, message, configuration := new(sessionMock), new(Message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		errorMessage := configuration.msgHeloBlacklistedDomain
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)
		handler.run(request)

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})
}

func TestHandlerHeloClearMessage(t *testing.T) {
	t.Run("erases all handler message data", func(t *testing.T) {
		notEmptyMessage := createNotEmptyMessage()
		handler, clearedMessage := newHandlerHelo(new(session), notEmptyMessage, new(configuration)), new(Message)
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
		message := new(Message)
		handler := newHandlerHelo(session, message, configuration)
		session.On("writeResponse", response, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.writeResult(true, request, response))
		assert.True(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, response, message.heloResponse)
	})

	t.Run("when failed request received", func(t *testing.T) {
		message, err := new(Message), errors.New(response)
		handler := newHandlerHelo(session, message, configuration)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", response, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.writeResult(false, request, response))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, response, message.heloResponse)
	})
}

func TestHandlerHeloIsInvalidCmdArg(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command HELO argument", func(t *testing.T) {
		request, message, errorMessage := "HELO name.zone42", new(Message), configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes valid command HELO argument", func(t *testing.T) {
		message := new(Message)
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("HELO example.com"))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})

	t.Run("when request includes localhost HELO argument", func(t *testing.T) {
		message := new(Message)
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("HELO localhost"))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})

	t.Run("when request includes valid ip address HELO argument", func(t *testing.T) {
		message := new(Message)
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("HELO 1.2.3.4"))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})

	t.Run("when request includes valid address literal HELO argument", func(t *testing.T) {
		message := new(Message)
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidCmdArg("HELO [1.2.3.4]"))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})

	t.Run("when request includes invalid ip address HELO argument", func(t *testing.T) {
		request, message, errorMessage := "HELO 999.999.999.999", new(Message), configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes invalid address literal HELO argument", func(t *testing.T) {
		request, message, errorMessage := "HELO [999.999.999.999]", new(Message), configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes malformed (left) address literal HELO argument", func(t *testing.T) {
		request, message, errorMessage := "HELO 1.2.3.4]", new(Message), configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request includes malformed (right) address literal HELO argument", func(t *testing.T) {
		request, message, errorMessage := "HELO [1.2.3.4", new(Message), configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})
}

func TestHandlerHeloHeloDomain(t *testing.T) {
	handler := new(handlerHelo)

	t.Run("when request includes valid domain name", func(t *testing.T) {
		validDomainName := "example.com"

		assert.Equal(t, validDomainName, handler.heloDomain("HELO "+validDomainName))
	})

	t.Run("when request includes localhost", func(t *testing.T) {
		localhost := "localhost"

		assert.Equal(t, localhost, handler.heloDomain("HELO "+localhost))
	})

	t.Run("when request not includes valid domain name", func(t *testing.T) {
		invalidDomainName := "name.42"

		assert.Equal(t, emptyString, handler.heloDomain("HELO "+invalidDomainName))
	})
}

func TestHandlerHeloIsBlacklistedDomain(t *testing.T) {
	domainName := "example.com"
	request := "EHLO " + domainName

	t.Run("when request includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(Message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		errorMessage := configuration.msgHeloBlacklistedDomain
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.isBlacklistedDomain(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	t.Run("when request not includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(Message), createConfiguration()
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
		session, message, errorMessage := new(sessionMock), new(Message), configuration.msgInvalidCmdHeloArg
		handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, errorMessage, message.heloResponse)
	})

	heloDomains := []string{"example.com", "localhost", "1.2.3.4", "[1.2.3.4]"}

	for _, heloDomain := range heloDomains {
		t.Run("when valid HELO request", func(t *testing.T) {
			request := "HELO " + heloDomain
			session, message := new(sessionMock), new(Message)
			handler := newHandlerHelo(session, message, configuration)

			assert.False(t, handler.isInvalidRequest(request))
			assert.False(t, message.helo)
			assert.Empty(t, message.heloRequest)
			assert.Empty(t, message.heloResponse)
		})
	}

	for _, blacklistedDomain := range heloDomains {
		t.Run("when request includes blacklisted HELO domain", func(t *testing.T) {
			configuration := createConfiguration()
			request := "HELO " + blacklistedDomain
			session, message, errorMessage := new(sessionMock), new(Message), configuration.msgHeloBlacklistedDomain
			configuration.blacklistedHeloDomains = heloDomains
			handler, err := newHandlerHelo(session, message, configuration), errors.New(errorMessage)
			session.On("addError", err).Once().Return(nil)
			session.On("writeResponse", errorMessage, configuration.responseDelayHelo).Once().Return(nil)

			assert.True(t, handler.isInvalidRequest(request))
			assert.False(t, message.helo)
			assert.Equal(t, request, message.heloRequest)
			assert.Equal(t, errorMessage, message.heloResponse)
		})
	}
}
