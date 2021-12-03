package smtpmock

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("creates new server", func(t *testing.T) {
		configuration := createConfiguration()
		server := newServer(configuration)

		assert.Same(t, configuration, server.configuration)
		assert.Equal(t, new(messages), server.messages)
		assert.Equal(t, newLogger(configuration.logToStdout, configuration.logServerActivity), server.logger)
		assert.Nil(t, server.listener)
		assert.NotNil(t, server.wg)
		assert.Nil(t, server.quit)
		assert.False(t, server.isStarted)
	})
}

func TestServerNewMessage(t *testing.T) {
	t.Run("pushes new message into server.messages, returns this message", func(t *testing.T) {
		server := &Server{messages: new(messages)}
		message, messages := server.newMessage(), server.messages.items

		assert.NotEmpty(t, messages)
		assert.Equal(t, message, messages[0])
	})
}

func TestServerIsInvalidCmd(t *testing.T) {
	availableComands, server := strings.Split("helo,ehlo,mail from:,rcpt to:,data,quit", ","), new(Server)

	for _, validCommand := range availableComands {
		t.Run("when valid command", func(t *testing.T) {
			assert.False(t, server.isInvalidCmd(validCommand))
		})
	}

	t.Run("when invalid command", func(t *testing.T) {
		assert.True(t, server.isInvalidCmd("some invalid command"))
	})
}

func TestServerRecognizeCommand(t *testing.T) {
	t.Run("captures the first word divided by spaces, converts it to upper case", func(t *testing.T) {
		firstWord, secondWord := "first", " command"
		command := firstWord + secondWord

		assert.Equal(t, strings.ToUpper(firstWord), new(Server).recognizeCommand(command))
	})
}

func TestServerAddToWaitGroup(t *testing.T) {
	waitGroup := new(waitGroupMock)
	server := &Server{wg: waitGroup}

	t.Run("increases count of goroutines by one", func(t *testing.T) {
		waitGroup.On("Add", 1).Once().Return(nil)
		server.addToWaitGroup()
	})
}

func TestServerRemoveFromWaitGroup(t *testing.T) {
	waitGroup := new(waitGroupMock)
	server := &Server{wg: waitGroup}

	t.Run("decreases count of goroutines by one", func(t *testing.T) {
		waitGroup.On("Done").Once().Return(nil)
		server.removeFromWaitGroup()
	})
}

func TestServerHandleSession(t *testing.T) {
	t.Run("when complex successful session", func(t *testing.T) {
		session, configuration := &sessionMock{}, createConfiguration()
		server := newServer(configuration)

		session.On("writeResponse", configuration.msgGreeting).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("helo example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgHeloReceived).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("ehlo example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgHeloReceived).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("mail from: receiver@example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgMailfromReceived).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("rcpt to: sender@example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgRcpttoReceived).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("data", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgDataReceived).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("readBytes").Once().Return([]uint8(".some message"), nil)
		session.On("readBytes").Once().Return([]uint8(".\r\n"), nil)
		session.On("writeResponse", configuration.msgMsgReceived).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("quit", nil)
		session.On("writeResponse", configuration.msgQuitCmd).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})

	t.Run("when invalid command, fail fast scenario disabled", func(t *testing.T) {
		session, configuration := &sessionMock{}, createConfiguration()
		server := newServer(configuration)

		session.On("writeResponse", configuration.msgGreeting).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("not implemented command", nil)
		session.On("writeResponse", configuration.msgInvalidCmd).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("quit", nil)
		session.On("writeResponse", configuration.msgQuitCmd).Once().Return(nil)

		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})

	t.Run("when invalid command, session error, fail fast scenario enabled", func(t *testing.T) {
		session, configuration := &sessionMock{}, newConfiguration(ConfigurationAttr{isCmdFailFast: true})
		server, errorMessage := newServer(configuration), configuration.msgInvalidCmdHeloArg

		session.On("writeResponse", configuration.msgGreeting).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("not implemented command", nil)
		session.On("writeResponse", configuration.msgInvalidCmd).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("helo 42", nil)
		session.On("clearError").Once().Return(nil)
		session.On("addError", errors.New(errorMessage)).Once().Return(nil)
		session.On("writeResponse", errorMessage).Once().Return(nil)

		session.On("isErrorFound").Once().Return(true)
		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})

	t.Run("when server quit channel was closed", func(t *testing.T) {
		session, configuration := &sessionMock{}, newConfiguration(ConfigurationAttr{isCmdFailFast: true})
		server := newServer(configuration)
		server.quit = make(chan interface{})
		close(server.quit)

		session.On("writeResponse", configuration.msgGreeting).Once().Return(nil)
		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})

	t.Run("when read request session error", func(t *testing.T) {
		session, configuration := &sessionMock{}, newConfiguration(ConfigurationAttr{isCmdFailFast: true})
		server := newServer(configuration)

		session.On("writeResponse", configuration.msgGreeting).Once().Return(nil)
		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return(emptyString, errors.New("some read request error"))
		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})
}

func TestServerStart(t *testing.T) {
	t.Run("when no errors happens during starting and running the server", func(t *testing.T) {
		configuration := createConfiguration()
		server := newServer(configuration)

		assert.NoError(t, server.Start())
		_ = runMinimalSuccessfulSMTPSession(configuration.hostAddress, configuration.portNumber)
		assert.NotEmpty(t, server.messages)
		assert.NotNil(t, server.quit)
		assert.True(t, server.isStarted)

		_ = server.Stop()
	})

	t.Run("when active server doesn't start current server", func(t *testing.T) {
		server := &Server{isStarted: true}

		assert.EqualError(t, server.Start(), serverStartErrorMsg)
	})

	t.Run("when listener error happens during starting the server doesn't start current server", func(t *testing.T) {
		configuration := createConfiguration()
		server, logger := newServer(configuration), new(loggerMock)
		errorMessage := fmt.Sprintf("%s: %d", serverErrorMsg, configuration.portNumber)
		listener, _ := net.Listen(networkProtocol, serverWithPortNumber(configuration.hostAddress, configuration.portNumber))
		server.logger = logger
		logger.On("error", errorMessage).Once().Return(nil)

		assert.EqualError(t, server.Start(), errorMessage)
		assert.False(t, server.isStarted)
		listener.Close()
	})
}

func TestServerStop(t *testing.T) {
	t.Run("when server active stops current server", func(t *testing.T) {
		logger, listener, waitGroup, quitChannel := new(loggerMock), new(listenerMock), new(waitGroupMock), make(chan interface{})
		server := &Server{logger: logger, listener: listener, wg: waitGroup, quit: quitChannel, isStarted: true}
		listener.On("Close").Once().Return(nil)
		waitGroup.On("Wait").Once().Return(nil)
		logger.On("infoActivity", serverStopMsg).Once().Return(nil)

		assert.NoError(t, server.Stop())
		assert.False(t, server.isStarted)
		_, isChannelOpened := <-server.quit
		assert.False(t, isChannelOpened)
	})

	t.Run("when server is inactive doesn't stop current server", func(t *testing.T) {
		assert.EqualError(t, new(Server).Stop(), serverStopErrorMsg)
	})
}
