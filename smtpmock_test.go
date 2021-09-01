package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicFunction(t *testing.T) { // TODO: remove during implementation
	t.Run("returns 42", func(t *testing.T) {
		assert.Equal(t, 42, PublicFunction())
	})
}
