package webserver

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime/debug"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func TestDefaultCustomization_PreBootstrap(t *testing.T) {
	// SUT + act
	var err = customizationDefault.PreBootstrap()

	// assert
	assert.NoError(t, err)
}

func TestDefaultCustomization_PostBootstrap(t *testing.T) {
	// SUT + act
	var err = customizationDefault.PostBootstrap()

	// assert
	assert.NoError(t, err)
}

func TestDefaultCustomization_AppClosing(t *testing.T) {
	// SUT + act
	var err = customizationDefault.AppClosing()

	// assert
	assert.NoError(t, err)
}

func TestDefaultCustomization_Log_NilSession(t *testing.T) {
	// arrange
	var dummySession *session
	var dummyLogType = LogType(rand.Intn(100))
	var dummyLogLevel = LogLevel(rand.Intn(100))
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyDescription = "some description"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummySession, i)
		return true
	})

	// SUT + act
	customizationDefault.Log(
		dummySession,
		dummyLogType,
		dummyLogLevel,
		dummyCategory,
		dummySubcategory,
		dummyDescription,
	)
}

func TestDefaultCustomization_Log_HappyPath(t *testing.T) {
	// arrange
	var dummySession = &session{}
	var dummyLogType = LogType(rand.Intn(100))
	var dummyLogLevel = LogLevel(rand.Intn(100))
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyDescription = "some description"
	var dummySessionID = uuid.New()
	var dummySessionName = "some session name"
	var dummyFormat = "<%v|%v> (%v|%v) [%v|%v] %v\n"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummySession, i)
		return false
	})
	m.ExpectMethod(dummySession, "GetID", 1, func() uuid.UUID {
		return dummySessionID
	})
	m.ExpectMethod(dummySession, "GetName", 1, func() string {
		return dummySessionName
	})
	m.ExpectFunc(fmt.Printf, 1, func(format string, a ...interface{}) (n int, err error) {
		assert.Equal(t, dummyFormat, format)
		assert.Equal(t, 7, len(a))
		assert.Equal(t, dummySessionID, a[0])
		assert.Equal(t, dummySessionName, a[1])
		assert.Equal(t, dummyLogType, a[2])
		assert.Equal(t, dummyLogLevel, a[3])
		assert.Equal(t, dummyCategory, a[4])
		assert.Equal(t, dummySubcategory, a[5])
		assert.Equal(t, dummyDescription, a[6])
		return rand.Int(), errors.New("some error")
	})

	// SUT + act
	customizationDefault.Log(
		dummySession,
		dummyLogType,
		dummyLogLevel,
		dummyCategory,
		dummySubcategory,
		dummyDescription,
	)
}

func TestDefaultCustomization_ServerCert(t *testing.T) {
	// SUT + act
	var result = customizationDefault.ServerCert()

	// assert
	assert.Nil(t, result)
}

func TestDefaultCustomization_CaCertPool(t *testing.T) {
	// SUT + act
	var result = customizationDefault.CaCertPool()

	// assert
	assert.Nil(t, result)
}

func TestDefaultCustomization_GraceShutdownWaitTime(t *testing.T) {
	// SUT + act
	var result = customizationDefault.GraceShutdownWaitTime()

	// assert
	assert.Equal(t, 3*time.Minute, result)
}

func TestDefaultCustomization_Routes(t *testing.T) {
	// SUT + act
	var results = customizationDefault.Routes()

	// assert
	assert.Empty(t, results)
}

func TestDefaultCustomization_Statics(t *testing.T) {
	// SUT + act
	var results = customizationDefault.Statics()

	// assert
	assert.Empty(t, results)
}

func TestDefaultCustomization_Middlewares(t *testing.T) {
	// SUT + act
	var results = customizationDefault.Middlewares()

	// assert
	assert.Empty(t, results)
}

func TestDefaultCustomization_InstrumentRouter(t *testing.T) {
	// arrange
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}

	// SUT + act
	var result = customizationDefault.InstrumentRouter(dummyRouter)

	// assert
	assert.Equal(t, dummyRouter, result)
}

func TestDefaultCustomization_WrapHandler(t *testing.T) {
	// arrange
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}

	// SUT + act
	var result = customizationDefault.WrapHandler(dummyRouter)

	// assert
	assert.Equal(t, dummyRouter, result)
}

func TestDefaultCustomization_Listener(t *testing.T) {
	// SUT + act
	var result = customizationDefault.Listener()

	// assert
	assert.Nil(t, result)
}

func TestDefaultCustomization_PreAction(t *testing.T) {
	// arrange
	var dummySession Session

	// SUT + act
	var err = customizationDefault.PreAction(dummySession)

	// assert
	assert.NoError(t, err)
}

func TestDefaultCustomization_PostAction(t *testing.T) {
	// arrange
	var dummySession Session

	// SUT + act
	var err = customizationDefault.PostAction(dummySession)

	// assert
	assert.NoError(t, err)
}

func TestDefaultCustomization_InterpretSuccess_NilResponseContent(t *testing.T) {
	// arrange
	var dummyResponseContent interface{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummyResponseContent, i)
		return true
	})

	// SUT + act
	var code, message = customizationDefault.InterpretSuccess(
		dummyResponseContent,
	)

	// assert
	assert.Equal(t, http.StatusNoContent, code)
	assert.Zero(t, message)
}

