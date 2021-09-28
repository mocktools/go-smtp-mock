package smtpmock

import "errors"

// MAILFROM command handler
type handlerMailfrom struct {
	*handler
}

// MAILFROM command handler builder. Returns pointer to new handlerHelo structure
func newHandlerMailfrom(session sessionInterface, message *message, configuration *configuration) *handlerMailfrom {
	return &handlerMailfrom{&handler{session: session, message: message, configuration: configuration}}
}

// MAILFROM handler methods

// Main MAILFROM handler runner
func (handler *handlerMailfrom) run() {
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

	if handler.isBlacklistedEmail(requestSnapshot) {
		return
	}

	handler.writeResult(true, requestSnapshot, handler.configuration.msgMailfromReceived)
}

// Writes hadled HELO result to session, message. Always returns true
func (handler *handlerMailfrom) writeResult(isSuccessful bool, request, response string) bool {
	session, message := handler.session, handler.message
	if !isSuccessful {
		session.addError(errors.New(response))
	}

	message.mailfromRequest, message.mailfromResponse, message.mailfrom = request, response, isSuccessful
	session.writeResponse(response)
	return true
}

// Invalid SMTP command predicate. Returns true and writes result for case when command is invalid,
// otherwise returns false.
func (handler *handlerMailfrom) isInvalidCmd(request string) bool {
	if !matchRegex(request, AvailableCmdsRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmd)
	}

	return false
}

// Invalid MAILFROM command sequence predicate. Returns true and writes result for case when MAILFROM
// command sequence is invalid, otherwise returns false
func (handler *handlerMailfrom) isInvalidCmdSequence(request string) bool {
	if !matchRegex(request, ValidMailfromCmdRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmdMailfromSequence)
	}

	return false
}

// Invalid MAILFROM command argument predicate. Returns true and writes result for case when MAILFROM
// command argument is invalid, otherwise returns false
func (handler *handlerMailfrom) isInvalidCmdArg(request string) bool {
	if !matchRegex(request, ValidMaifromCmdRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmdMailfromArg)
	}

	return false
}

// Invalid MAILFROM command request complex predicate. Returns true for case when one
// of the chain checks returns true, otherwise returns false
func (handler *handlerMailfrom) isInvalidRequest(request string) bool {
	return handler.isInvalidCmd(request) || handler.isInvalidCmdSequence(request) || handler.isInvalidCmdArg(request)
}

// Returns domain from HELO request
func (handler *handlerMailfrom) mailfromEmail(request string) string {
	return regexCaptureGroup(request, ValidMaifromCmdRegexPattern, 3)
}

// Custom behaviour for HELO domain predicate. Returns true and writes result for case when HELO domain
// is included in configuration.blacklistedHeloDomains slice
func (handler *handlerMailfrom) isBlacklistedEmail(request string) bool {
	configuration := handler.configuration
	if !isIncluded(configuration.blacklistedMailfromEmails, handler.mailfromEmail(request)) {
		return false
	}

	return handler.writeResult(false, request, configuration.msgMailfromBlacklistedEmail)
}
