package smtpmock

import (
	"regexp"
)

// Returns log message regex based on log level and message context
func loggerMessageRegex(logLevel, logMessage string) *regexp.Regexp {
	regex, _ := newRegex(logLevel + `: \d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2} ` + logMessage)
	return regex
}

// Creates configuration with default settings
func createConfiguration() *configuration {
	return NewConfiguration(ConfigurationAttr{})
}

// Creates not empty message
func createNotEmptyMessage() *message {
	return &message{
		heloRequest:      "a",
		heloResponse:     "b",
		mailfromRequest:  "c",
		mailfromResponse: "d",
		rcpttoRequest:    "a",
		rcpttoResponse:   "b",
		dataRequest:      "c",
		dataResponse:     "d",
		msgRequest:       "a",
		msgResponse:      "b",
		helo:             true,
		mailfrom:         true,
		rcptto:           true,
		data:             true,
		msg:              true,
	}
}