func TestDefaultCustomization_InterpretSuccess_EmptyResponseContent(t *testing.T) {
	// arrange
	var dummyResponseContent interface{}
	var dummyContent string

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummyResponseContent, i)
		return false
	})
	m.ExpectFunc(marshalIgnoreError, 1, func(v interface{}) string {
		assert.Equal(t, dummyResponseContent, v)
		return dummyContent
	})

	// SUT + act
	var code, message = customizationDefault.InterpretSuccess(
		dummyResponseContent,
	)

	// assert
	assert.Equal(t, http.StatusNoContent, code)
	assert.Zero(t, message)
}

func TestDefaultCustomization_InterpretSuccess_HappyPath(t *testing.T) {
	// arrange
	var dummyResponseContent interface{}
	var dummyContent = "some content"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummyResponseContent, i)
		return false
	})
	m.ExpectFunc(marshalIgnoreError, 1, func(v interface{}) string {
		assert.Equal(t, dummyResponseContent, v)
		return dummyContent
	})

	// SUT + act
	var code, message = customizationDefault.InterpretSuccess(
		dummyResponseContent,
	)

	// assert
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, dummyContent, message)
}

func TestDefaultCustomization_InterpretError_NormalError(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyError = errors.New(dummyErrorMessage)

	// SUT + act
	var code, message = customizationDefault.InterpretError(
		dummyError,
	)

	// assert
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, dummyErrorMessage, message)
}

func TestDefaultCustomization_InterpretError_AppError(t *testing.T) {
	// arrange
	var dummyStatusCode = rand.Intn(600)
	var dummyResponseMessage = "some response message"
	var dummyAppError = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyAppError, "HTTPStatusCode", 1, func() int {
		return dummyStatusCode
	})
	m.ExpectMethod(dummyAppError, "HTTPResponseMessage", 1, func() string {
		return dummyResponseMessage
	})

	// SUT + act
	var code, message = customizationDefault.InterpretError(
		dummyAppError,
	)

	// assert
	assert.Equal(t, dummyStatusCode, code)
	assert.Equal(t, dummyResponseMessage, message)
}

func TestGetRecoverError_Error(t *testing.T) {
	// arrange
	var dummyRecoverResult = errors.New("some error")

	// SUT + act
	var result = getRecoverError(
		dummyRecoverResult,
	)

	// assert
	assert.Equal(t, dummyRecoverResult, result)
}

func TestGetRecoverError_NonError(t *testing.T) {
	// arrange
	var dummyRecoverResult = "some recovery result"
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(fmt.Errorf, 1, func(format string, a ...interface{}) error {
		assert.Equal(t, "endpoint panic: %v", format)
		assert.Equal(t, 1, len(a))
		assert.Equal(t, dummyRecoverResult, a[0])
		return dummyError
	})

	// SUT + act
	var result = getRecoverError(
		dummyRecoverResult,
	)

	// assert
	assert.Equal(t, dummyError, result)
}

func TestDefaultCustomization_RecoverPanic(t *testing.T) {
	// arrange
	var dummySession = &session{}
	var dummyError = errors.New("some error")
	var dummyRecoverResult = dummyError.(interface{})
	var dummyDebugStackString = "some debug stack string"
	var dummyDebugStack = []byte(dummyDebugStackString)
	var dummyName = "some name"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(getRecoverError, 1, func(recoverResult interface{}) error {
		assert.Equal(t, dummyRecoverResult, recoverResult)
		return dummyError
	})
	m.ExpectFunc(debug.Stack, 1, func() []byte {
		return dummyDebugStack
	})
	m.ExpectMethod(dummySession, "GetName", 1, func(self *session) string {
		assert.Equal(t, dummySession, self)
		return dummyName
	})
	m.ExpectMethod(dummySession, "LogMethodLogic", 1, func(self *session, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, self)
		assert.Equal(t, LogLevelError, logLevel)
		assert.Equal(t, "RecoverPanic", category)
		assert.Equal(t, dummyName, subcategory)
		assert.Equal(t, "Error: %+v\nCallstack: %v", messageFormat)
		assert.Equal(t, 2, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
		assert.Equal(t, dummyDebugStackString, parameters[1])
	})

	// SUT + act
	var result, err = customizationDefault.RecoverPanic(
		dummySession,
		dummyRecoverResult,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyError, err)
}

func TestDefaultCustomization_NotFoundHandler(t *testing.T) {
	// SUT + act
	var result = customizationDefault.NotFoundHandler()

	// assert
	assert.Nil(t, result)
}

func TestDefaultCustomization_MethodNotAllowedHandler(t *testing.T) {
	// SUT + act
	var result = customizationDefault.MethodNotAllowedHandler()

	// assert
	assert.Nil(t, result)
}

func TestDefaultCustomization_ClientCert(t *testing.T) {
	// SUT + act
	var result = customizationDefault.ClientCert()

	// assert
	assert.Nil(t, result)
}

func TestDefaultCustomization_DefaultTimeout(t *testing.T) {
	// SUT + act
	var result = customizationDefault.DefaultTimeout()

	// assert
	assert.Equal(t, 3*time.Minute, result)
}

func TestDefaultCustomization_SkipServerCertVerification(t *testing.T) {
	// SUT + act
	var result = customizationDefault.SkipServerCertVerification()

	// assert
	assert.False(t, result)
}

func TestDefaultCustomization_RoundTripper(t *testing.T) {
	// arrange
	var dummyTransport = &http.Transport{}

	// SUT + act
	var result = customizationDefault.RoundTripper(dummyTransport)

	// assert
	assert.Equal(t, dummyTransport, result)
}

func TestDefaultCustomization_WrapRequest(t *testing.T) {
	// arrange
	var dummySession Session
	var dummyRequest = &http.Request{Host: "some host"}

	// SUT + act
	var result = customizationDefault.WrapRequest(dummySession, dummyRequest)

	// assert
	assert.Equal(t, dummyRequest, result)
}
