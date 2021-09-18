package smtpmock

import (
	"io"
	"log"
)

// Logger interface
type logger interface {
	info(string)
	warning(string)
	error(string)
}

// Custom logger that supports 3 different log levels (info, warning, error)
type eventLogger struct {
	eventInfo, eventWarning, eventError *log.Logger
	logToStdout                         bool
	flag                                int
	stdout, stderr                      io.Writer
}

// Logger builder. Returns pointer to builded new logger structure
func newLogger(logToStdout bool) *eventLogger {
	return &eventLogger{
		logToStdout: logToStdout,
		flag:        LogFlag,
	}
}

// logger methods

// Provides INFO log level. Writes to stdout for case when logger.logToStdout is enabled,
// suppressed otherwise
func (logger *eventLogger) info(message string) {
	if logger.logToStdout {
		if logger.eventInfo == nil {
			logger.eventInfo = log.New(logger.stdout, InfoLogLevel+": ", logger.flag)
		}

		logger.eventInfo.Println(message)
	}
}

// Provides WARNING log level. Writes to stdout for case when logger.logToStdout is enabled,
// suppressed otherwise
func (logger *eventLogger) warning(message string) {
	if logger.logToStdout {
		if logger.eventWarning == nil {
			logger.eventWarning = log.New(logger.stdout, WarningLogLevel+": ", logger.flag)
		}

		logger.eventWarning.Println(message)
	}
}

// Provides ERROR log level. Writes to stdout for case when logger.logToStdout is enabled,
// suppressed otherwise
func (logger *eventLogger) error(message string) {
	if logger.logToStdout {
		if logger.eventError == nil {
			logger.eventError = log.New(logger.stderr, ErrorLogLevel+": ", logger.flag)
		}

		logger.eventError.Println(message)
	}
}
