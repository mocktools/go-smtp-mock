package smtpmock

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	mail "github.com/xhit/go-simple-mail/v2"
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
		assert.Equal(t, 0, server.PortNumber)
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

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)

		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return("helo example.com", nil)
		session.On("clearError").Once().Return(nil)
		session.On("writeResponse", configuration.msgHeloReceived, configuration.responseDelayHelo).Once().Return(nil)
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
	})

	t.Run("when invalid command, fail fast scenario disabled", func(t *testing.T) {
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

	t.Run("when invalid command, session error, fail fast scenario enabled", func(t *testing.T) {
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

	t.Run("when server quit channel was closed", func(t *testing.T) {
		session, configuration := &sessionMock{}, newConfiguration(ConfigurationAttr{IsCmdFailFast: true})
		server := newServer(configuration)
		server.quit = make(chan interface{})
		close(server.quit)

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)
		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})

	t.Run("when read request session error", func(t *testing.T) {
		session, configuration := &sessionMock{}, newConfiguration(ConfigurationAttr{IsCmdFailFast: true})
		server := newServer(configuration)

		session.On("writeResponse", configuration.msgGreeting, defaultSessionResponseDelay).Once().Return(nil)
		session.On("setTimeout", defaultSessionTimeout).Once().Return(nil)
		session.On("readRequest").Once().Return(emptyString, errors.New("some read request error"))
		session.On("finish").Once().Return(nil)

		server.handleSession(session)
	})
}

func TestServerStart(t *testing.T) {
	t.Run("when no errors happens during starting and running the server with default port", func(t *testing.T) {
		configuration := createConfiguration()
		server := newServer(configuration)

		assert.NoError(t, server.Start())
		_ = runMinimalSuccessfulSMTPSession(configuration.hostAddress, server.PortNumber)
		assert.NotEmpty(t, server.messages)
		assert.NotNil(t, server.quit)
		assert.NotNil(t, server.quitTimeout)
		assert.True(t, server.isStarted)
		assert.Greater(t, server.PortNumber, 0)

		_ = server.Stop()
	})

	t.Run("when no errors happens during starting and running the server with custom port", func(t *testing.T) {
		configuration, portNumber := createConfiguration(), 2525
		configuration.portNumber = portNumber
		server := newServer(configuration)

		assert.NoError(t, server.Start())
		_ = runMinimalSuccessfulSMTPSession(configuration.hostAddress, portNumber)
		assert.NotEmpty(t, server.messages)
		assert.NotNil(t, server.quit)
		assert.NotNil(t, server.quitTimeout)
		assert.True(t, server.isStarted)
		assert.Equal(t, portNumber, server.PortNumber)

		_ = server.Stop()
	})

	t.Run("when active server doesn't start current server", func(t *testing.T) {
		server := &Server{isStarted: true}

		assert.EqualError(t, server.Start(), serverStartErrorMsg)
		assert.Equal(t, 0, server.PortNumber)
	})

	t.Run("when listener error happens during starting the server doesn't start current server", func(t *testing.T) {
		configuration := createConfiguration()
		server, logger := newServer(configuration), new(loggerMock)
		listener, _ := net.Listen(networkProtocol, emptyString)
		portNumber := listener.Addr().(*net.TCPAddr).Port
		errorMessage := fmt.Sprintf("%s: %d", serverErrorMsg, portNumber)
		configuration.portNumber, server.logger = portNumber, logger
		logger.On("error", errorMessage).Once().Return(nil)

		assert.EqualError(t, server.Start(), errorMessage)
		assert.False(t, server.isStarted)
		assert.Equal(t, 0, server.PortNumber)
		listener.Close()
	})
}

func TestServerStop(t *testing.T) {
	t.Run("when server active stops current server, graceful shutdown case", func(t *testing.T) {
		logger, listener, waitGroup, quitChannel := new(loggerMock), new(listenerMock), new(waitGroupMock), make(chan interface{})
		server := &Server{
			configuration: createConfiguration(),
			logger:        logger,
			listener:      listener,
			wg:            waitGroup,
			quit:          quitChannel,
			isStarted:     true,
			quitTimeout:   make(chan interface{}),
		}
		listener.On("Close").Once().Return(nil)
		waitGroup.On("Wait").Once().Return(nil)
		logger.On("infoActivity", serverStopMsg).Once().Return(nil)

		assert.NoError(t, server.Stop())
		assert.False(t, server.isStarted)
		_, isChannelOpened := <-server.quit
		assert.False(t, isChannelOpened)
	})

	t.Run("when server active stops current server, force shutdown case", func(t *testing.T) {
		logger, listener, waitGroup, quitChannel := new(loggerMock), new(listenerMock), new(waitGroupMock), make(chan interface{})
		server := &Server{
			configuration: createConfiguration(),
			logger:        logger,
			listener:      listener,
			wg:            waitGroup,
			quit:          quitChannel,
			isStarted:     true,
		}
		listener.On("Close").Once().Return(nil)
		waitGroup.On("Wait").Once().Return(nil)
		logger.On("infoActivity", serverForceStopMsg).Once().Return(nil)

		assert.NoError(t, server.Stop())
		assert.False(t, server.isStarted)
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
		server.newMessage()

		assert.NotEmpty(t, server.Messages())
	})
}

func TestServerHandleRSET(t *testing.T) {
	server := New(ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
	})
	err := server.Start()
	if err != nil {
		t.Fatalf("error while starting the server: %s", err.Error())
	}
	srv := mail.NewSMTPClient()
	srv.Host = "127.0.0.1"
	srv.Port = server.PortNumber

	// don't close the connection
	srv.KeepAlive = true
	client, err := srv.Connect()
	if err != nil {
		t.Fatal("could not connect to server")
	}

	t.Run("when multiple mails are send via one connection", func(t *testing.T) {
		email := mail.NewMSG()
		email.SetFrom("sender@test.com").
			AddTo("receiver@test.com").
			SetSubject("subject").
			SetBody(mail.TextHTML, "HTML-body").
			AddAlternative(mail.TextPlain, "TXT-alternative")
		if err := email.Send(client); err != nil {
			t.Fatalf("Error while sending email: %s", err.Error())
		}

		email.SetSubject("subject2")
		email.SetFrom("sender2@test.com")
		if err := email.Send(client); err != nil {
			t.Fatalf("Error while sending email: %s", err.Error())
		}

		email.SetSubject("subject3")
		email.SetFrom("sender3@test.com")
		if err := email.Send(client); err != nil {
			t.Fatalf("Error while sending email: %s", err.Error())
		}

		messages := server.Messages()
		// there should be 3 messages
		assert.Len(t, messages, 3)
		assert.Equal(t, "MAIL FROM:<sender@test.com>", messages[0].mailfromRequest)
		assert.Equal(t, "MAIL FROM:<sender2@test.com>", messages[1].mailfromRequest)
		assert.Equal(t, "MAIL FROM:<sender3@test.com>", messages[2].mailfromRequest)

		assert.Equal(t, "RSET", messages[0].RsetRequest())
		assert.Equal(t, "250 Ok", messages[0].RsetResponse())
		assert.Equal(t, "RSET", messages[1].RsetRequest())
		assert.Equal(t, "250 Ok", messages[1].RsetResponse())
		assert.Equal(t, "RSET", messages[2].RsetRequest())
		assert.Equal(t, "250 Ok", messages[2].RsetResponse())
	})
}
