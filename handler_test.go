package smtpmock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerClearError(t *testing.T) {
	t.Run("erases session error", func(t *testing.T) {
		session := &session{err: errors.New("some error")}
		handler := &handler{configuration: createConfiguration(), session: session}
		handler.clearError()

		assert.Nil(t, session.err)
	})
}
