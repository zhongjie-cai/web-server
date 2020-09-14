package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCustomization_PreBootstrap(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var err = customizationDefault.PreBootstrap()

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_PostBootstrap(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var err = customizationDefault.PostBootstrap()

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_AppClosing(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var err = customizationDefault.AppClosing()

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

type dummySessionLog struct {
	dummySession
	getID   func() uuid.UUID
	getName func() string
}

func (session *dummySessionLog) GetID() uuid.UUID {
	if session.getID != nil {
		return session.getID()
	}
	assert.Fail(session.t, "Unexpected call to GetID")
	return uuid.Nil
}

func (session *dummySessionLog) GetName() string {
	if session.getName != nil {
		return session.getName()
	}
	assert.Fail(session.t, "Unexpected call to GetName")
	return ""
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
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummySession, i)
		return true
	}

	// SUT + act
	customizationDefault.Log(
		dummySession,
		dummyLogType,
		dummyLogLevel,
		dummyCategory,
		dummySubcategory,
		dummyDescription,
	)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_Log_HappyPath(t *testing.T) {
	// arrange
	var dummySession = &dummySessionLog{dummySession: dummySession{t: t}}
	var dummyLogType = LogType(rand.Intn(100))
	var dummyLogLevel = LogLevel(rand.Intn(100))
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyDescription = "some description"
	var sessionGetIDExpected int
	var sessionGetIDCalled int
	var sessionGetNameExpected int
	var sessionGetNameCalled int
	var dummySessionID = uuid.New()
	var dummySessionName = "some session name"
	var dummyFormat = "<%v|%v> (%v|%v) [%v|%v] %v\n"

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummySession, i)
		return false
	}
	sessionGetIDExpected = 1
	dummySession.getID = func() uuid.UUID {
		sessionGetIDCalled++
		return dummySessionID
	}
	sessionGetNameExpected = 1
	dummySession.getName = func() string {
		sessionGetNameCalled++
		return dummySessionName
	}
	fmtPrintfExpected = 1
	fmtPrintf = func(format string, a ...interface{}) (n int, err error) {
		fmtPrintfCalled++
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
	}

	// SUT + act
	customizationDefault.Log(
		dummySession,
		dummyLogType,
		dummyLogLevel,
		dummyCategory,
		dummySubcategory,
		dummyDescription,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, sessionGetIDExpected, sessionGetIDCalled, "Unexpected number of calls to method sessionGetID")
	assert.Equal(t, sessionGetNameExpected, sessionGetNameCalled, "Unexpected number of calls to method sessionGetName")
}

