package smtpmock

import "log"

const (
	// SMTP mock default messages

	DefaultGreetingMsg               = "220 Welcome"
	DefaultReceivedMsg               = "250 Received"
	DefaultInvalidCmdHeloArgMsg      = "501 HELO requires domain address"
	DefaultInvalidCmdMsg             = "502 Command unrecognized. Available commands: HELO, EHLO, MAIL FROM:, RCPT TO:"
	DefaultInvalidCmdHeloSequenceMsg = "503 Bad sequence of commands. HELO should be the first"
	DefaultQuitMsg                   = "221 Closing connection"

	// Logger

	InfoLogLevel    = "INFO"
	WarningLogLevel = "WARNING"
	ErrorLogLevel   = "ERROR"
	LogFlag         = log.Ldate | log.Ltime

	// Helpers

	EmptyString = ""
)
