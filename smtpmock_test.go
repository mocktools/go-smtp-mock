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

		assert.Equal(t, defaultHostAddress, configuration.hostAddress)
		assert.False(t, configuration.logToStdout)
		assert.False(t, configuration.isCmdFailFast)
		assert.False(t, configuration.logServerActivity)
		assert.Equal(t, defaultGreetingMsg, configuration.msgGreeting)
		assert.Equal(t, defaultInvalidCmdMsg, configuration.msgInvalidCmd)
		assert.Equal(t, defaultQuitMsg, configuration.msgQuitCmd)
		assert.Equal(t, defaultSessionTimeout, configuration.sessionTimeout)

		assert.Equal(t, defaultInvalidCmdHeloSequenceMsg, configuration.msgInvalidCmdHeloSequence)
		assert.Equal(t, defaultInvalidCmdHeloArgMsg, configuration.msgInvalidCmdHeloArg)
		assert.Equal(t, defaultTransientNegativeMsg, configuration.msgHeloBlacklistedDomain)
		assert.Equal(t, defaultReceivedMsg, configuration.msgHeloReceived)

		assert.Equal(t, defaultInvalidCmdMailfromSequenceMsg, configuration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, defaultInvalidCmdMailfromArgMsg, configuration.msgInvalidCmdMailfromArg)
		assert.Equal(t, defaultTransientNegativeMsg, configuration.msgMailfromBlacklistedEmail)
		assert.Equal(t, defaultReceivedMsg, configuration.msgMailfromReceived)

		assert.Equal(t, defaultInvalidCmdRcpttoSequenceMsg, configuration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, defaultInvalidCmdRcpttoArgMsg, configuration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, defaultTransientNegativeMsg, configuration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, defaultNotRegistredRcpttoEmailMsg, configuration.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, defaultReceivedMsg, configuration.msgRcpttoReceived)

		assert.Equal(t, defaultInvalidCmdDataSequenceMsg, configuration.msgInvalidCmdDataSequence)
		assert.Equal(t, defaultReadyForReceiveMsg, configuration.msgDataReceived)

		assert.Equal(t, fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", defaultMessageSizeLimit), configuration.msgMsgSizeIsTooBig)
		assert.Equal(t, defaultReceivedMsg, configuration.msgMsgReceived)
		assert.Equal(t, defaultMessageSizeLimit, configuration.msqSizeLimit)

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
			HostAddress:                   "hostAddress",
			PortNumber:                    25,
			LogToStdout:                   true,
			LogServerActivity:             true,
			IsCmdFailFast:                 true,
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
			BlacklistedHeloDomains:        []string{},
			BlacklistedMailfromEmails:     []string{},
			NotRegisteredEmails:           []string{},
			BlacklistedRcpttoEmails:       []string{},
			MsqSizeLimit:                  42,
			SessionTimeout:                120,
		}
		server := New(configAttr)
		configuration := server.configuration

		assert.Equal(t, configAttr.HostAddress, configuration.hostAddress)
		assert.Equal(t, configAttr.PortNumber, configuration.portNumber)
		assert.Equal(t, configAttr.LogToStdout, configuration.logToStdout)
		assert.Equal(t, configAttr.IsCmdFailFast, configuration.isCmdFailFast)
		assert.Equal(t, configAttr.LogServerActivity, configuration.logServerActivity)
		assert.Equal(t, configAttr.MsgGreeting, configuration.msgGreeting)
		assert.Equal(t, configAttr.MsgInvalidCmd, configuration.msgInvalidCmd)
		assert.Equal(t, configAttr.MsgQuitCmd, configuration.msgQuitCmd)
		assert.Equal(t, configAttr.SessionTimeout, configuration.sessionTimeout)

		assert.Equal(t, configAttr.MsgInvalidCmdHeloSequence, configuration.msgInvalidCmdHeloSequence)
		assert.Equal(t, configAttr.MsgInvalidCmdHeloArg, configuration.msgInvalidCmdHeloArg)
		assert.Equal(t, configAttr.MsgHeloBlacklistedDomain, configuration.msgHeloBlacklistedDomain)
		assert.Equal(t, configAttr.MsgHeloReceived, configuration.msgHeloReceived)

		assert.Equal(t, configAttr.MsgInvalidCmdMailfromSequence, configuration.msgInvalidCmdMailfromSequence)
		assert.Equal(t, configAttr.MsgInvalidCmdMailfromArg, configuration.msgInvalidCmdMailfromArg)
		assert.Equal(t, configAttr.MsgMailfromBlacklistedEmail, configuration.msgMailfromBlacklistedEmail)
		assert.Equal(t, configAttr.MsgMailfromReceived, configuration.msgMailfromReceived)

		assert.Equal(t, configAttr.MsgInvalidCmdRcpttoSequence, configuration.msgInvalidCmdRcpttoSequence)
		assert.Equal(t, configAttr.MsgInvalidCmdRcpttoArg, configuration.msgInvalidCmdRcpttoArg)
		assert.Equal(t, configAttr.MsgRcpttoBlacklistedEmail, configuration.msgRcpttoBlacklistedEmail)
		assert.Equal(t, configAttr.MsgRcpttoNotRegisteredEmail, configuration.msgRcpttoNotRegisteredEmail)
		assert.Equal(t, configAttr.MsgRcpttoReceived, configuration.msgRcpttoReceived)

		assert.Equal(t, configAttr.MsgInvalidCmdDataSequence, configuration.msgInvalidCmdDataSequence)
		assert.Equal(t, configAttr.MsgDataReceived, configuration.msgDataReceived)

		assert.Equal(t, fmt.Sprintf(defaultMsgSizeIsTooBigMsg+" %d bytes", configAttr.MsqSizeLimit), configuration.msgMsgSizeIsTooBig)
		assert.Equal(t, configAttr.MsgMsgReceived, configuration.msgMsgReceived)
		assert.Equal(t, configAttr.MsqSizeLimit, configuration.msqSizeLimit)

		assert.Equal(t, configAttr.BlacklistedHeloDomains, configuration.blacklistedHeloDomains)
		assert.Equal(t, configAttr.BlacklistedMailfromEmails, configuration.blacklistedMailfromEmails)
		assert.Equal(t, configAttr.BlacklistedRcpttoEmails, configuration.blacklistedRcpttoEmails)
		assert.Equal(t, configAttr.NotRegisteredEmails, configuration.notRegisteredEmails)

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
		_ = runMinimalSuccessfulSMTPSession(configuration.hostAddress, server.PortNumber)
		_ = server.Stop()

		assert.NotEmpty(t, server.messages)
		assert.NotNil(t, server.quit)
		assert.False(t, server.isStarted)
		assert.Greater(t, server.PortNumber, 0)
	})
}