func TestDefaultCustomization_ServerCert(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.ServerCert()

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_CaCertPool(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.CaCertPool()

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_GraceShutdownWaitTime(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.GraceShutdownWaitTime()

	// assert
	assert.Equal(t, 3*time.Minute, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_Routes(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var results = customizationDefault.Routes()

	// assert
	assert.Empty(t, results)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_Statics(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var results = customizationDefault.Statics()

	// assert
	assert.Empty(t, results)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_Middlewares(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var results = customizationDefault.Middlewares()

	// assert
	assert.Empty(t, results)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_InstrumentRouter(t *testing.T) {
	// arrange
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}

	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.InstrumentRouter(dummyRouter)

	// assert
	assert.Equal(t, dummyRouter, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_PreAction(t *testing.T) {
	// arrange
	var dummySession Session

	// mock
	createMock(t)

	// SUT + act
	var err = customizationDefault.PreAction(dummySession)

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_PostAction(t *testing.T) {
	// arrange
	var dummySession Session

	// mock
	createMock(t)

	// SUT + act
	var err = customizationDefault.PostAction(dummySession)

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_InterpretSuccess_NilResponseContent(t *testing.T) {
	// arrange
	var dummyResponseContent interface{}

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyResponseContent, i)
		return true
	}

	// SUT + act
	var code, message = customizationDefault.InterpretSuccess(
		dummyResponseContent,
	)

	// assert
	assert.Equal(t, http.StatusNoContent, code)
	assert.Zero(t, message)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_InterpretSuccess_EmptyResponseContent(t *testing.T) {
	// arrange
	var dummyResponseContent interface{}
	var dummyContent string

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyResponseContent, i)
		return false
	}
	marshalIgnoreErrorFuncExpected = 1
	marshalIgnoreErrorFunc = func(v interface{}) string {
		marshalIgnoreErrorFuncCalled++
		assert.Equal(t, dummyResponseContent, v)
		return dummyContent
	}

	// SUT + act
	var code, message = customizationDefault.InterpretSuccess(
		dummyResponseContent,
	)

	// assert
	assert.Equal(t, http.StatusNoContent, code)
	assert.Zero(t, message)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_InterpretSuccess_HappyPath(t *testing.T) {
	// arrange
	var dummyResponseContent interface{}
	var dummyContent = "some content"

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyResponseContent, i)
		return false
	}
	marshalIgnoreErrorFuncExpected = 1
	marshalIgnoreErrorFunc = func(v interface{}) string {
		marshalIgnoreErrorFuncCalled++
		assert.Equal(t, dummyResponseContent, v)
		return dummyContent
	}

	// SUT + act
	var code, message = customizationDefault.InterpretSuccess(
		dummyResponseContent,
	)

	// assert
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, dummyContent, message)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_InterpretError_NormalError(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyError = errors.New(dummyErrorMessage)

	// mock
	createMock(t)

	// SUT + act
	var code, message = customizationDefault.InterpretError(
		dummyError,
	)

	// assert
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, dummyErrorMessage, message)

	// verify
	verifyAll(t)
}

type dummyAppErrorInterpretError struct {
	dummyAppError
	httpStatusCode      func() int
	httpResponseMessage func() string
}

func (dummyAppErrorInterpretError *dummyAppErrorInterpretError) HTTPStatusCode() int {
	if dummyAppErrorInterpretError.httpStatusCode != nil {
		return dummyAppErrorInterpretError.httpStatusCode()
	}
	assert.Fail(dummyAppErrorInterpretError.t, "Unexpected call to method AppError.HTTPStatusCode")
	return 0
}

func (dummyAppErrorInterpretError *dummyAppErrorInterpretError) HTTPResponseMessage() string {
	if dummyAppErrorInterpretError.httpResponseMessage != nil {
		return dummyAppErrorInterpretError.httpResponseMessage()
	}
	assert.Fail(dummyAppErrorInterpretError.t, "Unexpected call to method AppError.HTTPResponseMessage")
	return ""
}

func TestDefaultCustomization_InterpretError_AppError(t *testing.T) {
	// arrange
	var dummyStatusCode = rand.Intn(600)
	var dummyResponseMessage = "some response message"
	var dummyAppErrorInterpretError = &dummyAppErrorInterpretError{
		dummyAppError: dummyAppError{t: t},
	}
	var appErrorHTTPStatusCodeExpected int
	var appErrorHTTPStatusCodeCalled int
	var appErrorHTTPResponseMessageExpected int
	var appErrorHTTPResponseMessageCalled int

	// mock
	createMock(t)

	// expect
	appErrorHTTPStatusCodeExpected = 1
	dummyAppErrorInterpretError.httpStatusCode = func() int {
		appErrorHTTPStatusCodeCalled++
		return dummyStatusCode
	}
	appErrorHTTPResponseMessageExpected = 1
	dummyAppErrorInterpretError.httpResponseMessage = func() string {
		appErrorHTTPResponseMessageCalled++
		return dummyResponseMessage
	}

	// SUT + act
	var code, message = customizationDefault.InterpretError(
		dummyAppErrorInterpretError,
	)

	// assert
	assert.Equal(t, dummyStatusCode, code)
	assert.Equal(t, dummyResponseMessage, message)

	// verify
	verifyAll(t)
	assert.Equal(t, appErrorHTTPStatusCodeExpected, appErrorHTTPStatusCodeCalled, "Unexpected number of calls to appError.HTTPStatusCode")
	assert.Equal(t, appErrorHTTPResponseMessageExpected, appErrorHTTPResponseMessageCalled, "Unexpected number of calls to appError.HTTPResponseMessage")
}

func TestGetRecoverError_Error(t *testing.T) {
	// arrange
	var dummyRecoverResult = errors.New("some error")

	// mock
	createMock(t)

	// SUT + act
	var result = getRecoverError(
		dummyRecoverResult,
	)

	// assert
	assert.Equal(t, dummyRecoverResult, result)

	// verify
	verifyAll(t)
}

func TestGetRecoverError_NonError(t *testing.T) {
	// arrange
	var dummyRecoverResult = "some recovery result"
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	fmtErrorfExpected = 1
	fmtErrorf = func(format string, a ...interface{}) error {
		fmtErrorfCalled++
		assert.Equal(t, "Endpoint panic: %v", format)
		assert.Equal(t, 1, len(a))
		assert.Equal(t, dummyRecoverResult, a[0])
		return dummyError
	}

	// SUT + act
	var result = getRecoverError(
		dummyRecoverResult,
	)

	// assert
	assert.Equal(t, dummyError, result)

	// verify
	verifyAll(t)
}

type dummySessionRecoverPanic struct {
	dummySession
	getName        func() string
	logMethodLogic func(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{})
}

func (dummySessionRecoverPanic *dummySessionRecoverPanic) GetName() string {
	if dummySessionRecoverPanic.getName != nil {
		return dummySessionRecoverPanic.getName()
	}
	assert.Fail(dummySessionRecoverPanic.t, "Unexpected call to GetName")
	return ""
}

func (dummySessionRecoverPanic *dummySessionRecoverPanic) LogMethodLogic(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	if dummySessionRecoverPanic.logMethodLogic != nil {
		dummySessionRecoverPanic.logMethodLogic(logLevel, category, subcategory, messageFormat, parameters...)
		return
	}
	assert.Fail(dummySessionRecoverPanic.t, "Unexpected call to LogMethodLogic")
}

func TestDefaultCustomization_RecoverPanic(t *testing.T) {
	// arrange
	var dummySessionRecoverPanic = &dummySessionRecoverPanic{
		dummySession: dummySession{t: t},
	}
	var dummyError = errors.New("some error")
	var dummyRecoverResult = dummyError.(interface{})
	var dummyDebugStackString = "some debug stack string"
	var dummyDebugStack = []byte(dummyDebugStackString)
	var dummyName = "some name"
	var sessionGetNameExpected int
	var sessionGetNameCalled int
	var sessionLogMethodLogicExpected int
	var sessionLogMethodLogicCalled int

	// mock
	createMock(t)

	// expect
	getRecoverErrorFuncExpected = 1
	getRecoverErrorFunc = func(recoverResult interface{}) error {
		getRecoverErrorFuncCalled++
		assert.Equal(t, dummyRecoverResult, recoverResult)
		return dummyError
	}
	debugStackExpected = 1
	debugStack = func() []byte {
		debugStackCalled++
		return dummyDebugStack
	}
	sessionGetNameExpected = 1
	dummySessionRecoverPanic.getName = func() string {
		sessionGetNameCalled++
		return dummyName
	}
	sessionLogMethodLogicExpected = 1
	dummySessionRecoverPanic.logMethodLogic = func(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		sessionLogMethodLogicCalled++
		assert.Equal(t, LogLevelError, logLevel)
		assert.Equal(t, "RecoverPanic", category)
		assert.Equal(t, dummyName, subcategory)
		assert.Equal(t, "Error: %+v\nCallstack: %v", messageFormat)
		assert.Equal(t, 2, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
		assert.Equal(t, dummyDebugStackString, parameters[1])
	}

	// SUT + act
	var result, err = customizationDefault.RecoverPanic(
		dummySessionRecoverPanic,
		dummyRecoverResult,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyError, err)

	// verify
	verifyAll(t)
	assert.Equal(t, sessionGetNameExpected, sessionGetNameCalled, "Unexpected number of calls to method session.GetName")
	assert.Equal(t, sessionLogMethodLogicExpected, sessionLogMethodLogicCalled, "Unexpected number of calls to method session.LogMethodLogic")
}

func TestDefaultCustomization_NotFoundHandler(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.NotFoundHandler()

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_MethodNotAllowedHandler(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.MethodNotAllowedHandler()

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_ClientCert(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.ClientCert()

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_DefaultTimeout(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.DefaultTimeout()

	// assert
	assert.Equal(t, 3*time.Minute, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_SkipServerCertVerification(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.SkipServerCertVerification()

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_RoundTripper(t *testing.T) {
	// arrange
	var dummyTransport = &dummyTransport{t: t}

	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.RoundTripper(dummyTransport)

	// assert
	assert.Equal(t, dummyTransport, result)

	// verify
	verifyAll(t)
}

func TestDefaultCustomization_WrapRequest(t *testing.T) {
	// arrange
	var dummySession Session
	var dummyRequest = &http.Request{Host: "some host"}

	// mock
	createMock(t)

	// SUT + act
	var result = customizationDefault.WrapRequest(dummySession, dummyRequest)

	// assert
	assert.Equal(t, dummyRequest, result)

	// verify
	verifyAll(t)
}
