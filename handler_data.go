package smtpmock

import (
	"errors"
)

// DATA command handler
type handlerData struct {
	*handler
}

// DATA command handler builder. Returns pointer to new handlerData structure
func newHandlerData(session sessionInterface, message *message, configuration *configuration) *handlerData {
	return &handlerData{&handler{session: session, message: message, configuration: configuration}}
}

// DATA handler methods

// Main DATA handler runner
func (handler *handlerData) run() {
	var requestSnapshot string
	session := handler.session

	if handler.isFailFastScenario() {
		request, err := session.readRequest()
		if err != nil {
			return
		}

		if handler.isInvalidRequest(request) {
			return
		}
		requestSnapshot = request
	}

	if !handler.isFailFastScenario() {
		for {
			session.clearError()
			request, err := session.readRequest()
			if err != nil {
				return
			}

			if !handler.isInvalidRequest(request) {
				requestSnapshot = request
				break
			}
		}
	}

	handler.writeResult(true, requestSnapshot, handler.configuration.msgDataReceived)
}

// Writes handled DATA result to session, message. Always returns true
func (handler *handlerData) writeResult(isSuccessful bool, request, response string) bool {
	session, message := handler.session, handler.message
	if !isSuccessful {
		session.addError(errors.New(response))
	}

	message.dataRequest, message.dataResponse, message.data = request, response, isSuccessful
	session.writeResponse(response)
	return true
}

// Invalid SMTP command predicate. Returns true and writes result for case when command is invalid,
// otherwise returns false.
func (handler *handlerData) isInvalidCmd(request string) bool {
	if !matchRegex(request, AvailableCmdsRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmd)
	}

	return false
}

// Invalid DATA command sequence predicate. Returns true and writes result for case when DATA
// command sequence is invalid, otherwise returns false
func (handler *handlerData) isInvalidCmdSequence(request string) bool {
	if !matchRegex(request, ValidDataCmdRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmdDataSequence)
	}

	return false
}

// Invalid DATA command request complex predicate. Returns true for case when one
// of the chain checks returns true, otherwise returns false
func (handler *handlerData) isInvalidRequest(request string) bool {
	return handler.isInvalidCmd(request) || handler.isInvalidCmdSequence(request)
}
