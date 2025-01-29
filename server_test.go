package smtpmock

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("creates new server", func(t *testing.T) {
		configuration := createConfiguration()
		server := newServer(configuration)

		assert.Same(t, configuration, server.configuration)
		assert.Equal(t, new(messages), server.messages)
		assert.Equal(t, slog.Default(), server.logger)
		assert.Nil(t, server.listener)
		assert.NotNil(t, server.wg)
		assert.Nil(t, server.quit)
		assert.False(t, server.isStarted())
		assert.Equal(t, 0, server.PortNumber())
	})
}

func TestServerStart(t *testing.T) {
	t.Run("when no errors happens during starting and running the server with default port", func(t *testing.T) {
		configuration := createConfiguration()
		server := newServer(configuration)

		assert.NoError(t, server.Start())
		_ = runSuccessfulSMTPSession(configuration.hostAddress, server.PortNumber(), false, 0)
		assert.NotNil(t, server.messages)
		assert.NotNil(t, server.quit)
		assert.NotNil(t, server.quitTimeout)
		assert.True(t, server.isStarted())
		assert.Greater(t, server.PortNumber(), 0)

		_ = server.Stop()
	})

	t.Run("when no errors happens during starting and running the server with custom port", func(t *testing.T) {
		configuration, portNumber := createConfiguration(), 2525
		configuration.portNumber = portNumber
		server := newServer(configuration)

		assert.NoError(t, server.Start())
		_ = runSuccessfulSMTPSession(configuration.hostAddress, portNumber, false, 0)
		assert.NotNil(t, server.messages)
		assert.NotNil(t, server.quit)
		assert.NotNil(t, server.quitTimeout)
		assert.True(t, server.isStarted())
		assert.Equal(t, portNumber, server.PortNumber())

		_ = server.Stop()
	})

	t.Run("when active server doesn't start current server", func(t *testing.T) {
		server := &Server{started: true}

		assert.EqualError(t, server.Start(), serverStartErrorMsg)
		assert.Equal(t, 0, server.PortNumber())
	})

	t.Run("when listener error happens during starting the server doesn't start current server", func(t *testing.T) {
		configuration := createConfiguration()
		server, logger := newServer(configuration), slog.Default()
		listener, _ := net.Listen(networkProtocol, emptyString)
		portNumber := listener.Addr().(*net.TCPAddr).Port
		errorMessage := fmt.Sprintf("%s: %d", serverErrorMsg, portNumber)
		configuration.portNumber, server.logger = portNumber, logger

		assert.EqualError(t, server.Start(), errorMessage)
		assert.False(t, server.isStarted())
		assert.Equal(t, 0, server.PortNumber())
		listener.Close()
	})
}

func TestServerStop(t *testing.T) {
	t.Run("when server active stops current server, graceful shutdown case", func(t *testing.T) {
		logger, listener, waitGroup, quitChannel := slog.Default(), new(listenerMock), new(waitGroupMock), make(chan interface{})
		server := &Server{
			configuration: createConfiguration(),
			logger:        logger,
			listener:      listener,
			wg:            waitGroup,
			quit:          quitChannel,
			started:       true,
			quitTimeout:   make(chan interface{}),
		}
		listener.On("Close").Once().Return(nil)
		waitGroup.On("Wait").Once().Return(nil)

		assert.NoError(t, server.Stop())
		assert.False(t, server.isStarted())
		_, isChannelOpened := <-server.quit
		assert.False(t, isChannelOpened)
	})

	t.Run("when server active stops current server, force shutdown case", func(t *testing.T) {
		logger, listener, waitGroup, quitChannel := slog.Default(), new(listenerMock), new(waitGroupMock), make(chan interface{})
		server := &Server{
			configuration: createConfiguration(),
			logger:        logger,
			listener:      listener,
			wg:            waitGroup,
			quit:          quitChannel,
			started:       true,
		}
		listener.On("Close").Once().Return(nil)
		waitGroup.On("Wait").Once().Return(nil)

		assert.NoError(t, server.Stop())
		assert.False(t, server.isStarted())
		_, isChannelOpened := <-server.quit
		assert.False(t, isChannelOpened)
	})

	t.Run("when server is inactive doesn't stop current server", func(t *testing.T) {
		assert.EqualError(t, new(Server).Stop(), serverStopErrorMsg)
	})
}

