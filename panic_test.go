package webserver

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

func TestGetDebugStack(t *testing.T) {
	// SUT + act
	var result = getDebugStack()

	// assert
	assert.NotZero(t, result)

	// verify
	verifyAll(t)
}

func TestHandlePanic_NilRecoverResult(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyRecoverResult interface{}

	// mock
	createMock(t)

	// SUT + act
	handlePanic(
		dummySession,
		dummyRecoverResult,
	)

	// verify
	verifyAll(t)
}

func TestHandlePanic_HappyPath(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyError = errors.New("some error")
	var dummyRecoverResult = dummyError.(interface{})
	var dummyDebugStack = "some debug stack"

	// mock
	createMock(t)

	// expect
	getRecoverErrorFuncExpected = 1
	getRecoverErrorFunc = func(recoverResult interface{}) error {
		getRecoverErrorFuncCalled++
		assert.Equal(t, dummyRecoverResult, recoverResult)
		return dummyError
	}
	writeResponseFuncExpected = 1
	writeResponseFunc = func(session *session, responseObject interface{}, responseError error) {
		writeResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Nil(t, responseObject)
		assert.Equal(t, dummyError, responseError)
	}
	getDebugStackFuncExpected = 1
	getDebugStackFunc = func() string {
		getDebugStackFuncCalled++
		return dummyDebugStack
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "panic", category)
		assert.Equal(t, "Handle", subcategory)
		assert.Equal(t, "%+v\n%v", messageFormat)
		assert.Equal(t, 2, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
		assert.Equal(t, dummyDebugStack, parameters[1])
	}

	// SUT + act
	handlePanic(
		dummySession,
		dummyRecoverResult,
	)

	// verify
	verifyAll(t)
}
