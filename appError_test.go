package webserver

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppError(t *testing.T) {
	// arrange
	var dummyErrorCode = errorCode("some error code")
	var dummyErrorMessage = "some error message"
	var dummyInnerErrors = []error{
		errors.New("some inner error message 1"),
		errors.New("some inner error message 2"),
		errors.New("some inner error message 3"),
	}
	var dummyCleanedUpErrors = []*appError{
		&appError{Message: "some inner error message 1"},
		&appError{Message: "some inner error message 2"},
		&appError{Message: "some inner error message 3"},
	}

	// mock
	createMock(t)

	// expect
	cleanupInnerErrorsFuncExpected = 1
	cleanupInnerErrorsFunc = func(innerErrors []error) []*appError {
		cleanupInnerErrorsFuncCalled++
		assert.Equal(t, dummyInnerErrors, innerErrors)
		return dummyCleanedUpErrors
	}

	// SUT + act
	var result = newAppError(
		dummyErrorCode,
		dummyErrorMessage,
		dummyInnerErrors,
	)

	// assert
	assert.Equal(t, dummyErrorCode, result.Code)
	assert.Equal(t, dummyErrorMessage, result.Message)
	assert.Equal(t, dummyCleanedUpErrors, result.InnerErrors)

	// verify
	verifyAll(t)
}

func TestGetErrorMessage(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"

	// mock
	createMock(t)

	// SUT
	var sut = errors.New(dummyErrorMessage)

	// act
	var result = getErrorMessage(sut)

	// assert
	assert.Equal(t, dummyErrorMessage, result)

	// verify
	verifyAll(t)
}

func TestPrintInnerErrors_NilInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors []*appError

	// mock
	createMock(t)

	// SUT + act
	var result = printInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestPrintInnerErrors_EmptyInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors = []*appError{}

	// mock
	createMock(t)

	// SUT + act
	var result = printInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestPrintInnerErrors_ValidInnerErrors(t *testing.T) {
	// arrange
	var dummyErrorMessage1 = "some inner error 1"
	var dummyErrorMessage2 = "some inner error 2"
	var dummyErrorMessage3 = "some inner error 3"
	var dummyErrorMessages = []string{
		dummyErrorMessage1,
		dummyErrorMessage2,
		dummyErrorMessage3,
	}
	var dummyInnerError1 = &appError{Message: dummyErrorMessage1}
	var dummyInnerError2 = &appError{Message: dummyErrorMessage2}
	var dummyInnerError3 = &appError{Message: dummyErrorMessage3}
	var dummyInnerErrors = []*appError{
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	}
	var dummyJoinedMessage = "some joined message"
	var dummyResult = "some result"

	// mock
	createMock(t)

	// expect
	getErrorMessageFuncExpected = len(dummyInnerErrors)
	getErrorMessageFunc = func(err error) string {
		getErrorMessageFuncCalled++
		assert.Equal(t, dummyInnerErrors[getErrorMessageFuncCalled-1], err)
		return dummyErrorMessages[getErrorMessageFuncCalled-1]
	}
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		assert.ElementsMatch(t, dummyErrorMessages, a)
		assert.Equal(t, errorSeparator, sep)
		return dummyJoinedMessage
	}
	fmtSprintfExpected = 1
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		assert.Equal(t, errorJoiningFormat, format)
		assert.Equal(t, 1, len(a))
		assert.Equal(t, dummyJoinedMessage, a[0])
		return dummyResult
	}

	// SUT + act
	var result = printInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestAppError_Error(t *testing.T) {
	// arrange
	var dummyCode = errorCode("some error code")
	var dummyMessage = "some error message"
	var dummyInnerErrors = []*appError{
		&appError{Message: "some inner error 1"},
		&appError{Message: "some inner error 2"},
		&appError{Message: "some inner error 3"},
	}
	var dummyAppError = &appError{
		dummyCode,
		dummyMessage,
		dummyInnerErrors,
	}
	var dummyBaseErrorMessage = "some base error message"
	var dummyInnerErrorMessage = "some inner error message"
	var dummyResult = "some result"

	// mock
	createMock(t)

	// expect
	fmtSprintfExpected = 1
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		assert.Equal(t, errorMessageorFormat, format)
		assert.Equal(t, 2, len(a))
		assert.Equal(t, dummyCode, a[0])
		assert.Equal(t, dummyMessage, a[1])
		return dummyBaseErrorMessage
	}
	printInnerErrorsFuncExpected = 1
	printInnerErrorsFunc = func(innerErrors []*appError) string {
		printInnerErrorsFuncCalled++
		assert.Equal(t, dummyInnerErrors, innerErrors)
		return dummyInnerErrorMessage
	}
	fmtSprintExpected = 1
	fmtSprint = func(a ...interface{}) string {
		fmtSprintCalled++
		assert.Equal(t, 2, len(a))
		assert.Equal(t, dummyBaseErrorMessage, a[0])
		assert.Equal(t, dummyInnerErrorMessage, a[1])
		return dummyResult
	}

	// SUT + act
	var result = dummyAppError.Error()

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestAppError_ErrorCode(t *testing.T) {
	// arrange
	var dummyCode = errorCode("some error code")
	var dummyMessage = "some error message"

	// mock
	createMock(t)

	// SUT
	var appError = &appError{
		Code:    dummyCode,
		Message: dummyMessage,
	}

	// act
	var code = appError.ErrorCode()

	// assert
	assert.Equal(t, string(dummyCode), code)

	// verify
	verifyAll(t)
}

