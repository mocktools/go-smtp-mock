package smtpmock

// QUIT command handler
type handlerQuit struct {
	*handler
}

// QUIT command handler builder. Returns pointer to new handlerQuit structure
func newHandlerQuit(session sessionInterface, message *Message, configuration *configuration) *handlerQuit {
	return &handlerQuit{&handler{session: session, message: message, configuration: configuration}}
}

// QUIT handler methods

// Main QUIT handler runner. Appends the message to serverMessages before
// writing the QUIT response so that Messages() returns consistent results
// as soon as the client receives "221".
func (handler *handlerQuit) run(request string, serverMessages *messages) {
	if handler.isInvalidRequest(request) {
		return
	}

	handler.message.quitSent = true
	serverMessages.append(handler.message)
	configuration := handler.configuration
	handler.session.writeResponse(configuration.msgQuitCmd, configuration.responseDelayQuit)
}

// Invalid QUIT command predicate. Returns true when request is invalid, otherwise returns false
func (handler *handlerQuit) isInvalidRequest(request string) bool {
	return !validQuitCmdRegex.MatchString(request)
}
