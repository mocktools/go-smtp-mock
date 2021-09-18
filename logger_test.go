package smtpmock

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Run("when log to stdout enabled", func(t *testing.T) {
		isLogToStdoutEnabled := true
		logger := &eventLogger{logToStdout: isLogToStdoutEnabled, flag: LogFlag}

		assert.EqualValues(t, logger, newLogger(isLogToStdoutEnabled))
	})

	t.Run("when log to stdout disabled", func(t *testing.T) {
		isLogToStdoutEnabled := false
		logger := &eventLogger{logToStdout: isLogToStdoutEnabled, flag: LogFlag}

		assert.EqualValues(t, logger, newLogger(isLogToStdoutEnabled))
	})
}

func TestEventLoggerInfo(t *testing.T) {
	logMessage := "log message"

	t.Run("when log to stdout enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(true)
		logger.stdout = &buf
		logger.info(logMessage) // initializes and memoizes INFO logger during first function calling

		assert.Regexp(t, loggerMessageRegex(InfoLogLevel, logMessage), buf.String())
		assert.NotNil(t, logger.eventInfo)
		memoizedInfoLogger := logger.eventInfo
		logger.info(logMessage)
		assert.Same(t, memoizedInfoLogger, logger.eventInfo)
	})

	t.Run("when log to stdout disabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(false)
		logger.stdout = &buf
		logger.info(logMessage)

		assert.Nil(t, logger.eventInfo)
		assert.Empty(t, buf.String())
	})
}

func TestEventLoggerWarning(t *testing.T) {
	logMessage := "log message"

	t.Run("when log to stdout enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(true)
		logger.stdout = &buf
		logger.warning(logMessage) // initializes and memoizes WARNING logger during first function calling

		assert.Regexp(t, loggerMessageRegex(WarningLogLevel, logMessage), buf.String())
		assert.NotNil(t, logger.eventWarning)
		memoizedWarningLogger := logger.eventWarning
		logger.warning(logMessage)
		assert.Same(t, memoizedWarningLogger, logger.eventWarning)
	})

	t.Run("when log to stdout disabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(false)
		logger.stdout = &buf
		logger.warning(logMessage)

		assert.Nil(t, logger.eventWarning)
		assert.Empty(t, buf.String())
	})
}

func TestEventLoggerError(t *testing.T) {
	logMessage := "log message"

	t.Run("when log to stdout enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(true)
		logger.stderr = &buf
		logger.error(logMessage) // initializes and memoizes ERROR logger during first function calling

		assert.Regexp(t, loggerMessageRegex(ErrorLogLevel, logMessage), buf.String())
		assert.NotNil(t, logger.eventError)
		memoizedErrorLogger := logger.eventError
		logger.error(logMessage)
		assert.Same(t, memoizedErrorLogger, logger.eventError)
	})

	t.Run("when log to stdout disabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(false)
		logger.stderr = &buf
		logger.error(logMessage)

		assert.Nil(t, logger.eventError)
		assert.Empty(t, buf.String())
	})
}
