package smtpmock

import "errors"

// Base handler
type handler struct {
	session       sessionInterface
	message       *message
	configuration *configuration
}

// handler methods

// Writes hadled result to session, message. Always returns true
func (handler *handler) writeResult(isSuccessful bool, request, response string) bool {
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
func (handler *handler) isInvalidCmd(request string) bool {
	if !matchRegex(request, AvailableCmdsRegexPattern) {
		return handler.writeResult(false, request, handler.configuration.msgInvalidCmd)
	}

	return false
}