func TestServerMessages(t *testing.T) {
	configuration := createConfiguration()

	t.Run("when there are no messages on the server", func(t *testing.T) {
		server := newServer(configuration)

		assert.Empty(t, server.Messages())
	})

	t.Run("when there are messages on the server", func(t *testing.T) {
		server := newServer(configuration)
		message := new(Message)
		server.messages.append(message)

		assert.NotEmpty(t, server.Messages())
	})

	t.Run("message data are identical", func(t *testing.T) {
		server := newServer(configuration)

		server.messages.RLock()
		assert.Empty(t, server.messages.items)
		assert.Empty(t, server.Messages())
		assert.NotSame(t, server.messages.items, server.Messages())
		server.messages.RUnlock()

		message := new(Message)
		server.messages.append(message)

		server.messages.RLock()
		assert.Equal(t, []*Message{message}, server.messages.items)
		assert.Equal(t, []Message{*message}, server.Messages())
		assert.NotSame(t, server.messages.items, server.Messages())
		server.messages.RUnlock()
	})
}

func TestServerWaitForMessages(t *testing.T) {
	timeout := 1 * time.Millisecond

	t.Run("when expected number of messages is received without timeout", func(t *testing.T) {
		server, message := newServer(createConfiguration()), new(Message)
		server.messages.append(message)
		messages, err := server.WaitForMessages(len(server.messages.copy()), timeout)

		assert.Equal(t, []Message{*message}, messages)
		assert.NoError(t, err)
	})

	t.Run("when timeout occurs before receiving expected number of messages", func(t *testing.T) {
		server := newServer(createConfiguration())
		messages, err := server.WaitForMessages(1, timeout)

		assert.EqualError(t, err, fmt.Sprintf("timeout waiting for %d messages, got %d", 1, 0))
		assert.Empty(t, messages)
	})
}

func TestServerMessagesAndPurge(t *testing.T) {
	t.Run("returns empty messages after purge", func(t *testing.T) {
		server, message := newServer(createConfiguration()), new(Message)
		server.messages.append(message)

		assert.NotEmpty(t, server.Messages())
		assert.NotEmpty(t, server.MessagesAndPurge())
		assert.Empty(t, server.Messages())
	})
}

func TestServerWaitForMessagesAndPurge(t *testing.T) {
	timeout := 1 * time.Millisecond

	t.Run("when expected number of messages is received without timeout", func(t *testing.T) {
		server, message := newServer(createConfiguration()), new(Message)
		server.messages.append(message)
		messages, err := server.WaitForMessagesAndPurge(len(server.messages.copy()), timeout)

		assert.Equal(t, []Message{*message}, messages)
		assert.NoError(t, err)
		assert.Empty(t, server.Messages())
	})

	t.Run("when timeout occurs before receiving expected number of messages", func(t *testing.T) {
		server := newServer(createConfiguration())
		messages, err := server.WaitForMessagesAndPurge(1, timeout)

		assert.EqualError(t, err, fmt.Sprintf("timeout waiting for %d messages, got %d", 1, 0))
		assert.Empty(t, messages)
	})
}

func TestServerPortNumber(t *testing.T) {
	t.Run("returns server port number", func(t *testing.T) {
		portNumber := 2525
		server := &Server{portNumber: portNumber}

		assert.Equal(t, portNumber, server.PortNumber())
	})
}

func TestServerFetchMessages(t *testing.T) {
	timeout := 1 * time.Millisecond

	t.Run("when expected number of messages is received without timeout", func(t *testing.T) {
		server, message := newServer(createConfiguration()), new(Message)
		server.messages.append(message)
		messages, err := server.fetchMessages(len(server.messages.copy()), timeout, false)

		assert.Equal(t, []Message{*message}, messages)
		assert.NoError(t, err)
		assert.NotEmpty(t, server.Messages())
	})

	t.Run("when expected number of messages is received with purging", func(t *testing.T) {
		server, message := newServer(createConfiguration()), new(Message)
		server.messages.append(message)
		messages, err := server.fetchMessages(len(server.messages.copy()), timeout, true)

		assert.Equal(t, []Message{*message}, messages)
		assert.NoError(t, err)
		assert.Empty(t, server.Messages())
	})

	t.Run("when timeout occurs before receiving expected number of messages", func(t *testing.T) {
		server := newServer(createConfiguration())
		messages, err := server.fetchMessages(1, timeout, false)

		assert.EqualError(t, err, fmt.Sprintf("timeout waiting for %d messages, got %d", 1, 0))
		assert.Empty(t, messages)
	})
}

