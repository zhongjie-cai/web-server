package webserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_NonSupportedLogLevels(t *testing.T) {
	// arrange
	var unsupportedValue = maxLogLevel

	// SUT
	var sut = LogLevel(unsupportedValue)

	// act
	var result = sut.String()

	// assert
	assert.Equal(t, debugLogLevelName, result)
}

func TestString_SupportedLogLevel(t *testing.T) {
	// SUT
	var sut = LogLevelError

	// act
	var result = sut.String()

	// assert
	assert.Equal(t, errorLogLevelName, result)
}

func TestNewLogLevel_NoMatchFound(t *testing.T) {
	// arrange
	var dummyValue = "some value"

	// SUT + act
	var result = NewLogLevel(dummyValue)

	// assert
	assert.Equal(t, LogLevelDebug, result)
}

func TestNewLogLevel_HappyPath(t *testing.T) {
	for key, value := range logLevelNameMapping {
		// SUT + act
		var result = NewLogLevel(key)

		// assert
		assert.Equal(t, value, result)
	}
}
