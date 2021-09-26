package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerIsFailFastScenario(t *testing.T) {
	t.Run("when fail fast scenario enabled", func(t *testing.T) {
		handler := &handler{configuration: &configuration{isCmdFailFast: true}}

		assert.True(t, handler.isFailFastScenario())
	})

	t.Run("when fail fast scenario disabled", func(t *testing.T) {
		handler := &handler{configuration: new(configuration)}

		assert.False(t, handler.isFailFastScenario())
	})
}
