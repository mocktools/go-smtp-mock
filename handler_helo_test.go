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

func TestHandlerRun(t *testing.T) {
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
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultInvalidCmdMsg)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdMsg).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdMsg, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid HELO command sequence", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultInvalidCmdHeloSequenceMsg)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloSequenceMsg).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdHeloSequenceMsg, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, request includes invalid HELO command argument", func(t *testing.T) {
		request := "HELO user@example"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultInvalidCmdHeloArgMsg)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloArgMsg).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdHeloArgMsg, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, request includes blacklisted HELO domain", func(t *testing.T) {
		domainName := "example.com"
		request := "HELO " + domainName
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast, configuration.blacklistedHeloDomains = true, []string{domainName}
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultQuitMsg)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultQuitMsg).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultQuitMsg, message.heloResponse)
	})

	t.Run("when fail fast scenario enabled, successful HELO request", func(t *testing.T) {
		request := "HELO example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.isCmdFailFast = true
		handler := newHandlerHelo(session, message, configuration)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", DefaultReceivedMsg).Once().Return(nil)
		handler.run()

		assert.True(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultReceivedMsg, message.heloResponse)
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

		session.On("readRequest").Once().Return("HI example.com", nil)
		session.On("addError", errors.New(DefaultInvalidCmdMsg)).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdMsg).Once().Return(nil)

		session.On("readRequest").Once().Return("MAIL FROM: user@domain.com", nil)
		session.On("addError", errors.New(DefaultInvalidCmdHeloSequenceMsg)).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloSequenceMsg).Once().Return(nil)

		session.On("readRequest").Once().Return("HELO user@domain", nil)
		session.On("addError", errors.New(DefaultInvalidCmdHeloArgMsg)).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloArgMsg).Once().Return(nil)

		session.On("clearError").Times(4).Return(nil)
		session.On("readRequest").Once().Return(validHeloRequest, nil)
		session.On("writeResponse", DefaultReceivedMsg).Once().Return(nil)
		handler.run()

		assert.True(t, message.helo)
		assert.Equal(t, validHeloRequest, message.heloRequest)
		assert.Equal(t, DefaultReceivedMsg, message.heloResponse)
	})

	t.Run("when fail fast scenario disabled, request includes blacklisted HELO domain", func(t *testing.T) {
		domainName := "example.com"
		request := "EHLO " + domainName
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultQuitMsg)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultQuitMsg).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultQuitMsg, message.heloResponse)
	})

	t.Run("when fail fast scenario disabled, no read request errors, 4 failured HELO requests", func(t *testing.T) {
		domainName := "example.com"
		requestWithBlacklistedHeloDomain := "EHLO " + domainName
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		handler := newHandlerHelo(session, message, configuration)

		session.On("readRequest").Once().Return("HI example.com", nil)
		session.On("addError", errors.New(DefaultInvalidCmdMsg)).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdMsg).Once().Return(nil)

		session.On("readRequest").Once().Return("MAIL FROM: user@domain.com", nil)
		session.On("addError", errors.New(DefaultInvalidCmdHeloSequenceMsg)).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloSequenceMsg).Once().Return(nil)

		session.On("readRequest").Once().Return("HELO user@domain", nil)
		session.On("addError", errors.New(DefaultInvalidCmdHeloArgMsg)).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloArgMsg).Once().Return(nil)

		session.On("clearError").Times(4).Return(nil)
		session.On("readRequest").Once().Return(requestWithBlacklistedHeloDomain, nil)
		session.On("addError", errors.New(DefaultQuitMsg)).Once().Return(nil)
		session.On("writeResponse", DefaultQuitMsg).Once().Return(nil)
		handler.run()

		assert.False(t, message.helo)
		assert.Equal(t, requestWithBlacklistedHeloDomain, message.heloRequest)
		assert.Equal(t, DefaultQuitMsg, message.heloResponse)
	})

	t.Run("when fail fast scenario disabled, successful HELO request", func(t *testing.T) {
		request := "HELO example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler := newHandlerHelo(session, message, configuration)
		session.On("clearError").Once().Return(nil)
		session.On("readRequest").Once().Return(request, nil)
		session.On("writeResponse", DefaultReceivedMsg).Once().Return(nil)
		handler.run()

		assert.True(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultReceivedMsg, message.heloResponse)
	})
}

func TestHandlerHeloIsInvalidCmdSequence(t *testing.T) {
	configuration, session := createConfiguration(), &sessionMock{}

	t.Run("when request includes invalid command HELO sequence", func(t *testing.T) {
		request, message := "MAIL FROM:", new(message)
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultInvalidCmdHeloSequenceMsg)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloSequenceMsg).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdSequence(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdHeloSequenceMsg, message.heloResponse)
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
		request, message := "HELO name.zone42", new(message)
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultInvalidCmdHeloArgMsg)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloArgMsg).Once().Return(nil)

		assert.True(t, handler.isInvalidCmdArg(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdHeloArgMsg, message.heloResponse)
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

func TestHandlerheloDomain(t *testing.T) {
	handler := new(handlerHelo)

	t.Run("when request includes valid domain name", func(t *testing.T) {
		validDomainName := "example.com"

		assert.Equal(t, validDomainName, handler.heloDomain("HELO "+validDomainName))
	})

	t.Run("when request not includes valid domain name", func(t *testing.T) {
		invalidDomainName := "name.42"

		assert.Equal(t, EmptyString, handler.heloDomain("HELO "+invalidDomainName))
	})

	t.Run("when request includes partial match to valid domain name", func(t *testing.T) {
		validDomainName := "example.com"

		assert.Equal(t, validDomainName, handler.heloDomain("HELO "+validDomainName+"42"))
	})
}

func TestHandlerIsBlacklistedDomain(t *testing.T) {
	domainName := "example.com"
	request := "EHLO " + domainName

	t.Run("when request includes blacklisted domain name", func(t *testing.T) {
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		configuration.blacklistedHeloDomains = []string{domainName}
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultQuitMsg)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultQuitMsg).Once().Return(nil)

		assert.True(t, handler.isBlacklistedDomain(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultQuitMsg, message.heloResponse)
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

func TestHandlerIsInvalidRequest(t *testing.T) {
	t.Run("when request includes invalid SMTP command", func(t *testing.T) {
		request := "HI example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultInvalidCmdMsg)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdMsg).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdMsg, message.heloResponse)
	})

	t.Run("when request includes invalid HELO command sequence", func(t *testing.T) {
		request := "RCPT TO: user@example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultInvalidCmdHeloSequenceMsg)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloSequenceMsg).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdHeloSequenceMsg, message.heloResponse)
	})

	t.Run("when request includes invalid HELO command argument", func(t *testing.T) {
		request := "HELO user@example"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler, err := newHandlerHelo(session, message, configuration), errors.New(DefaultInvalidCmdHeloArgMsg)
		session.On("addError", err).Once().Return(nil)
		session.On("writeResponse", DefaultInvalidCmdHeloArgMsg).Once().Return(nil)

		assert.True(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Equal(t, request, message.heloRequest)
		assert.Equal(t, DefaultInvalidCmdHeloArgMsg, message.heloResponse)
	})

	t.Run("when valid HELO request", func(t *testing.T) {
		request := "HELO example.com"
		session, message, configuration := new(sessionMock), new(message), createConfiguration()
		handler := newHandlerHelo(session, message, configuration)

		assert.False(t, handler.isInvalidRequest(request))
		assert.False(t, message.helo)
		assert.Empty(t, message.heloRequest)
		assert.Empty(t, message.heloResponse)
	})
}
