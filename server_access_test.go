package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerAccess(t *testing.T) {

	configuration := createConfiguration()
	server := newServer(configuration)

	server.Clear()

	assert.Equal(t, 0, server.Count())
	assert.Empty(t, server.Message(0))
	assert.Empty(t, server.Message(1))

	assert.NoError(t, server.Start())

	_ = runMinimalSuccessfulSMTPSession(configuration.hostAddress, server.PortNumber)

	assert.Equal(t, 1, server.Count())
	assert.NotEmpty(t, server.Message(0))
	assert.Empty(t, server.Message(1))

	assert.NotEmpty(t, server.messages)
	assert.NotNil(t, server.quit)
	assert.NotNil(t, server.quitTimeout)
	assert.True(t, server.isStarted)
	assert.Greater(t, server.PortNumber, 0)

	_ = server.Stop()

	server.Clear()

	assert.Equal(t, 0, server.Count())
	assert.Empty(t, server.Message(0))
	assert.Empty(t, server.Message(1))
}
