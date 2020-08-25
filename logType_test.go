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
	assert.Equal(t, AppRoot, sut)
	assert.Equal(t, appRootName, result)

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
	var sut = MethodLogic

	// act
	var result = sut.String()

	// assert
	assert.Equal(t, methodLogicName, result)

	// verify
	verifyAll(t)
}

func TestString_MultipleSupportedLogTypes(t *testing.T) {
	// arrange
	var supportedValue = EndpointEnter | EndpointRequest | MethodLogic | EndpointResponse | EndpointExit

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
	assert.Equal(t, GeneralLogging, sut)
	assert.True(t, strings.Contains(result, apiEnterName))
	assert.True(t, strings.Contains(result, apiRequestName))
	assert.True(t, strings.Contains(result, methodLogicName))
	assert.True(t, strings.Contains(result, apiResponseName))
	assert.True(t, strings.Contains(result, apiExitName))

	// verify
	verifyAll(t)
}

func TestHasFlag_FlagMatch_AppRoot(t *testing.T) {
	// arrange
	var flag = AppRoot

	// mock
	createMock(t)

	// SUT
	var sut = AppRoot

	// act
	var result = sut.HasFlag(flag)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestHasFlag_FlagNoMatch_AppRoot(t *testing.T) {
	// arrange
	var flag = AppRoot

	// mock
	createMock(t)

	// SUT
	var sut = EndpointEnter | EndpointExit

	// act
	var result = sut.HasFlag(flag)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestHasFlag_FlagMatch_NotAppRoot(t *testing.T) {
	// arrange
	var flag = MethodLogic

	// mock
	createMock(t)

	// SUT
	var sut = EndpointEnter | MethodLogic | EndpointExit

	// act
	var result = sut.HasFlag(flag)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestHasFlag_FlagNoMatch_NotAppRoot(t *testing.T) {
	// arrange
	var flag = MethodLogic

	// mock
	createMock(t)

	// SUT
	var sut = EndpointEnter | EndpointExit

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
	assert.Equal(t, AppRoot, result)

	// tear down
	verifyAll(t)
}

func TestNewLogType_AppRoot(t *testing.T) {
	// arrange
	var dummyValue = appRootName

	// mock
	createMock(t)

	// expect
	stringsSplit = func(s, sep string) []string {
		return strings.Split(s, sep)
	}

	// SUT + act
	var result = NewLogType(dummyValue)

	// assert
	assert.Equal(t, AppRoot, result)

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
