package webserver

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_AppRoot(t *testing.T) {
	// arrange
	var appRootValue = 0

	// mock
	createMock(t)

	// SUT
	var sut = LogType(appRootValue)

	// act
	var result = sut.String()

	// assert
	assert.Equal(t, LogTypeAppRoot, sut)
	assert.Equal(t, appRootLogTypeName, result)

	// verify
	verifyAll(t)
}

func TestString_NonSupportedLogTypes(t *testing.T) {
	// arrange
	var unsupportedValue = 1 << 31

	// mock
	createMock(t)

	// expect
	sortStringsExpected = 1
	sortStrings = func(a []string) {
		sortStringsCalled++
		sort.Strings(a)
	}
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

	// SUT
	var sut = LogType(unsupportedValue)

	// act
	var result = sut.String()

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestString_SingleSupportedLogType(t *testing.T) {
	// mock
	createMock(t)

	// expect
	sortStringsExpected = 1
	sortStrings = func(a []string) {
		sortStringsCalled++
		sort.Strings(a)
	}
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

	// SUT
	var sut = LogTypeMethodLogic

	// act
	var result = sut.String()

	// assert
	assert.Equal(t, methodLogicLogTypeName, result)

	// verify
	verifyAll(t)
}

func TestString_MultipleSupportedLogTypes(t *testing.T) {
	// arrange
	var supportedValue = LogTypeEndpointEnter | LogTypeEndpointRequest | LogTypeMethodLogic | LogTypeEndpointResponse | LogTypeEndpointExit

	// mock
	createMock(t)

	// expect
	sortStringsExpected = 1
	sortStrings = func(a []string) {
		sortStringsCalled++
		sort.Strings(a)
	}
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

	// SUT
	var sut = LogType(supportedValue)

	// act
	var result = sut.String()

	// assert
	assert.Equal(t, LogTypeGeneralLogging, sut)
	assert.True(t, strings.Contains(result, apiEnterLogTypeName))
	assert.True(t, strings.Contains(result, apiRequestLogTypeName))
	assert.True(t, strings.Contains(result, methodLogicLogTypeName))
	assert.True(t, strings.Contains(result, apiResponseLogTypeName))
	assert.True(t, strings.Contains(result, apiExitLogTypeName))

	// verify
	verifyAll(t)
}

func TestHasFlag_FlagMatch_AppRoot(t *testing.T) {
	// arrange
	var flag = LogTypeAppRoot

	// mock
	createMock(t)

	// SUT
	var sut = LogTypeAppRoot

	// act
	var result = sut.HasFlag(flag)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestHasFlag_FlagNoMatch_AppRoot(t *testing.T) {
	// arrange
	var flag = LogTypeAppRoot

	// mock
	createMock(t)

	// SUT
	var sut = LogTypeEndpointEnter | LogTypeEndpointExit

	// act
	var result = sut.HasFlag(flag)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestHasFlag_FlagMatch_NotAppRoot(t *testing.T) {
	// arrange
	var flag = LogTypeMethodLogic

	// mock
	createMock(t)

	// SUT
	var sut = LogTypeEndpointEnter | LogTypeMethodLogic | LogTypeEndpointExit

	// act
	var result = sut.HasFlag(flag)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestHasFlag_FlagNoMatch_NotAppRoot(t *testing.T) {
	// arrange
	var flag = LogTypeMethodLogic

	// mock
	createMock(t)

	// SUT
	var sut = LogTypeEndpointEnter | LogTypeEndpointExit

	// act
	var result = sut.HasFlag(flag)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}

func TestNewLogType_NoMatchFound(t *testing.T) {
	// arrange
	var dummyValue = "some value"

	// mock
	createMock(t)

	// expect
	stringsSplit = func(s, sep string) []string {
		return strings.Split(s, sep)
	}

	// SUT + act
	var result = NewLogType(dummyValue)

	// assert
	assert.Equal(t, LogTypeAppRoot, result)

	// tear down
	verifyAll(t)
}

func TestNewLogType_AppRoot(t *testing.T) {
	// arrange
	var dummyValue = appRootLogTypeName

	// mock
	createMock(t)

	// expect
	stringsSplit = func(s, sep string) []string {
		return strings.Split(s, sep)
	}

	// SUT + act
	var result = NewLogType(dummyValue)

	// assert
	assert.Equal(t, LogTypeAppRoot, result)

	// tear down
	verifyAll(t)
}

func TestNewLogType_HappyPath(t *testing.T) {
	for key, value := range logTypeNameMapping {
		// mock
		createMock(t)

		// expect
		stringsSplit = func(s, sep string) []string {
			return strings.Split(s, sep)
		}

		// SUT + act
		var result = NewLogType(key)

		// assert
		assert.Equal(t, value, result)

		// tear down
		verifyAll(t)
	}
}
