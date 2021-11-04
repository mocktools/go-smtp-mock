package smtpmock

// SMTP mock configuration structure. Provides to configure mock behaviour
type configuration struct {
	hostAddress                   string
	portNumber                    int
	logToStdout                   bool
	logServerActivity             bool
	isCmdFailFast                 bool
	msgGreeting                   string
	msgInvalidCmd                 string
	msgQuitCmd                    string
	msgInvalidCmdHeloSequence     string
	msgInvalidCmdHeloArg          string
	msgHeloBlacklistedDomain      string
	msgHeloReceived               string
	msgInvalidCmdMailfromSequence string
	msgInvalidCmdMailfromArg      string
	msgMailfromBlacklistedEmail   string
	msgMailfromReceived           string
	msgInvalidCmdRcpttoSequence   string
	msgInvalidCmdRcpttoArg        string
	msgRcpttoNotRegisteredEmail   string
	msgRcpttoBlacklistedEmail     string
	msgRcpttoReceived             string
	msgInvalidCmdDataSequence     string
	msgDataReceived               string
	blacklistedHeloDomains        []string
	blacklistedMailfromEmails     []string
	blacklistedRcpttoEmails       []string
	notRegisteredEmails           []string
	// TODO: add ability to send 221 response before end of the session
}

// New configuration builder. Returns pointer to valid new configuration structure
func NewConfiguration(config ConfigurationAttr) *configuration {
	config.assignDefaultValues()

	return &configuration{
		hostAddress:                   config.hostAddress,
		portNumber:                    config.portNumber,
		logToStdout:                   config.logToStdout,
		logServerActivity:             config.logServerActivity,
		isCmdFailFast:                 config.isCmdFailFast,
		msgGreeting:                   config.msgGreeting,
		msgInvalidCmd:                 config.msgInvalidCmd,
		msgInvalidCmdHeloSequence:     config.msgInvalidCmdHeloSequence,
		msgInvalidCmdHeloArg:          config.msgInvalidCmdHeloArg,
		msgHeloBlacklistedDomain:      config.msgHeloBlacklistedDomain,
		msgHeloReceived:               config.msgHeloReceived,
		msgInvalidCmdMailfromSequence: config.msgInvalidCmdMailfromSequence,
		msgInvalidCmdMailfromArg:      config.msgInvalidCmdMailfromArg,
		msgMailfromBlacklistedEmail:   config.msgMailfromBlacklistedEmail,
		msgMailfromReceived:           config.msgMailfromReceived,
		msgInvalidCmdRcpttoSequence:   config.msgInvalidCmdRcpttoSequence,
		msgInvalidCmdRcpttoArg:        config.msgInvalidCmdRcpttoArg,
		msgRcpttoNotRegisteredEmail:   config.msgRcpttoNotRegisteredEmail,
		msgRcpttoBlacklistedEmail:     config.msgRcpttoBlacklistedEmail,
		msgRcpttoReceived:             config.msgRcpttoReceived,
		msgInvalidCmdDataSequence:     config.msgInvalidCmdDataSequence,
		msgDataReceived:               config.msgDataReceived,
		msgQuitCmd:                    config.msgQuitCmd,
		blacklistedHeloDomains:        config.blacklistedHeloDomains,
		blacklistedMailfromEmails:     config.blacklistedMailfromEmails,
		blacklistedRcpttoEmails:       config.blacklistedRcpttoEmails,
		notRegisteredEmails:           config.notRegisteredEmails,
	}
}