func TestAppError_HTTPStatusCode(t *testing.T) {
	// arrange
	var dummyCode = errorCode("some error code")
	var dummyMessage = "some error message"

	// mock
	createMock(t)

	// SUT
	var appError = &appError{
		Code:    dummyCode,
		Message: dummyMessage,
	}

	// act
	var code = appError.HTTPStatusCode()

	// assert
	assert.Equal(t, dummyCode.httpStatusCode(), code)

	// verify
	verifyAll(t)
}

func TestAppError_HTTPResponseMessage(t *testing.T) {
	// arrange
	var dummyCode = errorCode("some error code")
	var dummyMessage = "some error message"
	var dummyInnerErrors = []*appError{
		&appError{Message: "some inner error 1"},
		&appError{Message: "some inner error 2"},
		&appError{Message: "some inner error 3"},
	}
	var dummyAppError = &appError{
		dummyCode,
		dummyMessage,
		dummyInnerErrors,
	}
	var dummyResult = "some result"
	var dummyBytes = []byte(dummyResult)
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	jsonMarshalExpected = 1
	jsonMarshal = func(v interface{}) ([]byte, error) {
		jsonMarshalCalled++
		assert.Equal(t, dummyAppError, v)
		return dummyBytes, dummyError
	}

	// SUT
	var sut = dummyAppError

	// act
	var result = sut.HTTPResponseMessage()

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestEqualsError_SameError(t *testing.T) {
	// arrange
	var dummyError = errors.New("some error")
	var dummyTarget = dummyError

	// mock
	createMock(t)

	// expect

	// SUT + act
	var result = equalsError(
		dummyError,
		dummyTarget,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestEqualsError_SameMessage(t *testing.T) {
	// arrange
	var dummyError = errors.New("some error")
	var dummyTarget = errors.New("some error")

	// mock
	createMock(t)

	// expect

	// SUT + act
	var result = equalsError(
		dummyError,
		dummyTarget,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestEqualsError_ErrorIs(t *testing.T) {
	// arrange
	var dummyError = errors.New("some error")
	var dummyTarget = errors.New("some target")
	var dummyResult = rand.Intn(100) > 50

	// mock
	createMock(t)

	// expect
	errorsIsExpected = 1
	errorsIs = func(err, target error) bool {
		errorsIsCalled++
		return dummyResult
	}

	// SUT + act
	var result = equalsError(
		dummyError,
		dummyTarget,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestAppErrorContains_HappyPath(t *testing.T) {
	// arrange
	var dummyError = errors.New("some error")
	var dummyappError = &appError{
		Code:    errorCode("some error code"),
		Message: dummyError.Error(),
	}

	// mock
	createMock(t)

	// expect
	equalsErrorFuncExpected = 1
	equalsErrorFunc = func(err, target error) bool {
		equalsErrorFuncCalled++
		assert.Equal(t, dummyappError, err)
		assert.Equal(t, dummyError, target)
		return true
	}

	// SUT + act
	var result = appErrorContains(
		dummyappError,
		dummyError,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestInnerErrorContains_NilInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors []*appError
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// SUT + act
	var result = innerErrorContains(
		dummyInnerErrors,
		dummyError,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}

func TestInnerErrorContains_EmptyInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors = []*appError{}
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// SUT + act
	var result = innerErrorContains(
		dummyInnerErrors,
		dummyError,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}

func TestInnerErrorContains_ValidInnerError(t *testing.T) {
	// arrange
	var dummyInnerError = &appError{
		Code:    errorCode("some error code"),
		Message: "some error message",
	}
	var dummyInnerErrors = []*appError{
		dummyInnerError,
	}
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	appErrorContainsFuncExpected = 1
	appErrorContainsFunc = func(appError AppError, err error) bool {
		appErrorContainsFuncCalled++
		assert.Equal(t, dummyInnerError, appError)
		assert.Equal(t, dummyError, err)
		return true
	}

	// SUT + act
	var result = innerErrorContains(
		dummyInnerErrors,
		dummyError,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestInnerErrorContains_NoMatchingErrors(t *testing.T) {
	// arrange
	var dummyInnerError1 = &appError{
		Code:    errorCode("some error code 1"),
		Message: "some error message 1",
	}
	var dummyInnerError2 = &appError{
		Code:    errorCode("some error code 2"),
		Message: "some error message 2",
	}
	var dummyInnerErrors = []*appError{
		dummyInnerError1,
		dummyInnerError2,
	}
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	appErrorContainsFuncExpected = 2
	appErrorContainsFunc = func(appError AppError, err error) bool {
		appErrorContainsFuncCalled++
		assert.Equal(t, dummyError, err)
		if appErrorContainsFuncCalled == 1 {
			assert.Equal(t, dummyInnerError1, appError)
		} else if appErrorContainsFuncCalled == 2 {
			assert.Equal(t, dummyInnerError2, appError)
		}
		return false
	}

	// SUT + act
	var result = innerErrorContains(
		dummyInnerErrors,
		dummyError,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}

func TestAppError_Contains_ErrorEqual(t *testing.T) {
	// arrange
	var dummyappError = &appError{
		Code:    errorCode("some error code"),
		Message: "some error message",
	}
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	equalsErrorFuncExpected = 1
	equalsErrorFunc = func(err, target error) bool {
		equalsErrorFuncCalled++
		assert.Equal(t, dummyappError, err)
		assert.Equal(t, dummyError, target)
		return true
	}

	// SUT
	var sut = dummyappError

	// act
	var result = sut.Contains(
		dummyError,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestAppError_Contains_InnerErrorEqual(t *testing.T) {
	// arrange
	var dummyInnerErrors = []*appError{
		&appError{Message: "some inner error 1"},
		&appError{Message: "some inner error 2"},
		&appError{Message: "some inner error 3"},
	}
	var dummyappError = &appError{
		Code:        errorCode("some error code"),
		Message:     "some error message",
		InnerErrors: dummyInnerErrors,
	}
	var dummyError = errors.New("some error")
	var dummyResult = rand.Intn(100) > 50

	// mock
	createMock(t)

	// expect
	equalsErrorFuncExpected = 1
	equalsErrorFunc = func(err, target error) bool {
		equalsErrorFuncCalled++
		assert.Equal(t, dummyappError, err)
		assert.Equal(t, dummyError, target)
		return false
	}
	innerErrorContainsFuncExpected = 1
	innerErrorContainsFunc = func(innerErrors []*appError, err error) bool {
		innerErrorContainsFuncCalled++
		assert.Equal(t, dummyInnerErrors, innerErrors)
		return dummyResult
	}

	// SUT
	var sut = dummyappError

	// act
	var result = sut.Contains(
		dummyError,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestCleanupInnerErrors_NilInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors []error

	// mock
	createMock(t)

	// SUT + act
	var result = cleanupInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Empty(t, result)

	// verify
	verifyAll(t)
}

func TestCleanupInnerErrors_EmptyInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors = []error{}

	// mock
	createMock(t)

	// SUT + act
	var result = cleanupInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Empty(t, result)

	// verify
	verifyAll(t)
}

func TestCleanupInnerErrors_NoValidInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors = []error{
		nil,
		nil,
		nil,
	}

	// mock
	createMock(t)

	// SUT + act
	var result = cleanupInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Empty(t, result)

	// verify
	verifyAll(t)
}

func TestCleanupInnerErrors_HasValidInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerError1 = errors.New("some random error 1")
	var dummyInnerError2 = errors.New("some random error 2")
	var dummyInnerError3 = errors.New("some random error 3")
	var dummyInnerErrors = []error{
		dummyInnerError1,
		nil,
		dummyInnerError2,
		nil,
		dummyInnerError3,
	}

	// mock
	createMock(t)

	// SUT + act
	var result = cleanupInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Equal(t, 3, len(result))
	assert.Equal(t, dummyInnerError1.Error(), result[0].Message)
	assert.Equal(t, dummyInnerError2.Error(), result[1].Message)
	assert.Equal(t, dummyInnerError3.Error(), result[2].Message)

	// verify
	verifyAll(t)
}

func TestAppErrorWrap_NoInnerError(t *testing.T) {
	// arrange
	var dummyMessage = "dummy error"
	var dummyCode = errorCodeGeneralFailure
	var dummyInnerErrorMessage = "dummy inner error"
	var dummyInnerMostErrorMessage = "dummy inner most error"
	var dummyInnerError1 = &appError{
		Message: "dummy inner error 1",
	}
	var dummyInnerError2 = &appError{
		Code:    errorCodeGeneralFailure,
		Message: dummyInnerErrorMessage,
		InnerErrors: []*appError{
			&appError{Message: dummyInnerMostErrorMessage},
		},
	}
	var dummyInnerError3 = &appError{
		Message: "dummy inner error 3",
	}
	var dummyInnerErrors = []*appError{
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	}
	var dummyCleanedInnerErrors = []*appError{}

	// mock
	createMock(t)

	// expect
	cleanupInnerErrorsFuncExpected = 1
	cleanupInnerErrorsFunc = func(innerErrors []error) []*appError {
		cleanupInnerErrorsFuncCalled++
		assert.Equal(t, 3, len(innerErrors))
		assert.NoError(t, innerErrors[0])
		assert.NoError(t, innerErrors[1])
		assert.NoError(t, innerErrors[2])
		return dummyCleanedInnerErrors
	}

	// SUT
	var appError = &appError{
		Code:        dummyCode,
		Message:     dummyMessage,
		InnerErrors: dummyInnerErrors,
	}

	// act
	appError.Wrap(
		nil,
		nil,
		nil,
	)

	// assert
	assert.Equal(t, dummyInnerErrors, appError.InnerErrors)

	// verify
	verifyAll(t)
}

func TestAppErrorWrap_HasInnerError(t *testing.T) {
	// arrange
	var dummyMessage = "dummy error"
	var dummyCode = errorCodeGeneralFailure
	var dummyInnerErrorMessage = "dummy inner error"
	var dummyInnerMostErrorMessage = "dummy inner most error"
	var dummyInnerError1 = &appError{
		Message: "dummy inner error 1",
	}
	var dummyInnerError2 = &appError{
		Code:    errorCodeGeneralFailure,
		Message: dummyInnerErrorMessage,
		InnerErrors: []*appError{
			&appError{Message: dummyInnerMostErrorMessage},
		},
	}
	var dummyInnerError3 = &appError{
		Message: "dummy inner error 3",
	}
	var dummyInnerErrors = []*appError{
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	}
	var dummyNewInnerError1 = errors.New("some new error 1")
	var dummyNewInnerError2 = errors.New("some new error 2")
	var dummyNewInnerError3 = errors.New("some new error 3")
	var dummyNewInnerErrors = []error{
		dummyNewInnerError1,
		nil,
		dummyNewInnerError2,
		nil,
		dummyNewInnerError3,
	}
	var dummyCleanedInnerErrors = []*appError{
		&appError{Message: dummyNewInnerError1.Error()},
		&appError{Message: dummyNewInnerError2.Error()},
		&appError{Message: dummyNewInnerError3.Error()},
	}

	// mock
	createMock(t)

	// expect
	cleanupInnerErrorsFuncExpected = 1
	cleanupInnerErrorsFunc = func(innerErrors []error) []*appError {
		cleanupInnerErrorsFuncCalled++
		assert.Equal(t, dummyNewInnerErrors, innerErrors)
		return dummyCleanedInnerErrors
	}

	// SUT
	var appError = &appError{
		Code:        dummyCode,
		Message:     dummyMessage,
		InnerErrors: dummyInnerErrors,
	}

	// act
	appError.Wrap(
		dummyNewInnerErrors...,
	)

	// assert
	assert.Equal(t, 6, len(appError.InnerErrors))
	assert.Equal(t, dummyInnerErrors[0], appError.InnerErrors[0])
	assert.Equal(t, dummyInnerErrors[1], appError.InnerErrors[1])
	assert.Equal(t, dummyInnerErrors[2], appError.InnerErrors[2])
	assert.Equal(t, dummyCleanedInnerErrors[0], appError.InnerErrors[3])
	assert.Equal(t, dummyCleanedInnerErrors[1], appError.InnerErrors[4])
	assert.Equal(t, dummyCleanedInnerErrors[2], appError.InnerErrors[5])

	// verify
	verifyAll(t)
}

func TestGetGeneralFailure(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetGeneralFailure(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetUnauthorized(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeUnauthorized, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetUnauthorized(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetInvalidOperation(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeInvalidOperation, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetInvalidOperation(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetBadRequest(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetBadRequest(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetNotFound(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeNotFound, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetNotFound(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetCircuitBreak(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeCircuitBreak, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetCircuitBreak(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetOperationLock(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeOperationLock, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetOperationLock(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetAccessForbidden(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeAccessForbidden, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetAccessForbidden(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetDataCorruption(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeDataCorruption, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetDataCorruption(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}

func TestGetNotImplemented(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeNotImplemented, errorCode)
		assert.Equal(t, dummyErrorMessage, errorMessageor)
		assert.Equal(t, 3, len(innerErrors))
		assert.Equal(t, dummyInnerError1, innerErrors[0])
		assert.Equal(t, dummyInnerError2, innerErrors[1])
		assert.Equal(t, dummyInnerError3, innerErrors[2])
		return dummyResult
	}

	// SUT + act
	var appError, ok = GetNotImplemented(
		dummyErrorMessage,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	).(*appError)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyResult, appError)

	// verify
	verifyAll(t)
}
