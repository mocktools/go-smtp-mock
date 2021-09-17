package smtpmock

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegex(t *testing.T) {
	t.Run("valid regex pattern", func(t *testing.T) {
		regexPattern := EmptyString
		actualRegex, err := newRegex(regexPattern)
		expectedRegex, _ := regexp.Compile(regexPattern)

		assert.Equal(t, expectedRegex, actualRegex)
		assert.NoError(t, err)
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		actualRegex, err := newRegex(`\K`)

		assert.Nil(t, actualRegex)
		assert.Error(t, err)
	})
}

func TestMatchRegex(t *testing.T) {
	t.Run("valid regex pattern, matched string", func(t *testing.T) {
		assert.True(t, matchRegex(EmptyString, EmptyString))
	})

	t.Run("valid regex pattern, not matched string", func(t *testing.T) {
		assert.False(t, matchRegex("42", `\D+`))
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		assert.False(t, matchRegex(EmptyString, `\K`))
	})
}

func TestRegexCaptureGroup(t *testing.T) {
	str, regexPattern := "abbc", `\A(a)(b{2}).+\z`

	t.Run("returns string when regex capture group found", func(t *testing.T) {
		assert.Equal(t, "bb", regexCaptureGroup(str, regexPattern, 2))
	})

	t.Run("panics when regex capture group not found", func(t *testing.T) {
		assert.Panics(t, func() { regexCaptureGroup(str, regexPattern, 3) })
	})
}

func TestIsIncluded(t *testing.T) {
	var item string

	t.Run("item found in slice", func(t *testing.T) {
		assert.True(t, isIncluded([]string{item}, item))
	})

	t.Run("item not found in slice", func(t *testing.T) {
		assert.False(t, isIncluded([]string{}, item))
	})
}