// ConfigurationAttr kwargs structure for configuration builder
type ConfigurationAttr struct {
	hostAddress                   string
	portNumber                    int
	logToStdout                   bool
	logServerActivity             bool
	isCmdFailFast                 bool
	msgGreeting                   string
	msgInvalidCmd                 string
	msgQuitCmd                    string
	msgInvalidCmdHeloSequence     string
	msgInvalidCmdHeloArg          string
	msgHeloBlacklistedDomain      string
	msgHeloReceived               string
	msgInvalidCmdMailfromSequence string
	msgInvalidCmdMailfromArg      string
	msgMailfromBlacklistedEmail   string
	msgMailfromReceived           string
	msgInvalidCmdRcpttoSequence   string
	msgInvalidCmdRcpttoArg        string
	msgRcpttoNotRegisteredEmail   string
	msgRcpttoBlacklistedEmail     string
	msgRcpttoReceived             string
	msgInvalidCmdDataSequence     string
	msgDataReceived               string
	blacklistedHeloDomains        []string
	blacklistedMailfromEmails     []string
	blacklistedRcpttoEmails       []string
	notRegisteredEmails           []string
}

// ConfigurationAttr methods

// assigns default values to ConfigurationAttr fields
func (config *ConfigurationAttr) assignDefaultValues() {
	if config.hostAddress == EmptyString {
		config.hostAddress = HostAddress
	}
	if config.portNumber == 0 {
		config.portNumber = PortNuber
	}
	if config.msgGreeting == EmptyString {
		config.msgGreeting = DefaultGreetingMsg
	}
	if config.msgInvalidCmd == EmptyString {
		config.msgInvalidCmd = DefaultInvalidCmdMsg
	}
	if config.msgQuitCmd == EmptyString {
		config.msgQuitCmd = DefaultQuitMsg
	}

	// HELO defaults
	if config.msgInvalidCmdHeloSequence == EmptyString {
		config.msgInvalidCmdHeloSequence = DefaultInvalidCmdHeloSequenceMsg
	}
	if config.msgInvalidCmdHeloArg == EmptyString {
		config.msgInvalidCmdHeloArg = DefaultInvalidCmdHeloArgMsg
	}
	if config.msgHeloBlacklistedDomain == EmptyString {
		config.msgHeloBlacklistedDomain = DefaultQuitMsg
	}
	if config.msgHeloReceived == EmptyString {
		config.msgHeloReceived = DefaultReceivedMsg
	}

	// MAIL FROM defaults
	if config.msgInvalidCmdMailfromSequence == EmptyString {
		config.msgInvalidCmdMailfromSequence = DefaultInvalidCmdMailfromSequenceMsg
	}
	if config.msgInvalidCmdMailfromArg == EmptyString {
		config.msgInvalidCmdMailfromArg = DefaultInvalidCmdMailfromArgMsg
	}
	if config.msgMailfromBlacklistedEmail == EmptyString {
		config.msgMailfromBlacklistedEmail = DefaultQuitMsg
	}
	if config.msgMailfromReceived == EmptyString {
		config.msgMailfromReceived = DefaultReceivedMsg
	}

	// RCPT TO defaults
	if config.msgInvalidCmdRcpttoSequence == EmptyString {
		config.msgInvalidCmdRcpttoSequence = DefaultInvalidCmdRcpttoSequenceMsg
	}
	if config.msgInvalidCmdRcpttoArg == EmptyString {
		config.msgInvalidCmdRcpttoArg = DefaultInvalidCmdRcpttoArgMsg
	}
	if config.msgRcpttoBlacklistedEmail == EmptyString {
		config.msgRcpttoBlacklistedEmail = DefaultQuitMsg
	}
	if config.msgRcpttoNotRegisteredEmail == EmptyString {
		config.msgRcpttoNotRegisteredEmail = DefaultNotRegistredRcpttoEmailMsg
	}
	if config.msgRcpttoReceived == EmptyString {
		config.msgRcpttoReceived = DefaultReceivedMsg
	}

	// DATA defaults
	if config.msgInvalidCmdDataSequence == EmptyString {
		config.msgInvalidCmdDataSequence = DefaultInvalidCmdDataSequenceMsg
	}
	if config.msgDataReceived == EmptyString {
		config.msgDataReceived = DefaultReadyForReceiveMsg
	}
}
