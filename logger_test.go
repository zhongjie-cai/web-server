package webserver

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

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
}

func TestPrepareLoggingFunc_HappyPath(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyLogType = LogType(rand.Intn(100))
	var dummyLogLevel = LogLevel(rand.Intn(100))
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "%v %v %v"
	var dummyParameter1 = "some parameter 1"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("some parameter 3")
	var dummyDescription = fmt.Sprintf(dummyMessageFormat, dummyParameter1, dummyParameter2, dummyParameter3)

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "Log", 1, func(self *DefaultCustomization, session Session, logType LogType, logLevel LogLevel, category, subcategory, description string) {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyLogType, logType)
		assert.Equal(t, dummyLogLevel, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyDescription, description)
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(prepareLogging, 1, func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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
	})

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
}
