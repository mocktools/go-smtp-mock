package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessagesAppend(t *testing.T) {
	t.Run("addes message pointer into items slice", func(t *testing.T) {
		message, messages := new(message), new(messages)
		messages.append(message)

		assert.Same(t, message, messages.items[0])
	})
}
