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
		assert.Equal(t, defaultPortNuber, buildedConfiguration.portNumber)
		assert.False(t, buildedConfiguration.logToStdout)
		assert.False(t, buildedConfiguration.isCmdFailFast)
		assert.False(t, buildedConfiguration.logServerActivity)
		assert.Equal(t, defaultGreetingMsg, buildedConfiguration.msgGreeting)
		assert.Equal(t, defaultInvalidCmdMsg, buildedConfiguration.msgInvalidCmd)
		assert.Equal(t, defaultQuitMsg, buildedConfiguration.msgQuitCmd)
		assert.Equal(t, defaultSessionTimeout, buildedConfiguration.sessionTimeout)

		assert.Equal(t, defaultInvalidCmdHeloSequenceMsg, buildedConfiguration.msgInvalidCmdHeloSequence)
		assert.Equal(t, defaultInvalidCmdHeloArgMsg, buildedConfiguration.msgInvalidCmdHeloArg)
		assert.Equal(t, defaultQuitMsg, buildedConfiguration.msgHeloBlacklistedDomain)
		assert.Equal(t, defaultReceivedMsg, buildedConfiguration.msgHeloReceived)

		assert.Equal(t, defaultInvalidCmdMailfromSequenceMsg, buildedConfiguration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, defaultInvalidCmdMailfromArgMsg, buildedConfiguration.msgInvalidCmdMailfromArg)
		assert.Equal(t, defaultQuitMsg, buildedConfiguration.msgMailfromBlacklistedEmail)
		assert.Equal(t, defaultReceivedMsg, buildedConfiguration.msgMailfromReceived)

		assert.Equal(t, defaultInvalidCmdRcpttoSequenceMsg, buildedConfiguration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, defaultInvalidCmdRcpttoArgMsg, buildedConfiguration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, defaultQuitMsg, buildedConfiguration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, defaultNotRegistredRcpttoEmailMsg, buildedConfiguration.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, defaultReceivedMsg, buildedConfiguration.msgRcpttoReceived)

		assert.Equal(t, defaultInvalidCmdDataSequenceMsg, buildedConfiguration.msgInvalidCmdDataSequence)
		assert.Equal(t, defaultReadyForReceiveMsg, buildedConfiguration.msgDataReceived)

		assert.Equal(t, fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", defaultMessageSizeLimit), buildedConfiguration.msgMsgSizeIsTooBig)
		assert.Equal(t, defaultReceivedMsg, buildedConfiguration.msgMsgReceived)
		assert.Equal(t, defaultMessageSizeLimit, buildedConfiguration.msqSizeLimit)

		assert.Empty(t, buildedConfiguration.blacklistedHeloDomains)
		assert.Empty(t, buildedConfiguration.blacklistedMailfromEmails)
		assert.Empty(t, buildedConfiguration.blacklistedRcpttoEmails)
		assert.Empty(t, buildedConfiguration.notRegisteredEmails)
	})

	t.Run("creates new configuration with custom settings", func(t *testing.T) {
		configAttr := ConfigurationAttr{
			hostAddress:                   "hostAddress",
			portNumber:                    25,
			logToStdout:                   true,
			logServerActivity:             true,
			isCmdFailFast:                 true,
			msgGreeting:                   "msgGreeting",
			msgInvalidCmd:                 "msgInvalidCmd",
			msgQuitCmd:                    "msgQuitCmd",
			msgInvalidCmdHeloSequence:     "msgInvalidCmdHeloSequence",
			msgInvalidCmdHeloArg:          "msgInvalidCmdHeloArg",
			msgHeloBlacklistedDomain:      "msgHeloBlacklistedDomain",
			msgHeloReceived:               "msgHeloReceived",
			msgInvalidCmdMailfromSequence: "msgInvalidCmdMailfromSequence",
			msgInvalidCmdMailfromArg:      "msgInvalidCmdMailfromArg",
			msgMailfromBlacklistedEmail:   "msgMailfromBlacklistedEmail",
			msgMailfromReceived:           "msgMailfromReceived",
			msgInvalidCmdRcpttoSequence:   "msgInvalidCmdRcpttoSequence",
			msgInvalidCmdRcpttoArg:        "msgInvalidCmdRcpttoArg",
			msgRcpttoNotRegisteredEmail:   "msgRcpttoNotRegisteredEmail",
			msgRcpttoBlacklistedEmail:     "msgRcpttoBlacklistedEmail",
			msgRcpttoReceived:             "msgRcpttoReceived",
			msgInvalidCmdDataSequence:     "msgInvalidCmdDataSequence",
			msgDataReceived:               "msgDataReceived",
			msgMsgSizeIsTooBig:            emptyString,
			msgMsgReceived:                "msgMsgReceived",
			blacklistedHeloDomains:        []string{},
			blacklistedMailfromEmails:     []string{},
			notRegisteredEmails:           []string{},
			blacklistedRcpttoEmails:       []string{},
			msqSizeLimit:                  42,
			sessionTimeout:                120,
		}
		buildedConfiguration := newConfiguration(configAttr)

		assert.Equal(t, configAttr.hostAddress, buildedConfiguration.hostAddress)
		assert.Equal(t, configAttr.portNumber, buildedConfiguration.portNumber)
		assert.Equal(t, configAttr.logToStdout, buildedConfiguration.logToStdout)
		assert.Equal(t, configAttr.isCmdFailFast, buildedConfiguration.isCmdFailFast)
		assert.Equal(t, configAttr.logServerActivity, buildedConfiguration.logServerActivity)
		assert.Equal(t, configAttr.msgGreeting, buildedConfiguration.msgGreeting)
		assert.Equal(t, configAttr.msgInvalidCmd, buildedConfiguration.msgInvalidCmd)
		assert.Equal(t, configAttr.msgQuitCmd, buildedConfiguration.msgQuitCmd)
		assert.Equal(t, configAttr.sessionTimeout, buildedConfiguration.sessionTimeout)

		assert.Equal(t, configAttr.msgInvalidCmdHeloSequence, buildedConfiguration.msgInvalidCmdHeloSequence)
		assert.Equal(t, configAttr.msgInvalidCmdHeloArg, buildedConfiguration.msgInvalidCmdHeloArg)
		assert.Equal(t, configAttr.msgHeloBlacklistedDomain, buildedConfiguration.msgHeloBlacklistedDomain)
		assert.Equal(t, configAttr.msgHeloReceived, buildedConfiguration.msgHeloReceived)

		assert.Equal(t, configAttr.msgInvalidCmdMailfromSequence, buildedConfiguration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, configAttr.msgInvalidCmdMailfromArg, buildedConfiguration.msgInvalidCmdMailfromArg)
		assert.Equal(t, configAttr.msgMailfromBlacklistedEmail, buildedConfiguration.msgMailfromBlacklistedEmail)
		assert.Equal(t, configAttr.msgMailfromReceived, buildedConfiguration.msgMailfromReceived)

		assert.Equal(t, configAttr.msgInvalidCmdRcpttoSequence, buildedConfiguration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, configAttr.msgInvalidCmdRcpttoArg, buildedConfiguration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, configAttr.msgRcpttoBlacklistedEmail, buildedConfiguration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, configAttr.msgRcpttoNotRegisteredEmail, buildedConfiguration.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, configAttr.msgRcpttoReceived, buildedConfiguration.msgRcpttoReceived)

		assert.Equal(t, configAttr.msgInvalidCmdDataSequence, buildedConfiguration.msgInvalidCmdDataSequence)
		assert.Equal(t, configAttr.msgDataReceived, buildedConfiguration.msgDataReceived)

		assert.Equal(t, fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", configAttr.msqSizeLimit), buildedConfiguration.msgMsgSizeIsTooBig)
		assert.Equal(t, configAttr.msgMsgReceived, buildedConfiguration.msgMsgReceived)
		assert.Equal(t, configAttr.msqSizeLimit, buildedConfiguration.msqSizeLimit)

		assert.Equal(t, configAttr.blacklistedHeloDomains, buildedConfiguration.blacklistedHeloDomains)
		assert.Equal(t, configAttr.blacklistedMailfromEmails, buildedConfiguration.blacklistedMailfromEmails)
		assert.Equal(t, configAttr.blacklistedRcpttoEmails, buildedConfiguration.blacklistedRcpttoEmails)
		assert.Equal(t, configAttr.notRegisteredEmails, buildedConfiguration.notRegisteredEmails)
	})
}

