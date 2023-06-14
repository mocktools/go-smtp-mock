package smtpmock

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	t.Run("creates new configuration with default settings", func(t *testing.T) {
		buildedConfiguration := newConfiguration(ConfigurationAttr{})

		assert.Equal(t, defaultHostAddress, buildedConfiguration.hostAddress)
		assert.False(t, buildedConfiguration.logToStdout)
		assert.False(t, buildedConfiguration.isCmdFailFast)
		assert.False(t, buildedConfiguration.multipleRcptto)
		assert.False(t, buildedConfiguration.multipleMessageReceiving)
		assert.False(t, buildedConfiguration.logServerActivity)
		assert.Equal(t, defaultGreetingMsg, buildedConfiguration.msgGreeting)
		assert.Equal(t, defaultInvalidCmdMsg, buildedConfiguration.msgInvalidCmd)
		assert.Equal(t, defaultQuitMsg, buildedConfiguration.msgQuitCmd)
		assert.Equal(t, defaultSessionTimeout, buildedConfiguration.sessionTimeout)
		assert.Equal(t, defaultShutdownTimeout, buildedConfiguration.shutdownTimeout)

		assert.Equal(t, defaultInvalidCmdHeloSequenceMsg, buildedConfiguration.msgInvalidCmdHeloSequence)
		assert.Equal(t, defaultInvalidCmdHeloArgMsg, buildedConfiguration.msgInvalidCmdHeloArg)
		assert.Equal(t, defaultTransientNegativeMsg, buildedConfiguration.msgHeloBlacklistedDomain)
		assert.Equal(t, defaultReceivedMsg, buildedConfiguration.msgHeloReceived)

		assert.Equal(t, defaultInvalidCmdMailfromSequenceMsg, buildedConfiguration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, defaultInvalidCmdMailfromArgMsg, buildedConfiguration.msgInvalidCmdMailfromArg)
		assert.Equal(t, defaultTransientNegativeMsg, buildedConfiguration.msgMailfromBlacklistedEmail)
		assert.Equal(t, defaultReceivedMsg, buildedConfiguration.msgMailfromReceived)

		assert.Equal(t, defaultInvalidCmdRcpttoSequenceMsg, buildedConfiguration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, defaultInvalidCmdRcpttoArgMsg, buildedConfiguration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, defaultTransientNegativeMsg, buildedConfiguration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, defaultNotRegistredRcpttoEmailMsg, buildedConfiguration.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, defaultReceivedMsg, buildedConfiguration.msgRcpttoReceived)

		assert.Equal(t, defaultInvalidCmdDataSequenceMsg, buildedConfiguration.msgInvalidCmdDataSequence)
		assert.Equal(t, defaultReadyForReceiveMsg, buildedConfiguration.msgDataReceived)

		assert.Equal(t, defaultInvalidCmdHeloSequenceMsg, buildedConfiguration.msgInvalidCmdRsetSequence)
		assert.Equal(t, defaultInvalidCmdMsg, buildedConfiguration.msgInvalidCmdRsetArg)
		assert.Equal(t, defaultOkMsg, buildedConfiguration.msgRsetReceived)

		assert.Equal(t, defaultOkMsg, buildedConfiguration.msgNoopReceived)

		assert.Equal(t, fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", defaultMessageSizeLimit), buildedConfiguration.msgMsgSizeIsTooBig)
		assert.Equal(t, defaultReceivedMsg, buildedConfiguration.msgMsgReceived)
		assert.Equal(t, defaultMessageSizeLimit, buildedConfiguration.msgSizeLimit)

		assert.Empty(t, buildedConfiguration.blacklistedHeloDomains)
		assert.Empty(t, buildedConfiguration.blacklistedMailfromEmails)
		assert.Empty(t, buildedConfiguration.blacklistedRcpttoEmails)
		assert.Empty(t, buildedConfiguration.notRegisteredEmails)

		assert.Equal(t, defaultSessionResponseDelay, buildedConfiguration.responseDelayHelo)
		assert.Equal(t, defaultSessionResponseDelay, buildedConfiguration.responseDelayMailfrom)
		assert.Equal(t, defaultSessionResponseDelay, buildedConfiguration.responseDelayRcptto)
		assert.Equal(t, defaultSessionResponseDelay, buildedConfiguration.responseDelayData)
		assert.Equal(t, defaultSessionResponseDelay, buildedConfiguration.responseDelayMessage)
		assert.Equal(t, defaultSessionResponseDelay, buildedConfiguration.responseDelayRset)
		assert.Equal(t, defaultSessionResponseDelay, buildedConfiguration.responseDelayNoop)
		assert.Equal(t, defaultSessionResponseDelay, buildedConfiguration.responseDelayQuit)
	})

	t.Run("creates new configuration with custom settings", func(t *testing.T) {
		configAttr := ConfigurationAttr{
			HostAddress:                   "hostAddress",
			PortNumber:                    25,
			LogToStdout:                   true,
			LogServerActivity:             true,
			IsCmdFailFast:                 true,
			MultipleRcptto:                true,
			MultipleMessageReceiving:      true,
			MsgGreeting:                   "msgGreeting",
			MsgInvalidCmd:                 "msgInvalidCmd",
			MsgQuitCmd:                    "msgQuitCmd",
			MsgInvalidCmdHeloSequence:     "msgInvalidCmdHeloSequence",
			MsgInvalidCmdHeloArg:          "msgInvalidCmdHeloArg",
			MsgHeloBlacklistedDomain:      "msgHeloBlacklistedDomain",
			MsgHeloReceived:               "msgHeloReceived",
			MsgInvalidCmdMailfromSequence: "msgInvalidCmdMailfromSequence",
			MsgInvalidCmdMailfromArg:      "msgInvalidCmdMailfromArg",
			MsgMailfromBlacklistedEmail:   "msgMailfromBlacklistedEmail",
			MsgMailfromReceived:           "msgMailfromReceived",
			MsgInvalidCmdRcpttoSequence:   "msgInvalidCmdRcpttoSequence",
			MsgInvalidCmdRcpttoArg:        "msgInvalidCmdRcpttoArg",
			MsgRcpttoNotRegisteredEmail:   "msgRcpttoNotRegisteredEmail",
			MsgRcpttoBlacklistedEmail:     "msgRcpttoBlacklistedEmail",
			MsgRcpttoReceived:             "msgRcpttoReceived",
			MsgInvalidCmdDataSequence:     "msgInvalidCmdDataSequence",
			MsgDataReceived:               "msgDataReceived",
			MsgMsgSizeIsTooBig:            emptyString,
			MsgMsgReceived:                "msgMsgReceived",
			MsgInvalidCmdRsetSequence:     "msgInvalidCmdRsetSequence",
			MsgInvalidCmdRsetArg:          "msgInvalidCmdRsetArg",
			MsgRsetReceived:               "msgRsetReceived",
			MsgNoopReceived:               "msgNoopReceived",
			BlacklistedHeloDomains:        []string{},
			BlacklistedMailfromEmails:     []string{},
			NotRegisteredEmails:           []string{},
			BlacklistedRcpttoEmails:       []string{},
			ResponseDelayHelo:             2,
			ResponseDelayMailfrom:         2,
			ResponseDelayRcptto:           2,
			ResponseDelayData:             2,
			ResponseDelayMessage:          2,
			ResponseDelayRset:             2,
			ResponseDelayNoop:             2,
			ResponseDelayQuit:             2,
			MsgSizeLimit:                  42,
			SessionTimeout:                120,
			ShutdownTimeout:               2,
		}
		buildedConfiguration := newConfiguration(configAttr)

		assert.Equal(t, configAttr.HostAddress, buildedConfiguration.hostAddress)
		assert.Equal(t, configAttr.PortNumber, buildedConfiguration.portNumber)
		assert.Equal(t, configAttr.LogToStdout, buildedConfiguration.logToStdout)
		assert.Equal(t, configAttr.IsCmdFailFast, buildedConfiguration.isCmdFailFast)
		assert.Equal(t, configAttr.MultipleRcptto, buildedConfiguration.multipleRcptto)
		assert.Equal(t, configAttr.MultipleMessageReceiving, buildedConfiguration.multipleMessageReceiving)
		assert.Equal(t, configAttr.LogServerActivity, buildedConfiguration.logServerActivity)
		assert.Equal(t, configAttr.MsgGreeting, buildedConfiguration.msgGreeting)
		assert.Equal(t, configAttr.MsgInvalidCmd, buildedConfiguration.msgInvalidCmd)
		assert.Equal(t, configAttr.MsgQuitCmd, buildedConfiguration.msgQuitCmd)
		assert.Equal(t, configAttr.SessionTimeout, buildedConfiguration.sessionTimeout)
		assert.Equal(t, configAttr.ShutdownTimeout, buildedConfiguration.shutdownTimeout)

		assert.Equal(t, configAttr.MsgInvalidCmdHeloSequence, buildedConfiguration.msgInvalidCmdHeloSequence)
		assert.Equal(t, configAttr.MsgInvalidCmdHeloArg, buildedConfiguration.msgInvalidCmdHeloArg)
		assert.Equal(t, configAttr.MsgHeloBlacklistedDomain, buildedConfiguration.msgHeloBlacklistedDomain)
		assert.Equal(t, configAttr.MsgHeloReceived, buildedConfiguration.msgHeloReceived)

		assert.Equal(t, configAttr.MsgInvalidCmdMailfromSequence, buildedConfiguration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, configAttr.MsgInvalidCmdMailfromArg, buildedConfiguration.msgInvalidCmdMailfromArg)
		assert.Equal(t, configAttr.MsgMailfromBlacklistedEmail, buildedConfiguration.msgMailfromBlacklistedEmail)
		assert.Equal(t, configAttr.MsgMailfromReceived, buildedConfiguration.msgMailfromReceived)

		assert.Equal(t, configAttr.MsgInvalidCmdRcpttoSequence, buildedConfiguration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, configAttr.MsgInvalidCmdRcpttoArg, buildedConfiguration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, configAttr.MsgRcpttoBlacklistedEmail, buildedConfiguration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, configAttr.MsgRcpttoNotRegisteredEmail, buildedConfiguration.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, configAttr.MsgRcpttoReceived, buildedConfiguration.msgRcpttoReceived)

		assert.Equal(t, configAttr.MsgInvalidCmdDataSequence, buildedConfiguration.msgInvalidCmdDataSequence)
		assert.Equal(t, configAttr.MsgDataReceived, buildedConfiguration.msgDataReceived)

		assert.Equal(t, configAttr.MsgInvalidCmdRsetSequence, buildedConfiguration.msgInvalidCmdRsetSequence)
		assert.Equal(t, configAttr.MsgInvalidCmdRsetArg, buildedConfiguration.msgInvalidCmdRsetArg)
		assert.Equal(t, configAttr.MsgRsetReceived, buildedConfiguration.msgRsetReceived)

		assert.Equal(t, configAttr.MsgNoopReceived, buildedConfiguration.msgNoopReceived)

		assert.Equal(t, fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", configAttr.MsgSizeLimit), buildedConfiguration.msgMsgSizeIsTooBig)
		assert.Equal(t, configAttr.MsgMsgReceived, buildedConfiguration.msgMsgReceived)
		assert.Equal(t, configAttr.MsgSizeLimit, buildedConfiguration.msgSizeLimit)

		assert.Equal(t, configAttr.BlacklistedHeloDomains, buildedConfiguration.blacklistedHeloDomains)
		assert.Equal(t, configAttr.BlacklistedMailfromEmails, buildedConfiguration.blacklistedMailfromEmails)
		assert.Equal(t, configAttr.BlacklistedRcpttoEmails, buildedConfiguration.blacklistedRcpttoEmails)
		assert.Equal(t, configAttr.NotRegisteredEmails, buildedConfiguration.notRegisteredEmails)

		assert.Equal(t, configAttr.ResponseDelayHelo, buildedConfiguration.responseDelayHelo)
		assert.Equal(t, configAttr.ResponseDelayMailfrom, buildedConfiguration.responseDelayMailfrom)
		assert.Equal(t, configAttr.ResponseDelayRcptto, buildedConfiguration.responseDelayRcptto)
		assert.Equal(t, configAttr.ResponseDelayData, buildedConfiguration.responseDelayData)
		assert.Equal(t, configAttr.ResponseDelayMessage, buildedConfiguration.responseDelayMessage)
		assert.Equal(t, configAttr.ResponseDelayRset, buildedConfiguration.responseDelayRset)
		assert.Equal(t, configAttr.ResponseDelayNoop, buildedConfiguration.responseDelayNoop)
		assert.Equal(t, configAttr.ResponseDelayQuit, buildedConfiguration.responseDelayQuit)
	})
}

func TestConfigurationAttrAssignDefaultValues(t *testing.T) {
	t.Run("assignes default values", func(t *testing.T) {
		configurationAttr := new(ConfigurationAttr)
		configurationAttr.assignDefaultValues()

		assert.Equal(t, defaultHostAddress, configurationAttr.HostAddress)
		assert.Equal(t, defaultGreetingMsg, configurationAttr.MsgGreeting)
		assert.Equal(t, defaultInvalidCmdMsg, configurationAttr.MsgInvalidCmd)
		assert.Equal(t, defaultQuitMsg, configurationAttr.MsgQuitCmd)
		assert.Equal(t, defaultSessionTimeout, configurationAttr.SessionTimeout)
		assert.Equal(t, defaultShutdownTimeout, configurationAttr.ShutdownTimeout)

		assert.Equal(t, defaultInvalidCmdHeloSequenceMsg, configurationAttr.MsgInvalidCmdHeloSequence)
		assert.Equal(t, defaultInvalidCmdHeloArgMsg, configurationAttr.MsgInvalidCmdHeloArg)
		assert.Equal(t, defaultTransientNegativeMsg, configurationAttr.MsgHeloBlacklistedDomain)
		assert.Equal(t, defaultReceivedMsg, configurationAttr.MsgHeloReceived)

		assert.Equal(t, defaultInvalidCmdMailfromSequenceMsg, configurationAttr.MsgInvalidCmdMailfromSequence)
		assert.Equal(t, defaultInvalidCmdMailfromArgMsg, configurationAttr.MsgInvalidCmdMailfromArg)
		assert.Equal(t, defaultTransientNegativeMsg, configurationAttr.MsgMailfromBlacklistedEmail)
		assert.Equal(t, defaultReceivedMsg, configurationAttr.MsgMailfromReceived)

		assert.Equal(t, defaultInvalidCmdRcpttoSequenceMsg, configurationAttr.MsgInvalidCmdRcpttoSequence)
		assert.Equal(t, defaultInvalidCmdRcpttoArgMsg, configurationAttr.MsgInvalidCmdRcpttoArg)
		assert.Equal(t, defaultTransientNegativeMsg, configurationAttr.MsgRcpttoBlacklistedEmail)
		assert.Equal(t, defaultNotRegistredRcpttoEmailMsg, configurationAttr.MsgRcpttoNotRegisteredEmail)
		assert.Equal(t, defaultReceivedMsg, configurationAttr.MsgRcpttoReceived)

		assert.Equal(t, defaultInvalidCmdDataSequenceMsg, configurationAttr.MsgInvalidCmdDataSequence)
		assert.Equal(t, defaultReadyForReceiveMsg, configurationAttr.MsgDataReceived)

		assert.Equal(t, defaultInvalidCmdHeloSequenceMsg, configurationAttr.MsgInvalidCmdRsetSequence)
		assert.Equal(t, defaultInvalidCmdMsg, configurationAttr.MsgInvalidCmdRsetArg)
		assert.Equal(t, defaultOkMsg, configurationAttr.MsgRsetReceived)

		assert.Equal(t, defaultOkMsg, configurationAttr.MsgNoopReceived)

		assert.Equal(t, fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", defaultMessageSizeLimit), configurationAttr.MsgMsgSizeIsTooBig)
		assert.Equal(t, defaultReceivedMsg, configurationAttr.MsgMsgReceived)
		assert.Equal(t, defaultMessageSizeLimit, configurationAttr.MsgSizeLimit)
	})
}
