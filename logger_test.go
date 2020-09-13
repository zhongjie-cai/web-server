package webserver

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type dummyCustomizationLoger struct {
	dummyCustomization
	log func(session Session, logType LogType, logLevel LogLevel, category, subcategory, description string)
}

func (customization *dummyCustomizationLoger) Log(session Session, logType LogType, logLevel LogLevel, category, subcategory, description string) {
	if customization.log != nil {
		customization.log(session, logType, logLevel, category, subcategory, description)
		return
	}
	assert.Fail(customization.t, "Unexpected call to Log")
}

func TestPrepareLoggingFunc_NilSession(t *testing.T) {
	// arrange
	var dummySession *session
	var dummyLogType = LogType(rand.Intn(100))
	var dummyLogLevel = LogLevel(rand.Intn(100))
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// SUT + act
	prepareLogging(
		dummySession,
		dummyLogType,
		dummyLogLevel,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestPrepareLoggingFunc_HappyPath(t *testing.T) {
	// arrange
	var dummyCustomizationLoger = &dummyCustomizationLoger{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationLoger,
	}
	var dummyLogType = LogType(rand.Intn(100))
	var dummyLogLevel = LogLevel(rand.Intn(100))
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")
	var dummyDescription = "some description"
	var customizationLogExpected int
	var customizationLogCalled int

	// mock
	createMock(t)

	// expect
	fmtSprintfExpected = 1
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		assert.Equal(t, dummyMessageFormat, format)
		assert.Equal(t, 3, len(a))
		assert.Equal(t, dummyParameter1, a[0])
		assert.Equal(t, dummyParameter2, a[1])
		assert.Equal(t, dummyParameter3, a[2])
		return dummyDescription
	}
	customizationLogExpected = 1
	dummyCustomizationLoger.log = func(session Session, logType LogType, logLevel LogLevel, category, subcategory, description string) {
		customizationLogCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyLogType, logType)
		assert.Equal(t, dummyLogLevel, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyDescription, description)
	}

	// SUT + act
	prepareLogging(
		dummySession,
		dummyLogType,
		dummyLogLevel,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationLogExpected, customizationLogCalled, "Unexpected number of calls to method customization.Log")
}

func TestLogAppRoot(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeAppRoot, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logAppRoot(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogEndpointEnter(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeEndpointEnter, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logEndpointEnter(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogEndpointRequest(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeEndpointRequest, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logEndpointRequest(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogMethodEnter(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeMethodEnter, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logMethodEnter(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogMethodParameter(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeMethodParameter, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logMethodParameter(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogMethodLogic(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyLogLevel = LogLevel(rand.Intn(100))
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeMethodLogic, logType)
		assert.Equal(t, dummyLogLevel, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logMethodLogic(
		dummySession,
		dummyLogLevel,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogWebcallStart(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeWebcallStart, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logWebcallStart(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogWebcallRequest(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeWebcallRequest, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logWebcallRequest(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogWebcallResponse(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeWebcallResponse, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logWebcallResponse(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogWebcallFinish(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeWebcallFinish, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logWebcallFinish(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogMethodReturn(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeMethodReturn, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logMethodReturn(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogMethodExit(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeMethodExit, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logMethodExit(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogEndpointResponse(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeEndpointResponse, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logEndpointResponse(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}

func TestLogEndpointExit(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")

	// mock
	createMock(t)

	// expect
	prepareLoggingFuncExpected = 1
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, LogTypeEndpointExit, logType)
		assert.Equal(t, LogLevelInfo, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

	// SUT + act
	logEndpointExit(
		dummySession,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
}
