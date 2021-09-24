package smtpmock

// SMTP mock configuration structure. Provides to configure mock behaviour
type configuration struct {
	hostAddress               string
	portNumber                int
	logToStdout               bool
	logServerActivity         bool
	isCmdFailFast             bool
	msgGreeting               string
	msgHeloReceived           string
	msgInvalidCmdHeloArg      string
	msgInvalidCmd             string
	msgInvalidCmdHeloSequence string
	msgQuit                   string
	blacklistedHeloDomains    []string
}

// New configuration builder. Returns pointer to valid new configuration structure
func NewConfiguration(config ConfigurationAttr) *configuration {
	config.assignDefaultValues()

	return &configuration{
		hostAddress:               config.hostAddress,
		portNumber:                config.portNumber,
		logToStdout:               config.logToStdout,
		logServerActivity:         config.logServerActivity,
		isCmdFailFast:             config.isCmdFailFast,
		msgGreeting:               config.msgGreeting,
		msgHeloReceived:           config.msgHeloReceived,
		msgInvalidCmdHeloArg:      config.msgInvalidCmdHeloArg,
		msgInvalidCmd:             config.msgInvalidCmd,
		msgInvalidCmdHeloSequence: config.msgInvalidCmdHeloSequence,
		msgQuit:                   config.msgQuit,
		blacklistedHeloDomains:    config.blacklistedHeloDomains,
	}
}

// ConfigurationAttr kwargs structure for configuration builder
type ConfigurationAttr struct {
	hostAddress               string
	portNumber                int
	logToStdout               bool
	logServerActivity         bool
	isCmdFailFast             bool
	msgGreeting               string
	msgHeloReceived           string
	msgInvalidCmdHeloArg      string
	msgInvalidCmd             string
	msgInvalidCmdHeloSequence string
	msgQuit                   string
	blacklistedHeloDomains    []string
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
	if config.msgHeloReceived == EmptyString {
		config.msgHeloReceived = DefaultReceivedMsg
	}
	if config.msgInvalidCmdHeloArg == EmptyString {
		config.msgInvalidCmdHeloArg = DefaultInvalidCmdHeloArgMsg
	}
	if config.msgInvalidCmd == EmptyString {
		config.msgInvalidCmd = DefaultInvalidCmdMsg
	}
	if config.msgInvalidCmdHeloSequence == EmptyString {
		config.msgInvalidCmdHeloSequence = DefaultInvalidCmdHeloSequenceMsg
	}
	if config.msgQuit == EmptyString {
		config.msgQuit = DefaultQuitMsg
	}
}