func TestConfigurationAttrAssignDefaultValues(t *testing.T) {
	t.Run("assignes default values", func(t *testing.T) {
		configurationAttr := new(ConfigurationAttr)
		configurationAttr.assignDefaultValues()

		assert.Equal(t, defaultHostAddress, configurationAttr.hostAddress)
		assert.Equal(t, defaultPortNuber, configurationAttr.portNumber)
		assert.Equal(t, defaultGreetingMsg, configurationAttr.msgGreeting)
		assert.Equal(t, defaultInvalidCmdMsg, configurationAttr.msgInvalidCmd)
		assert.Equal(t, defaultQuitMsg, configurationAttr.msgQuitCmd)
		assert.Equal(t, defaultSessionTimeout, configurationAttr.sessionTimeout)

		assert.Equal(t, defaultInvalidCmdHeloSequenceMsg, configurationAttr.msgInvalidCmdHeloSequence)
		assert.Equal(t, defaultInvalidCmdHeloArgMsg, configurationAttr.msgInvalidCmdHeloArg)
		assert.Equal(t, defaultQuitMsg, configurationAttr.msgHeloBlacklistedDomain)
		assert.Equal(t, defaultReceivedMsg, configurationAttr.msgHeloReceived)

		assert.Equal(t, defaultInvalidCmdMailfromSequenceMsg, configurationAttr.msgInvalidCmdMailfromSequence)
		assert.Equal(t, defaultInvalidCmdMailfromArgMsg, configurationAttr.msgInvalidCmdMailfromArg)
		assert.Equal(t, defaultQuitMsg, configurationAttr.msgMailfromBlacklistedEmail)
		assert.Equal(t, defaultReceivedMsg, configurationAttr.msgMailfromReceived)

		assert.Equal(t, defaultInvalidCmdRcpttoSequenceMsg, configurationAttr.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, defaultInvalidCmdRcpttoArgMsg, configurationAttr.msgInvalidCmdRcpttoArg)
		assert.Equal(t, defaultQuitMsg, configurationAttr.msgRcpttoBlacklistedEmail)
		assert.Equal(t, defaultNotRegistredRcpttoEmailMsg, configurationAttr.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, defaultReceivedMsg, configurationAttr.msgRcpttoReceived)

		assert.Equal(t, defaultInvalidCmdDataSequenceMsg, configurationAttr.msgInvalidCmdDataSequence)
		assert.Equal(t, defaultReadyForReceiveMsg, configurationAttr.msgDataReceived)

		assert.Equal(t, fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", defaultMessageSizeLimit), configurationAttr.msgMsgSizeIsTooBig)
		assert.Equal(t, defaultReceivedMsg, configurationAttr.msgMsgReceived)
		assert.Equal(t, defaultMessageSizeLimit, configurationAttr.msqSizeLimit)
	})
}
