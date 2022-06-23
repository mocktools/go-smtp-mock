package smtpmock

// RSET command handler
type handlerRset struct {
	*handler
}

// RSET command handler builder. Returns pointer to new handlerRset structure
func newHandlerRset(session sessionInterface, message *message, configuration *configuration) *handlerRset {
	return &handlerRset{&handler{session: session, message: message, configuration: configuration}}
}

// RSET handler methods

// Main RSET handler runner
func (handler *handlerRset) run(request string) {
	if handler.isInvalidRequest(request) {
		return
	}
	configuration := handler.configuration
	handler.message.rset = true
	handler.message.rsetRequest = request
	handler.message.rsetResponse = configuration.msgRsetReceived

	handler.session.writeResponse(configuration.msgRsetReceived, configuration.responseDelayRset)
}

// Invalid QUIT command predicate. Returns true when request is invalid,
// otherwise returns false.
func (handler *handlerRset) isInvalidRequest(request string) bool {
	return !matchRegex(request, validRsetCmdsRegexPattern)
}
