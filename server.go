package smtpmock

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

// WaitGroup interface
type waitGroup interface {
	Add(int)
	Done()
	Wait()
}

// SMTP mock server that will be listening for SMTP connections
type server struct {
	configuration *configuration
	messages      *messages
	logger        logger
	listener      net.Listener
	wg            waitGroup
	quit          chan interface{}
	isStarted     bool
}

// SMTP mock server builder, creates new server
func newServer(configuration *configuration) *server {
	return &server{
		configuration: configuration,
		messages:      new(messages),
		logger:        newLogger(configuration.logToStdout, configuration.logServerActivity),
		wg:            new(sync.WaitGroup),
		quit:          make(chan interface{}),
	}
}

// server methods

// Creates and assigns new message to server.messages
func (server *server) newMessage() *message {
	newMessage := new(message)
	server.messages.append(newMessage)
	return newMessage
}

// Binds and runs SMTP mock server on specified port
func (server *server) Start() (err error) {
	if server.isStarted {
		return errors.New(ServerStartErrorMsg)
	}

	configuration, logger := server.configuration, server.logger
	portNumber := configuration.portNumber

	listener, err := net.Listen(NetworkProtocol, serverWithPortNumber(configuration.hostAddress, portNumber))
	if err != nil {
		errorMessage := fmt.Sprintf("%s: %d", ServerErrorMsg, server.configuration.portNumber)
		logger.error(errorMessage)
		return errors.New(errorMessage)
	}

	server.listener, server.isStarted = listener, true
	logger.infoActivity(fmt.Sprintf("%s: %d", ServerStartMsg, portNumber))

	server.addToWaitGroup()
	go func() {
		defer server.removeFromWaitGroup()
		for {
			connection, err := server.listener.Accept()
			if err != nil {
				if _, ok := <-server.quit; !ok {
					logger.warning(ServerNotAcceptNewConnectionsMsg)
				}
				return
			}

			server.addToWaitGroup()
			go func() {
				server.handleSession(newSession(connection, logger))
				server.removeFromWaitGroup()
			}()

			logger.infoActivity(SessionStartMsg)
		}
	}()

	return err
}

// Stops server gracefully. Returns error for case when server is not active
func (server *server) Stop() (err error) {
	if server.isStarted {
		close(server.quit)
		server.listener.Close()
		server.wg.Wait()
		server.isStarted = false
		server.logger.infoActivity(ServerStopMsg)
		return
	}

	return errors.New(ServerStopErrorMsg)
}

// Invalid SMTP command predicate. Returns true when command is invalid, otherwise returns false
func (server *server) isInvalidCmd(request string) bool {
	return !matchRegex(request, AvailableCmdsRegexPattern)
}

// Recognizes command implemented commands. Captures the first word divided by spaces,
// converts it to upper case
func (server *server) recognizeCommand(request string) string {
	command := strings.Split(request, " ")[0]
	return strings.ToUpper(command)
}

// Addes goroutine to WaitGroup
func (server *server) addToWaitGroup() {
	server.wg.Add(1)
}

// Removes goroutine from WaitGroup
func (server *server) removeFromWaitGroup() {
	server.wg.Done()
}

// SMTP client-server session handler
func (server *server) handleSession(session sessionInterface) {
	defer session.finish()
	message, configuration := server.newMessage(), server.configuration
	session.writeResponse(configuration.msgGreeting)

	for {
		select {
		case <-server.quit:
			return
		default:
			session.setTimeout(configuration.sessionTimeout)
			request, err := session.readRequest()
			if err != nil {
				return
			}

			if server.isInvalidCmd(request) {
				session.writeResponse(configuration.msgInvalidCmd)
				continue
			}

			switch server.recognizeCommand(request) {
			case "HELO", "EHLO":
				newHandlerHelo(session, message, configuration).run(request)
			case "MAIL":
				newHandlerMailfrom(session, message, configuration).run(request)
			case "RCPT":
				newHandlerRcptto(session, message, configuration).run(request)
			case "DATA":
				newHandlerData(session, message, configuration).run(request)
			case "QUIT":
				newHandlerQuit(session, message, configuration).run(request)
			}

			if message.quitSent || (session.isErrorFound() && configuration.isCmdFailFast) {
				return
			}
		}
	}
}
