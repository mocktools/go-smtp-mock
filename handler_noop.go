package smtpmock

import "fmt"

// NOOP command handler
type handlerNoop struct {
	*handler
}

// NOOP command handler builder. Returns pointer to new handlerNoop structure
func newHandlerNoop(session sessionInterface, message *Message, configuration *configuration) *handlerNoop {
	return &handlerNoop{&handler{session: session, message: message, configuration: configuration}}
}

// NOOP handler methods

// Main NOOP handler runner
func (handler *handlerNoop) run(request string) {
	fmt.Println("NOOP incoming: ", request)
	if handler.isInvalidRequest(request) {
		fmt.Println("NOOP invalid")
		return
	}
	fmt.Println("Increase noop")
	handler.message.noopCount++ // Not realy thread save
	configuration := handler.configuration
	fmt.Println("START return NOOP response: ", configuration.msgNoopCmd)
	handler.session.writeResponse(configuration.msgNoopCmd, configuration.responseDelayNoop)
	// handler.session.writeResponse(defaultOkMsg, configuration.responseDelayNoop)
	fmt.Println("END return NOOP response")
}

// Invalid NOOP command predicate. Returns true when request is invalid, otherwise returns false
func (handler *handlerNoop) isInvalidRequest(request string) bool {
	return !matchRegex(request, validNoopCmdRegexPattern)
}
