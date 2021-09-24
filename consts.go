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

	// Session

	SessionStartMsg    = "New SMTP session started"
	SessionRequestMsg  = "SMTP request: "
	SessionResponseMsg = "SMTP response: "
	SessionEndMsg      = "SMTP session finished"

	// Server

	NetworkProtocol = "tcp"
	HostAddress     = "0.0.0.0"
	PortNuber       = 2525
	ServerStartMsg  = "SMTP mock server started on port"
	ServerErrorMsg  = "Failed to start SMTP mock server on port"

	// Regex patterns

	AvailableCmdsRegexPattern = `(?i)helo|ehlo|mail from:|rcpt to:`
	ValidHeloCmdsRegexPattern = `(?i)helo|ehlo`
	DomainRegexPattern        = `(?i)([\p{L}0-9]+([\-.]{1}[\p{L}0-9]+)*\.\p{L}{2,63})`
	ValidHeloCmdRegexPattern  = `^(?i)(helo|ehlo) ` + DomainRegexPattern + `$`

	// Helpers

	EmptyString = ""
)
