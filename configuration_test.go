package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	t.Run("creates new configuration with default settings", func(t *testing.T) {
		buildedConfiguration := NewConfiguration(ConfigurationAttr{})

		assert.False(t, buildedConfiguration.logToStdout)
		assert.False(t, buildedConfiguration.isCmdFailFast)
		assert.Equal(t, DefaultGreetingMsg, buildedConfiguration.msgGreeting)
		assert.Equal(t, DefaultInvalidCmdMsg, buildedConfiguration.msgInvalidCmd)
		assert.Equal(t, DefaultQuitMsg, buildedConfiguration.msgQuitCmd)

		assert.Equal(t, DefaultInvalidCmdHeloSequenceMsg, buildedConfiguration.msgInvalidCmdHeloSequence)
		assert.Equal(t, DefaultInvalidCmdHeloArgMsg, buildedConfiguration.msgInvalidCmdHeloArg)
		assert.Equal(t, DefaultQuitMsg, buildedConfiguration.msgHeloBlacklistedDomain)
		assert.Equal(t, DefaultReceivedMsg, buildedConfiguration.msgHeloReceived)

		assert.Equal(t, DefaultInvalidCmdMailfromSequenceMsg, buildedConfiguration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, DefaultInvalidCmdMailfromArgMsg, buildedConfiguration.msgInvalidCmdMailfromArg)
		assert.Equal(t, DefaultQuitMsg, buildedConfiguration.msgMailfromBlacklistedEmail)
		assert.Equal(t, DefaultReceivedMsg, buildedConfiguration.msgMailfromReceived)

		assert.Equal(t, DefaultInvalidCmdRcpttoSequenceMsg, buildedConfiguration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, DefaultInvalidCmdRcpttoArgMsg, buildedConfiguration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, DefaultNotRegistredRcpttoEmailMsg, buildedConfiguration.msgRcpttoNotRegistredEmail)
		assert.Equal(t, DefaultQuitMsg, buildedConfiguration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, DefaultReceivedMsg, buildedConfiguration.msgRcpttoReceived)

		assert.Empty(t, buildedConfiguration.blacklistedHeloDomains)
		assert.Empty(t, buildedConfiguration.blacklistedMailfromEmails)
		assert.Empty(t, buildedConfiguration.registredRcpttoEmails)
		assert.Empty(t, buildedConfiguration.blacklistedRcpttoEmails)
	})

	t.Run("creates new configuration with custom settings", func(t *testing.T) {
		configAttr := ConfigurationAttr{
			logToStdout:                   true,
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
			msgRcpttoNotRegistredEmail:    "msgRcpttoNotRegistredEmail",
			msgRcpttoBlacklistedEmail:     "msgRcpttoBlacklistedEmail",
			msgRcpttoReceived:             "msgRcpttoReceived",
			blacklistedHeloDomains:        []string{},
			blacklistedMailfromEmails:     []string{},
			registredRcpttoEmails:         []string{},
			blacklistedRcpttoEmails:       []string{},
		}
		buildedConfiguration := NewConfiguration(configAttr)

		assert.Equal(t, configAttr.logToStdout, buildedConfiguration.logToStdout)
		assert.Equal(t, configAttr.isCmdFailFast, buildedConfiguration.isCmdFailFast)
		assert.Equal(t, configAttr.msgGreeting, buildedConfiguration.msgGreeting)
		assert.Equal(t, configAttr.msgInvalidCmd, buildedConfiguration.msgInvalidCmd)
		assert.Equal(t, configAttr.msgQuitCmd, buildedConfiguration.msgQuitCmd)

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
		assert.Equal(t, configAttr.msgRcpttoNotRegistredEmail, buildedConfiguration.msgRcpttoNotRegistredEmail)
		assert.Equal(t, configAttr.msgRcpttoBlacklistedEmail, buildedConfiguration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, configAttr.msgRcpttoReceived, buildedConfiguration.msgRcpttoReceived)

		assert.Equal(t, configAttr.blacklistedHeloDomains, buildedConfiguration.blacklistedHeloDomains)
		assert.Equal(t, configAttr.blacklistedMailfromEmails, buildedConfiguration.blacklistedMailfromEmails)
		assert.Equal(t, configAttr.registredRcpttoEmails, buildedConfiguration.registredRcpttoEmails)
		assert.Equal(t, configAttr.blacklistedRcpttoEmails, buildedConfiguration.blacklistedRcpttoEmails)
	})
}

func TestConfigurationAttrAssignDefaultValues(t *testing.T) {
	t.Run("assignes default values", func(t *testing.T) {
		configurationAttr := new(ConfigurationAttr)
		configurationAttr.assignDefaultValues()

		assert.Equal(t, DefaultGreetingMsg, configurationAttr.msgGreeting)
		assert.Equal(t, DefaultInvalidCmdMsg, configurationAttr.msgInvalidCmd)
		assert.Equal(t, DefaultQuitMsg, configurationAttr.msgQuitCmd)

		assert.Equal(t, DefaultInvalidCmdHeloSequenceMsg, configurationAttr.msgInvalidCmdHeloSequence)
		assert.Equal(t, DefaultInvalidCmdHeloArgMsg, configurationAttr.msgInvalidCmdHeloArg)
		assert.Equal(t, DefaultQuitMsg, configurationAttr.msgHeloBlacklistedDomain)
		assert.Equal(t, DefaultReceivedMsg, configurationAttr.msgHeloReceived)

		assert.Equal(t, DefaultInvalidCmdMailfromSequenceMsg, configurationAttr.msgInvalidCmdMailfromSequence)
		assert.Equal(t, DefaultInvalidCmdMailfromArgMsg, configurationAttr.msgInvalidCmdMailfromArg)
		assert.Equal(t, DefaultQuitMsg, configurationAttr.msgMailfromBlacklistedEmail)
		assert.Equal(t, DefaultReceivedMsg, configurationAttr.msgMailfromReceived)

		assert.Equal(t, DefaultInvalidCmdRcpttoSequenceMsg, configurationAttr.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, DefaultInvalidCmdRcpttoArgMsg, configurationAttr.msgInvalidCmdRcpttoArg)
		assert.Equal(t, DefaultNotRegistredRcpttoEmailMsg, configurationAttr.msgRcpttoNotRegistredEmail)
		assert.Equal(t, DefaultQuitMsg, configurationAttr.msgRcpttoBlacklistedEmail)
		assert.Equal(t, DefaultReceivedMsg, configurationAttr.msgRcpttoReceived)
	})
}
