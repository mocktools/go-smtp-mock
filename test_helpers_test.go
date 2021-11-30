package smtpmock

import (
	"net"
	"net/smtp"
	"regexp"
	"time"
)

// Returns log message regex based on log level and message context
func loggerMessageRegex(logLevel, logMessage string) *regexp.Regexp {
	regex, _ := newRegex(logLevel + `: \d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}\.\d{6} ` + logMessage)
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

// Runs minimal successfull SMTP session with target host
func runMinimalSuccessfulSmtpSession(hostAddress string, portNumber int) error {
	connection, _ := net.DialTimeout(NetworkProtocol, serverWithPortNumber(hostAddress, portNumber), time.Duration(2)*time.Second)
	client, _ := smtp.NewClient(connection, hostAddress)

	if err := client.Hello("example.com"); err != nil {
		return err
	}
	if err := client.Quit(); err != nil {
		return err
	}
	if err := client.Close(); err != nil {
		return err
	}

	return nil
}
