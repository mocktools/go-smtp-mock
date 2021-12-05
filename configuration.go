package smtpmock

import "fmt"

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
	msgMsgSizeIsTooBig            string
	msgMsgReceived                string
	blacklistedHeloDomains        []string
	blacklistedMailfromEmails     []string
	blacklistedRcpttoEmails       []string
	notRegisteredEmails           []string
	msqSizeLimit                  int
	sessionTimeout                int

	// TODO: add ability to send 221 response before end of session for case when fail fast scenario enabled
}

// New configuration builder. Returns pointer to valid new configuration structure
func newConfiguration(config ConfigurationAttr) *configuration {
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
		msgMsgSizeIsTooBig:            config.msgMsgSizeIsTooBig,
		msgMsgReceived:                config.msgMsgReceived,
		msgQuitCmd:                    config.msgQuitCmd,
		blacklistedHeloDomains:        config.blacklistedHeloDomains,
		blacklistedMailfromEmails:     config.blacklistedMailfromEmails,
		blacklistedRcpttoEmails:       config.blacklistedRcpttoEmails,
		notRegisteredEmails:           config.notRegisteredEmails,
		msqSizeLimit:                  config.msqSizeLimit,
		sessionTimeout:                config.sessionTimeout,
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
	msgMsgSizeIsTooBig            string
	msgMsgReceived                string
	blacklistedHeloDomains        []string
	blacklistedMailfromEmails     []string
	blacklistedRcpttoEmails       []string
	notRegisteredEmails           []string
	msqSizeLimit                  int
	sessionTimeout                int
}

// ConfigurationAttr methods

// Assigns server defaults
func (config *ConfigurationAttr) assignServerDefaultValues() {
	if config.hostAddress == emptyString {
		config.hostAddress = defaultHostAddress
	}
	if config.portNumber == 0 {
		config.portNumber = defaultPortNuber
	}
	if config.msgGreeting == emptyString {
		config.msgGreeting = defaultGreetingMsg
	}
	if config.msgInvalidCmd == emptyString {
		config.msgInvalidCmd = defaultInvalidCmdMsg
	}
	if config.msgQuitCmd == emptyString {
		config.msgQuitCmd = defaultQuitMsg
	}
	if config.sessionTimeout == 0 {
		config.sessionTimeout = defaultSessionTimeout
	}
}

// Assigns handlerHelo defaults
func (config *ConfigurationAttr) assignHandlerHeloDefaultValues() {
	if config.msgInvalidCmdHeloSequence == emptyString {
		config.msgInvalidCmdHeloSequence = defaultInvalidCmdHeloSequenceMsg
	}
	if config.msgInvalidCmdHeloArg == emptyString {
		config.msgInvalidCmdHeloArg = defaultInvalidCmdHeloArgMsg
	}
	if config.msgHeloBlacklistedDomain == emptyString {
		config.msgHeloBlacklistedDomain = defaultQuitMsg
	}
	if config.msgHeloReceived == emptyString {
		config.msgHeloReceived = defaultReceivedMsg
	}
}

// Assigns handlerMailfrom defaults
func (config *ConfigurationAttr) assignHandlerMailfromDefaultValues() {
	if config.msgInvalidCmdMailfromSequence == emptyString {
		config.msgInvalidCmdMailfromSequence = defaultInvalidCmdMailfromSequenceMsg
	}
	if config.msgInvalidCmdMailfromArg == emptyString {
		config.msgInvalidCmdMailfromArg = defaultInvalidCmdMailfromArgMsg
	}
	if config.msgMailfromBlacklistedEmail == emptyString {
		config.msgMailfromBlacklistedEmail = defaultQuitMsg
	}
	if config.msgMailfromReceived == emptyString {
		config.msgMailfromReceived = defaultReceivedMsg
	}
}

// Assigns handlerRcptto defaults
func (config *ConfigurationAttr) assignHandlerRcpttoDefaultValues() {
	if config.msgInvalidCmdRcpttoSequence == emptyString {
		config.msgInvalidCmdRcpttoSequence = defaultInvalidCmdRcpttoSequenceMsg
	}
	if config.msgInvalidCmdRcpttoArg == emptyString {
		config.msgInvalidCmdRcpttoArg = defaultInvalidCmdRcpttoArgMsg
	}
	if config.msgRcpttoBlacklistedEmail == emptyString {
		config.msgRcpttoBlacklistedEmail = defaultQuitMsg
	}
	if config.msgRcpttoNotRegisteredEmail == emptyString {
		config.msgRcpttoNotRegisteredEmail = defaultNotRegistredRcpttoEmailMsg
	}
	if config.msgRcpttoReceived == emptyString {
		config.msgRcpttoReceived = defaultReceivedMsg
	}
}

// Assigns handlerData defaults
func (config *ConfigurationAttr) assignHandlerDataDefaultValues() {
	if config.msgInvalidCmdDataSequence == emptyString {
		config.msgInvalidCmdDataSequence = defaultInvalidCmdDataSequenceMsg
	}
	if config.msgDataReceived == emptyString {
		config.msgDataReceived = defaultReadyForReceiveMsg
	}
}

// Assigns handlerMessage defaults
func (config *ConfigurationAttr) assignHandlerMessageDefaultValues() {
	if config.msqSizeLimit == 0 {
		config.msqSizeLimit = defaultMessageSizeLimit
	}
	if config.msgMsgSizeIsTooBig == emptyString {
		config.msgMsgSizeIsTooBig = fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", config.msqSizeLimit)
	}
	if config.msgMsgReceived == emptyString {
		config.msgMsgReceived = defaultReceivedMsg
	}
}

// Assigns default values to ConfigurationAttr fields
func (config *ConfigurationAttr) assignDefaultValues() {
	config.assignServerDefaultValues()
	config.assignHandlerHeloDefaultValues()
	config.assignHandlerMailfromDefaultValues()
	config.assignHandlerRcpttoDefaultValues()
	config.assignHandlerDataDefaultValues()
	config.assignHandlerMessageDefaultValues()
}
