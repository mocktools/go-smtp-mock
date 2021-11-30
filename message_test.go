package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageIsCleared(t *testing.T) {
	t.Run("when cleared status true", func(t *testing.T) {
		message := &message{cleared: true}

		assert.True(t, message.isCleared())
	})

	t.Run("when cleared status false", func(t *testing.T) {
		assert.False(t, new(message).isCleared())
	})
}

func TestMessagesAppend(t *testing.T) {
	t.Run("addes message pointer into items slice", func(t *testing.T) {
		message, messages := new(message), new(messages)
		messages.append(message)

		assert.Same(t, message, messages.items[0])
	})
}
