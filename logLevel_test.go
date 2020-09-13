package webserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_NonSupportedLogLevels(t *testing.T) {
	// arrange
	var unsupportedValue = maxLogLevel

	// mock
	createMock(t)

	// SUT
	var sut = LogLevel(unsupportedValue)

	// act
	var result = sut.String()

	// assert
	assert.Equal(t, debugLogLevelName, result)

	// verify
	verifyAll(t)
}

func TestString_SupportedLogLevel(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var sut = LogLevelError

	// act
	var result = sut.String()

	// assert
	assert.Equal(t, errorLogLevelName, result)

	// verify
	verifyAll(t)
}

func TestNewLogLevel_NoMatchFound(t *testing.T) {
	// arrange
	var dummyValue = "some value"

	// mock
	createMock(t)

	// SUT + act
	var result = NewLogLevel(dummyValue)

	// assert
	assert.Equal(t, LogLevelDebug, result)

	// tear down
	verifyAll(t)
}

func TestNewLogLevel_HappyPath(t *testing.T) {
	for key, value := range logLevelNameMapping {
		// mock
		createMock(t)

		// SUT + act
		var result = NewLogLevel(key)

		// assert
		assert.Equal(t, value, result)

		// tear down
		verifyAll(t)
	}
}
