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

// Server structure which implements SMTP mock server
type Server struct {
	configuration *configuration
	messages      *messages
	logger        logger
	listener      net.Listener
	wg            waitGroup
	quit          chan interface{}
	isStarted     bool
	PortNumber    int
}

// SMTP mock server builder, creates new server
func newServer(configuration *configuration) *Server {
	return &Server{
		configuration: configuration,
		messages:      new(messages),
		logger:        newLogger(configuration.logToStdout, configuration.logServerActivity),
		wg:            new(sync.WaitGroup),
	}
}

// server methods

// Start binds and runs SMTP mock server on specified port or random free port. Returns error for
// case when server is active. Server port number will be assigned after successful start only
func (server *Server) Start() (err error) {
	if server.isStarted {
		return errors.New(serverStartErrorMsg)
	}

	configuration, logger := server.configuration, server.logger
	portNumber := configuration.portNumber

	listener, err := net.Listen(networkProtocol, serverWithPortNumber(configuration.hostAddress, portNumber))
	if err != nil {
		errorMessage := fmt.Sprintf("%s: %d", serverErrorMsg, portNumber)
		logger.error(errorMessage)
		return errors.New(errorMessage)
	}

	portNumber = listener.Addr().(*net.TCPAddr).Port
	server.listener, server.isStarted, server.quit, server.PortNumber = listener, true, make(chan interface{}), portNumber
	logger.infoActivity(fmt.Sprintf("%s: %d", serverStartMsg, portNumber))

	server.addToWaitGroup()
	go func() {
		defer server.removeFromWaitGroup()
		for {
			connection, err := server.listener.Accept()
			if err != nil {
				if _, ok := <-server.quit; !ok {
					logger.warning(serverNotAcceptNewConnectionsMsg)
				}
				return
			}

			server.addToWaitGroup()
			go func() {
				server.handleSession(newSession(connection, logger))
				server.removeFromWaitGroup()
			}()

			logger.infoActivity(sessionStartMsg)
		}
	}()

	return err
}

// Stop shutdowns server gracefully. Returns error for case when server is not active
func (server *Server) Stop() (err error) {
	if server.isStarted {
		close(server.quit)
		server.listener.Close()
		server.wg.Wait()
		server.isStarted = false
		server.logger.infoActivity(serverStopMsg)
		return
	}

	return errors.New(serverStopErrorMsg)
}

// Creates and assigns new message to server.messages
func (server *Server) newMessage() *message {
	newMessage := new(message)
	server.messages.append(newMessage)
	return newMessage
}

// Invalid SMTP command predicate. Returns true when command is invalid, otherwise returns false
func (server *Server) isInvalidCmd(request string) bool {
	return !matchRegex(request, availableCmdsRegexPattern)
}

// Recognizes command implemented commands. Captures the first word divided by spaces,
// converts it to upper case
func (server *Server) recognizeCommand(request string) string {
	command := strings.Split(request, " ")[0]
	return strings.ToUpper(command)
}

// Addes goroutine to WaitGroup
func (server *Server) addToWaitGroup() {
	server.wg.Add(1)
}

// Removes goroutine from WaitGroup
func (server *Server) removeFromWaitGroup() {
	server.wg.Done()
}

// SMTP client-server session handler
func (server *Server) handleSession(session sessionInterface) {
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
