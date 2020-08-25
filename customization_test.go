package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/google/uuid"

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

func TestInterpretSuccess_NilResponseContent(t *testing.T) {
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

func TestInterpretSuccess_EmptyResponseContent(t *testing.T) {
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

func TestInterpretSuccess_HappyPath(t *testing.T) {
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

func TestInterpretError(t *testing.T) {
	// arrange
	var testData = map[error]int{
		errSessionNil:            http.StatusInternalServerError,
		errRouteRegistration:     http.StatusInternalServerError,
		errRouteNotFound:         http.StatusNotFound,
		errHostServer:            http.StatusInternalServerError,
		ErrRequestBodyEmpty:      http.StatusBadRequest,
		ErrRequestBodyInvalid:    http.StatusBadRequest,
		ErrParameterNotFound:     http.StatusBadRequest,
		ErrParameterInvalid:      http.StatusBadRequest,
		ErrQueryNotFound:         http.StatusBadRequest,
		ErrQueryInvalid:          http.StatusBadRequest,
		ErrHeaderNotFound:        http.StatusBadRequest,
		ErrHeaderInvalid:         http.StatusBadRequest,
		ErrWebRequestNil:         http.StatusInternalServerError,
		ErrResponseInvalid:       http.StatusInternalServerError,
		errors.New("some error"): http.StatusInternalServerError,
	}
	var dummyMessage = "some message"

	for dummyError, dummyCode := range testData {
		// mock
		createMock(t)

		// expect
		fmtSprintfExpected = 1
		fmtSprintf = func(format string, a ...interface{}) string {
			fmtSprintfCalled++
			assert.Equal(t, "%+v", format)
			assert.Equal(t, 1, len(a))
			assert.Equal(t, dummyError, a[0])
			return dummyMessage
		}

		// SUT + act
		var code, message = customizationDefault.InterpretError(
			dummyError,
		)

		// assert
		assert.Equal(t, dummyCode, code)
		assert.Equal(t, dummyMessage, message)

		// verify
		verifyAll(t)
	}
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
