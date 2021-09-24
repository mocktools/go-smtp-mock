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
		assert.Equal(t, DefaultReceivedMsg, buildedConfiguration.msgHeloReceived)
		assert.Equal(t, DefaultInvalidCmdHeloArgMsg, buildedConfiguration.msgInvalidCmdHeloArg)
		assert.Equal(t, DefaultInvalidCmdMsg, buildedConfiguration.msgInvalidCmd)
		assert.Equal(t, DefaultInvalidCmdHeloSequenceMsg, buildedConfiguration.msgInvalidCmdHeloSequence)
		assert.Equal(t, DefaultQuitMsg, buildedConfiguration.msgQuit)
		assert.Empty(t, buildedConfiguration.blacklistedHeloDomains)
	})

	t.Run("creates new configuration with custom settings", func(t *testing.T) {
		configAttr := ConfigurationAttr{
			logToStdout:               true,
			isCmdFailFast:             true,
			msgGreeting:               "a",
			msgHeloReceived:           "b",
			msgInvalidCmdHeloArg:      "c",
			msgInvalidCmd:             "d",
			msgInvalidCmdHeloSequence: "e",
			msgQuit:                   "f",
			blacklistedHeloDomains:    []string{},
		}
		buildedConfiguration := NewConfiguration(configAttr)

		assert.Equal(t, configAttr.logToStdout, buildedConfiguration.logToStdout)
		assert.Equal(t, configAttr.isCmdFailFast, buildedConfiguration.isCmdFailFast)
		assert.Equal(t, configAttr.msgGreeting, buildedConfiguration.msgGreeting)
		assert.Equal(t, configAttr.msgHeloReceived, buildedConfiguration.msgHeloReceived)
		assert.Equal(t, configAttr.msgInvalidCmdHeloArg, buildedConfiguration.msgInvalidCmdHeloArg)
		assert.Equal(t, configAttr.msgInvalidCmd, buildedConfiguration.msgInvalidCmd)
		assert.Equal(t, configAttr.msgInvalidCmdHeloSequence, buildedConfiguration.msgInvalidCmdHeloSequence)
		assert.Equal(t, configAttr.msgQuit, buildedConfiguration.msgQuit)
		assert.Equal(t, configAttr.blacklistedHeloDomains, buildedConfiguration.blacklistedHeloDomains)
	})
}

func TestConfigurationAttrAssignDefaultValues(t *testing.T) {
	t.Run("assignes default values", func(t *testing.T) {
		configurationAttr := new(ConfigurationAttr)
		configurationAttr.assignDefaultValues()

		assert.Equal(t, DefaultGreetingMsg, configurationAttr.msgGreeting)
		assert.Equal(t, DefaultReceivedMsg, configurationAttr.msgHeloReceived)
		assert.Equal(t, DefaultInvalidCmdHeloArgMsg, configurationAttr.msgInvalidCmdHeloArg)
		assert.Equal(t, DefaultInvalidCmdMsg, configurationAttr.msgInvalidCmd)
		assert.Equal(t, DefaultInvalidCmdHeloSequenceMsg, configurationAttr.msgInvalidCmdHeloSequence)
		assert.Equal(t, DefaultQuitMsg, configurationAttr.msgQuit)
	})
}
