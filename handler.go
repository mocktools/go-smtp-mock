package smtpmock

// Base handler
type handler struct {
	session       sessionInterface
	message       *message
	configuration *configuration
}

// handler methods

// Fail fast scenario predicate. Returns true if fail fast enabled in configuration,
// otherwise returns false
func (handler *handler) isFailFastScenario() bool {
	return handler.configuration.isCmdFailFast
}
