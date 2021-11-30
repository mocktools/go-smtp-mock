package smtpmock

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("creates new server with default configuration settings", func(t *testing.T) {
		server := New(ConfigurationAttr{})
		configuration := server.configuration

		assert.Equal(t, DefaultHostAddress, configuration.hostAddress)
		assert.Equal(t, DefaultPortNuber, configuration.portNumber)
		assert.False(t, configuration.logToStdout)
		assert.False(t, configuration.isCmdFailFast)
		assert.False(t, configuration.logServerActivity)
		assert.Equal(t, DefaultGreetingMsg, configuration.msgGreeting)
		assert.Equal(t, DefaultInvalidCmdMsg, configuration.msgInvalidCmd)
		assert.Equal(t, DefaultQuitMsg, configuration.msgQuitCmd)
		assert.Equal(t, DefaultSessionTimeout, configuration.sessionTimeout)

		assert.Equal(t, DefaultInvalidCmdHeloSequenceMsg, configuration.msgInvalidCmdHeloSequence)
		assert.Equal(t, DefaultInvalidCmdHeloArgMsg, configuration.msgInvalidCmdHeloArg)
		assert.Equal(t, DefaultQuitMsg, configuration.msgHeloBlacklistedDomain)
		assert.Equal(t, DefaultReceivedMsg, configuration.msgHeloReceived)

		assert.Equal(t, DefaultInvalidCmdMailfromSequenceMsg, configuration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, DefaultInvalidCmdMailfromArgMsg, configuration.msgInvalidCmdMailfromArg)
		assert.Equal(t, DefaultQuitMsg, configuration.msgMailfromBlacklistedEmail)
		assert.Equal(t, DefaultReceivedMsg, configuration.msgMailfromReceived)

		assert.Equal(t, DefaultInvalidCmdRcpttoSequenceMsg, configuration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, DefaultInvalidCmdRcpttoArgMsg, configuration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, DefaultQuitMsg, configuration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, DefaultNotRegistredRcpttoEmailMsg, configuration.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, DefaultReceivedMsg, configuration.msgRcpttoReceived)

		assert.Equal(t, DefaultInvalidCmdDataSequenceMsg, configuration.msgInvalidCmdDataSequence)
		assert.Equal(t, DefaultReadyForReceiveMsg, configuration.msgDataReceived)

		assert.Equal(t, fmt.Sprintf(DefaultMsgSizeIsTooBigMsg+" %d bytes", DefaultMessageSizeLimit), configuration.msgMsgSizeIsTooBig)
		assert.Equal(t, DefaultReceivedMsg, configuration.msgMsgReceived)
		assert.Equal(t, DefaultMessageSizeLimit, configuration.msqSizeLimit)

		assert.Empty(t, configuration.blacklistedHeloDomains)
		assert.Empty(t, configuration.blacklistedMailfromEmails)
		assert.Empty(t, configuration.blacklistedRcpttoEmails)
		assert.Empty(t, configuration.notRegisteredEmails)

		assert.Empty(t, server.messages)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.wg)
		assert.Nil(t, server.quit)
		assert.False(t, server.isStarted)
	})

	t.Run("creates new server with custom configuration settings", func(t *testing.T) {
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
			msgMsgSizeIsTooBig:            EmptyString,
			msgMsgReceived:                "msgMsgReceived",
			blacklistedHeloDomains:        []string{},
			blacklistedMailfromEmails:     []string{},
			notRegisteredEmails:           []string{},
			blacklistedRcpttoEmails:       []string{},
			msqSizeLimit:                  42,
			sessionTimeout:                120,
		}
		server := New(configAttr)
		configuration := server.configuration

		assert.Equal(t, configAttr.hostAddress, configuration.hostAddress)
		assert.Equal(t, configAttr.portNumber, configuration.portNumber)
		assert.Equal(t, configAttr.logToStdout, configuration.logToStdout)
		assert.Equal(t, configAttr.isCmdFailFast, configuration.isCmdFailFast)
		assert.Equal(t, configAttr.logServerActivity, configuration.logServerActivity)
		assert.Equal(t, configAttr.msgGreeting, configuration.msgGreeting)
		assert.Equal(t, configAttr.msgInvalidCmd, configuration.msgInvalidCmd)
		assert.Equal(t, configAttr.msgQuitCmd, configuration.msgQuitCmd)
		assert.Equal(t, configAttr.sessionTimeout, configuration.sessionTimeout)

		assert.Equal(t, configAttr.msgInvalidCmdHeloSequence, configuration.msgInvalidCmdHeloSequence)
		assert.Equal(t, configAttr.msgInvalidCmdHeloArg, configuration.msgInvalidCmdHeloArg)
		assert.Equal(t, configAttr.msgHeloBlacklistedDomain, configuration.msgHeloBlacklistedDomain)
		assert.Equal(t, configAttr.msgHeloReceived, configuration.msgHeloReceived)

		assert.Equal(t, configAttr.msgInvalidCmdMailfromSequence, configuration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, configAttr.msgInvalidCmdMailfromArg, configuration.msgInvalidCmdMailfromArg)
		assert.Equal(t, configAttr.msgMailfromBlacklistedEmail, configuration.msgMailfromBlacklistedEmail)
		assert.Equal(t, configAttr.msgMailfromReceived, configuration.msgMailfromReceived)

		assert.Equal(t, configAttr.msgInvalidCmdRcpttoSequence, configuration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, configAttr.msgInvalidCmdRcpttoArg, configuration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, configAttr.msgRcpttoBlacklistedEmail, configuration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, configAttr.msgRcpttoNotRegisteredEmail, configuration.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, configAttr.msgRcpttoReceived, configuration.msgRcpttoReceived)

		assert.Equal(t, configAttr.msgInvalidCmdDataSequence, configuration.msgInvalidCmdDataSequence)
		assert.Equal(t, configAttr.msgDataReceived, configuration.msgDataReceived)

		assert.Equal(t, fmt.Sprintf(DefaultMsgSizeIsTooBigMsg+" %d bytes", configAttr.msqSizeLimit), configuration.msgMsgSizeIsTooBig)
		assert.Equal(t, configAttr.msgMsgReceived, configuration.msgMsgReceived)
		assert.Equal(t, configAttr.msqSizeLimit, configuration.msqSizeLimit)

		assert.Equal(t, configAttr.blacklistedHeloDomains, configuration.blacklistedHeloDomains)
		assert.Equal(t, configAttr.blacklistedMailfromEmails, configuration.blacklistedMailfromEmails)
		assert.Equal(t, configAttr.blacklistedRcpttoEmails, configuration.blacklistedRcpttoEmails)
		assert.Equal(t, configAttr.notRegisteredEmails, configuration.notRegisteredEmails)

		assert.Empty(t, server.messages)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.wg)
		assert.Nil(t, server.quit)
		assert.False(t, server.isStarted)
	})

	t.Run("successful iteration with new server", func(t *testing.T) {
		server := New(ConfigurationAttr{})
		configuration := server.configuration

		assert.Empty(t, server.messages)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.wg)
		assert.Nil(t, server.quit)
		assert.False(t, server.isStarted)

		assert.NoError(t, server.Start())
		assert.True(t, server.isStarted)
		_ = runMinimalSuccessfulSmtpSession(configuration.hostAddress, configuration.portNumber)
		_ = server.Stop()

		assert.NotEmpty(t, server.messages)
		assert.NotNil(t, server.quit)
		assert.False(t, server.isStarted)
	})
}
