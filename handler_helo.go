package smtpmock

import "errors"

// HELO command handler
type handlerHelo struct {
	session       sessionInterface
	message       *message
	configuration *configuration
}

// HELO command handler builder. Returns pointer to new handlerHelo structure
func newHandlerHelo(session sessionInterface, message *message, configuration *configuration) *handlerHelo {
	return &handlerHelo{session: session, message: message, configuration: configuration}
}

// HELO handler methods

// Main HELO handler runner
func (handler *handlerHelo) run() {
	var requestSnapshot string
	configuration, session := handler.configuration, handler.session
	failFastScenario := configuration.isCmdFailFast

	if failFastScenario {
		request, err := session.readRequest()
		if err != nil {
			return
		}

		if handler.isInvalidRequest(request) {
			return
		}
		requestSnapshot = request
	}

	if !failFastScenario {
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

	if handler.isBlacklistedDomain(requestSnapshot) {
		return
	}

	handler.writeResult(true, requestSnapshot, configuration.msgHeloReceived)
}

// Writes hadled HELO result to session, message. Always returns true
func (handler *handlerHelo) writeResult(isSuccessful bool, request, response string) bool {
	session, message := handler.session, handler.message
	if !isSuccessful {
		session.addError(errors.New(response))
	}

	message.heloRequest, message.heloResponse, message.helo = request, response, isSuccessful
	session.writeResponse(response)
	return true
}

// Invalid SMTP command predicate. Returns true and writes result for case when command is invalid,
// otherwise returns false.
func (handler *handlerHelo) isInvalidCmd(request string) bool {
	if !matchRegex(request, AvailableCmdsRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmd)
	}

	return false
}

// Invalid HELO command sequence predicate. Returns true and writes result for case when HELO command
// sequence is invalid, otherwise returns false
func (handler *handlerHelo) isInvalidCmdSequence(request string) bool {
	if !matchRegex(request, ValidHeloCmdsRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmdHeloSequence)
	}

	return false
}

// Invalid HELO command argument predicate. Returns true and writes result for case when HELO command
// argument is invalid, otherwise returns false
func (handler *handlerHelo) isInvalidCmdArg(request string) bool {
	if !matchRegex(request, ValidHeloCmdRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmdHeloArg)
	}

	return false
}

// Invalid HELO command request complex predicate. Returns true for case when one
// of the chain checks returns true, otherwise returns false
func (handler *handlerHelo) isInvalidRequest(request string) bool {
	return handler.isInvalidCmd(request) || handler.isInvalidCmdSequence(request) || handler.isInvalidCmdArg(request)
}

// Returns domain from HELO request
func (handler *handlerHelo) heloDomain(request string) string {
	return regexCaptureGroup(request, DomainRegexPattern, 0)
}

// Custom behaviour for HELO domain predicate. Returns true and writes result for case when HELO domain
// is included in configuration.blacklistedHeloDomains slice
func (handler *handlerHelo) isBlacklistedDomain(request string) bool {
	configuration := handler.configuration
	if !isIncluded(configuration.blacklistedHeloDomains, handler.heloDomain(request)) {
		return false
	}

	return handler.writeResult(false, request, configuration.msgQuit)
}