func TestServerIsStarted(t *testing.T) {
	t.Run("returns current server started-flag status", func(t *testing.T) {
		server := &Server{started: true}

		assert.True(t, server.isStarted())
	})
}

func TestServerSetListener(t *testing.T) {
	t.Run("sets server listener", func(t *testing.T) {
		server := new(Server)
		listener, _ := net.Listen("tcp", "localhost:2526")
		server.setListener(listener)

		assert.Equal(t, listener, server.listener)
	})
}

func TestServerSetPortNumber(t *testing.T) {
	t.Run("sets server listener", func(t *testing.T) {
		server, portNumber := new(Server), 2525
		server.setPortNumber(portNumber)

		assert.Equal(t, portNumber, server.PortNumber())
	})
}

func TestServerStartFlag(t *testing.T) {
	t.Run("sets server started-flag status to true", func(t *testing.T) {
		server := new(Server)
		server.start()

		assert.True(t, server.isStarted())
	})
}

func TestServerStopFlag(t *testing.T) {
	t.Run("sets server started-flag status to false", func(t *testing.T) {
		server := &Server{started: true}
		server.stop()

		assert.False(t, server.isStarted())
	})
}

func TestServerNewMessageWithHeloContext(t *testing.T) {
	t.Run("pushes new message into server.messages with helo context from other message, returns this message", func(t *testing.T) {
		server := &Server{messages: new(messages)}
		message, heloRequest, heloResponse, helo := new(Message), "heloRequest", "heloResponse", true
		message.heloRequest, message.heloResponse, message.helo = heloRequest, heloResponse, helo
		newMessage := server.newMessageWithHeloContext(message)

		server.messages.RLock()
		messages := server.messages.items
		assert.Equal(t, heloRequest, newMessage.heloRequest)
		assert.Equal(t, heloResponse, newMessage.heloResponse)
		assert.Equal(t, helo, newMessage.helo)
		assert.Equal(t, newMessage, messages[0])
		assert.Equal(t, 1, len(messages))
		server.messages.RUnlock()
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

	t.Run("increases count of goroutines by one", func(*testing.T) {
		waitGroup.On("Add", 1).Once().Return(nil)
		server.addToWaitGroup()
	})
}

func TestServerRemoveFromWaitGroup(t *testing.T) {
	waitGroup := new(waitGroupMock)
	server := &Server{wg: waitGroup}

	t.Run("decreases count of goroutines by one", func(*testing.T) {
		waitGroup.On("Done").Once().Return(nil)
		server.removeFromWaitGroup()
	})
}

func TestServerIsAbleToEndSession(t *testing.T) {
	t.Run("when quit command has been sent", func(t *testing.T) {
		server, message, session := newServer(createConfiguration()), &Message{quitSent: true}, new(session)
		server.messages.append(message)

		assert.True(t, server.isAbleToEndSession(message, session))
	})

	t.Run("when quit command has not been sent, error has been found, fail fast scenario has been enabled", func(t *testing.T) {
		server, message, session := newServer(createConfiguration()), new(Message), new(session)
		server.messages.append(message)
		session.err = errors.New("some error")
		server.configuration.isCmdFailFast = true

		assert.True(t, server.isAbleToEndSession(message, session))
	})

	t.Run("when quit command has not been sent, no errors", func(t *testing.T) {
		server, message, session := newServer(createConfiguration()), new(Message), new(session)
		server.messages.append(message)

		assert.False(t, server.isAbleToEndSession(message, session))
	})

	t.Run("when quit command has not been sent, error has been found, fail fast scenario has not been enabled", func(t *testing.T) {
		server, message, session := newServer(createConfiguration()), new(Message), new(session)
		server.messages.append(message)
		session.err = errors.New("some error")

		assert.False(t, server.isAbleToEndSession(message, session))
	})
}

func TestServerHandleSession(t *testing.T) {
	t.Run("when complex successful session, multiple message receiving scenario disabled", func(t *testing.T) {
		session, configuration := &sessionMock{}, createConfiguration()
		server := newServer(configuration)

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("helo example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgHeloReceived, configuration.responseDelayHelo).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("noop", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgNoopReceived, configuration.responseDelayNoop).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("ehlo example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgHeloReceived, configuration.responseDelayHelo).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("rset", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgRsetReceived, configuration.responseDelayRset).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("mail from: receiver@example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgMailfromReceived, configuration.responseDelayMailfrom).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("rcpt to: sender@example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgRcpttoReceived, configuration.responseDelayRcptto).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("data", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgDataReceived, configuration.responseDelayData).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("readBytes").Once().Return([]uint8(".some message"), nil)
		session.On("readBytes").Once().Return([]uint8(".\r\n"), nil)
		session.On("writeResponse", configuration.msgMsgReceived, configuration.responseDelayMessage).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("quit", nil)
		session.On("writeResponse", configuration.msgQuitCmd, configuration.responseDelayQuit).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("finish").Once().Return(nil)

		server.handleSession(session)
		assert.Equal(t, 1, len(server.Messages()))
	})

	t.Run("when complex successful session, multiple message receiving scenario enabled", func(t *testing.T) {
		session, configuration := &sessionMock{}, createConfiguration()
		configuration.multipleMessageReceiving = true
		server := newServer(configuration)

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("helo example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgHeloReceived, configuration.responseDelayHelo).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("noop", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgNoopReceived, configuration.responseDelayNoop).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("ehlo example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgHeloReceived, configuration.responseDelayHelo).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("mail from: receiver@example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgMailfromReceived, configuration.responseDelayMailfrom).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("rcpt to: sender1@example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgRcpttoReceived, configuration.responseDelayRcptto).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("data", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgDataReceived, configuration.responseDelayData).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("readBytes").Once().Return([]uint8(".some message"), nil)
		session.On("readBytes").Once().Return([]uint8(".\r\n"), nil)
		session.On("writeResponse", configuration.msgMsgReceived, configuration.responseDelayMessage).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("rset", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgRsetReceived, configuration.responseDelayRset).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("mail from: receiver@example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgMailfromReceived, configuration.responseDelayMailfrom).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("rcpt to: sender1@example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgRcpttoReceived, configuration.responseDelayRcptto).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("data", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgDataReceived, configuration.responseDelayData).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("readBytes").Once().Return([]uint8(".some message"), nil)
		session.On("readBytes").Once().Return([]uint8(".\r\n"), nil)
		session.On("writeResponse", configuration.msgMsgReceived, configuration.responseDelayMessage).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("quit", nil)
		session.On("writeResponse", configuration.msgQuitCmd, configuration.responseDelayQuit).Once().Return(nil)
		session.On("isErrorFound").Once().Return(false)

		session.On("finish").Once().Return(nil)

		server.handleSession(session)
		assert.Equal(t, 2, len(server.Messages()))
	})

	t.Run("when invalid command, fail fast scenario disabled", func(*testing.T) {
		session, configuration := &sessionMock{}, createConfiguration()
		server := newServer(configuration)

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("not implemented command", nil)
		session.On("writeResponse", configuration.msgInvalidCmd, defaultSessionResponseDelay).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("quit", nil)
		session.On("writeResponse", configuration.msgQuitCmd, configuration.responseDelayQuit).Once().Return(nil)

		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})

	t.Run("when invalid command, session error, fail fast scenario enabled", func(*testing.T) {
		session, configuration := &sessionMock{}, newConfiguration(ConfigurationAttr{IsCmdFailFast: true})
		server, errorMessage := newServer(configuration), configuration.msgInvalidCmdHeloArg

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("not implemented command", nil)
		session.On("writeResponse", configuration.msgInvalidCmd, defaultSessionResponseDelay).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("helo 42", nil)
		session.On("clearError").Once().Return(nil)
		session.On("addError", errors.New(errorMessage)).Once().Return(nil)
		session.On("writeResponse", errorMessage, defaultSessionResponseDelay).Once().Return(nil)

		session.On("isErrorFound").Once().Return(true)
		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})

	t.Run("when server quit channel was closed", func(*testing.T) {
		session, configuration := &sessionMock{}, newConfiguration(ConfigurationAttr{IsCmdFailFast: true})
		server := newServer(configuration)
		server.quit = make(chan interface{})
		close(server.quit)

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)
		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})

	t.Run("when read request session error", func(*testing.T) {
		session, configuration := &sessionMock{}, newConfiguration(ConfigurationAttr{IsCmdFailFast: true})
		server := newServer(configuration)

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)
		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return(emptyString, errors.New("some read request error"))
		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})
}
