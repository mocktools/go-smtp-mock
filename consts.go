package smtpmock

import "log"

const (
	// SMTP mock default messages

	DefaultGreetingMsg                   = "220 Welcome"
	DefaultQuitMsg                       = "221 Closing connection"
	DefaultReceivedMsg                   = "250 Received"
	DefaultReadyForReceiveMsg            = "354 Ready for receive message. End data with <CR><LF>.<CR><LF>"
	DefaultInvalidCmdHeloArgMsg          = "501 HELO requires domain address"
	DefaultInvalidCmdMailfromArgMsg      = "501 MAIL FROM requires valid email address"
	DefaultInvalidCmdRcpttoArgMsg        = "501 RCPT TO requires valid email address"
	DefaultInvalidCmdMsg                 = "502 Command unrecognized. Available commands: HELO, EHLO, MAIL FROM:, RCPT TO:, DATA"
	DefaultInvalidCmdHeloSequenceMsg     = "503 Bad sequence of commands. HELO should be the first"
	DefaultInvalidCmdMailfromSequenceMsg = "503 Bad sequence of commands. MAIL FROM should be used after HELO"
	DefaultInvalidCmdRcpttoSequenceMsg   = "503 Bad sequence of commands. RCPT TO should be used after MAIL FROM"
	DefaultInvalidCmdDataSequenceMsg     = "503 Bad sequence of commands. DATA should be used after RCPT TO"
	DefaultNotRegistredRcpttoEmailMsg    = "550 User not found"
	DefaultMsgSizeIsTooBigMsg            = "552 Message exceeded max size of"

	// Logger

	InfoLogLevel    = "INFO"
	WarningLogLevel = "WARNING"
	ErrorLogLevel   = "ERROR"
	LogFlag         = log.Ldate | log.Lmicroseconds

	// Session

	SessionStartMsg      = "New SMTP session started"
	SessionRequestMsg    = "SMTP request: "
	SessionResponseMsg   = "SMTP response: "
	SessionEndMsg        = "SMTP session finished"
	SessionBinaryDataMsg = "message binary data portion"

	// Server

	NetworkProtocol               = "tcp"
	DefaultHostAddress            = "0.0.0.0"
	DefaultPortNuber              = 2525
	DefaultMessageSizeLimit       = 10485760 // in bytes (10MB)
	DefaultSessionTimeout         = 30       // in seconds
	ServerStartMsg                = "SMTP mock server started on port"
	ServerErrorMsg                = "Failed to start SMTP mock server on port"
	ServerNotAcceptNewConnections = "SMTP mock server is in the shutdown mode and won't accept new connections"

	// Regex patterns

	AvailableCmdsRegexPattern          = `(?i)helo|ehlo|mail from:|rcpt to:|data|quit`
	DomainRegexPattern                 = `(?i)([\p{L}0-9]+([\-.]{1}[\p{L}0-9]+)*\.\p{L}{2,63})`
	EmailRegexPattern                  = `(?i)<?((.+)@` + DomainRegexPattern + `)>?`
	ValidHeloCmdsRegexPattern          = `(?i)helo|ehlo`
	ValidMailfromCmdRegexPattern       = `(?i)mail from:`
	ValidRcpttoCmdRegexPattern         = `(?i)rcpt to:`
	ValidDataCmdRegexPattern           = `\A(?i)data\z`
	ValidQuitCmdRegexPattern           = `\A(?i)quit\z`
	ValidHeloComplexCmdRegexPattern    = `\A(` + ValidHeloCmdsRegexPattern + `) (` + DomainRegexPattern + `)\z`
	ValidMailromComplexCmdRegexPattern = `\A(` + ValidMailfromCmdRegexPattern + `) ?(` + EmailRegexPattern + `)\z`
	ValidRcpttoComplexCmdRegexPattern  = `\A(` + ValidRcpttoCmdRegexPattern + `) ?(` + EmailRegexPattern + `)\z`

	// Helpers

	EmptyString = ""
)
