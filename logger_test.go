package smtpmock

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Run("when log to stdout, server activity enabled", func(t *testing.T) {
		isLogToStdoutEnabled, logServerActivity := true, true
		logger := &eventLogger{
			logToStdout:       isLogToStdoutEnabled,
			logServerActivity: logServerActivity,
			flag:              logFlag,
			stdout:            os.Stdout,
			stderr:            os.Stderr,
		}

		assert.EqualValues(t, logger, newLogger(isLogToStdoutEnabled, logServerActivity))
	})

	t.Run("when log to stdout, server activity disabled", func(t *testing.T) {
		isLogToStdoutEnabled, logServerActivity := false, false
		logger := &eventLogger{
			flag:   logFlag,
			stdout: os.Stdout,
			stderr: os.Stderr,
		}

		assert.EqualValues(t, logger, newLogger(isLogToStdoutEnabled, logServerActivity))
	})
}

func TestEventLoggerInfoActivity(t *testing.T) {
	logMessage := "log message"

	t.Run("when log to stdout and server activity enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(true, true)
		logger.stdout = &buf
		logger.InfoActivity(logMessage) // initializes and memoizes INFO logger during first function calling

		assert.Regexp(t, loggerMessageRegex(infoLogLevel, logMessage), buf.String())
		assert.NotNil(t, logger.eventInfo)
		memoizedInfoLogger := logger.eventInfo
		logger.InfoActivity(logMessage)
		assert.Same(t, memoizedInfoLogger, logger.eventInfo)
	})

	t.Run("when log to stdout disabled, server activity enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(false, true)
		logger.stdout = &buf
		logger.InfoActivity(logMessage)

		assert.Nil(t, logger.eventInfo)
		assert.Empty(t, buf.String())
	})

	t.Run("when log to stdout enabled, server activity disabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(true, false)
		logger.stdout = &buf
		logger.InfoActivity(logMessage)

		assert.Nil(t, logger.eventInfo)
		assert.Empty(t, buf.String())
	})
}

func TestEventLoggerInfo(t *testing.T) {
	logMessage := "log message"

	t.Run("when log to stdout enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(true, false)
		logger.stdout = &buf
		logger.Info(logMessage) // initializes and memoizes INFO logger during first function calling

		assert.Regexp(t, loggerMessageRegex(infoLogLevel, logMessage), buf.String())
		assert.NotNil(t, logger.eventInfo)
		memoizedInfoLogger := logger.eventInfo
		logger.Info(logMessage)
		assert.Same(t, memoizedInfoLogger, logger.eventInfo)
	})

	t.Run("when log to stdout disabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(false, false)
		logger.stdout = &buf
		logger.Info(logMessage)

		assert.Nil(t, logger.eventInfo)
		assert.Empty(t, buf.String())
	})
}

func TestEventLoggerWarning(t *testing.T) {
	logMessage := "log message"

	t.Run("when log to stdout enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(true, false)
		logger.stdout = &buf
		logger.Warning(logMessage) // initializes and memoizes WARNING logger during first function calling

		assert.Regexp(t, loggerMessageRegex(warningLogLevel, logMessage), buf.String())
		assert.NotNil(t, logger.eventWarning)
		memoizedWarningLogger := logger.eventWarning
		logger.Warning(logMessage)
		assert.Same(t, memoizedWarningLogger, logger.eventWarning)
	})

	t.Run("when log to stdout disabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(false, false)
		logger.stdout = &buf
		logger.Warning(logMessage)

		assert.Nil(t, logger.eventWarning)
		assert.Empty(t, buf.String())
	})
}

func TestEventLoggerError(t *testing.T) {
	logMessage := "log message"

	t.Run("when log to stdout enabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(true, false)
		logger.stderr = &buf
		logger.Error(logMessage) // initializes and memoizes ERROR logger during first function calling

		assert.Regexp(t, loggerMessageRegex(errorLogLevel, logMessage), buf.String())
		assert.NotNil(t, logger.eventError)
		memoizedErrorLogger := logger.eventError
		logger.Error(logMessage)
		assert.Same(t, memoizedErrorLogger, logger.eventError)
	})

	t.Run("when log to stdout disabled", func(t *testing.T) {
		var buf bytes.Buffer
		logger := newLogger(false, false)
		logger.stderr = &buf
		logger.Error(logMessage)

		assert.Nil(t, logger.eventError)
		assert.Empty(t, buf.String())
	})
}
