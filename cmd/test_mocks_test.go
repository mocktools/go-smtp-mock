package main

import (
	"github.com/stretchr/testify/mock"
)

// log mock
type logMock struct {
	mock.Mock
}

func (logMock *logMock) Fatalf(format string, v ...interface{}) {
	logMock.Called(format, v)
}
