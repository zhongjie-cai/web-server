package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker/v2"
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
		{Message: "some inner error message 1"},
		{Message: "some inner error message 2"},
		{Message: "some inner error message 3"},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(cleanupInnerErrors).Expects(dummyInnerErrors).Returns(dummyCleanedUpErrors).Once()

	// SUT + act
	var result = newAppError(
		dummyErrorCode,
		dummyErrorMessage,
		dummyInnerErrors...,
	)

	// assert
	assert.Equal(t, dummyErrorCode, result.Code)
	assert.Equal(t, dummyErrorMessage, result.Message)
	assert.Equal(t, dummyCleanedUpErrors, result.InnerErrors)
}

func TestGetErrorMessage(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"

	// SUT
	var sut = errors.New(dummyErrorMessage)

	// act
	var result = getErrorMessage(sut)

	// assert
	assert.Equal(t, dummyErrorMessage, result)
}

func TestPrintInnerErrors_NilInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors []*appError

	// SUT + act
	var result = printInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Zero(t, result)
}

func TestPrintInnerErrors_EmptyInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors = []*appError{}

	// SUT + act
	var result = printInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Zero(t, result)
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
	var dummyResult = " -> [ some joined message ]"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	for i, dummyInnerError := range dummyInnerErrors {
		m.Mock(getErrorMessage).Expects(dummyInnerError).Returns(dummyErrorMessages[i]).Once()
	}
	m.Mock(strings.Join).Expects(dummyErrorMessages, errorSeparator).Returns(dummyJoinedMessage).Once()

	// SUT + act
	var result = printInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestAppError_Error(t *testing.T) {
	// arrange
	var dummyCode = errorCode("some error code")
	var dummyMessage = "some error message"
	var dummyInnerErrors = []*appError{
		{Message: "some inner error 1"},
		{Message: "some inner error 2"},
		{Message: "some inner error 3"},
	}
	var dummyAppError = &appError{
		dummyCode,
		dummyMessage,
		dummyInnerErrors,
	}
	var dummyBaseErrorMessage = "(some error code) some error message"
	var dummyInnerErrorMessage = "some inner error message"
	var dummyResult = "some result"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(printInnerErrors).Expects(dummyInnerErrors).Returns(dummyInnerErrorMessage).Once()
	m.Mock(fmt.Sprint).Expects(dummyBaseErrorMessage, dummyInnerErrorMessage).Returns(dummyResult).Once()

	// SUT + act
	var result = dummyAppError.Error()

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestAppError_ErrorCode(t *testing.T) {
	// arrange
	var dummyCode = errorCode("some error code")
	var dummyMessage = "some error message"

	// SUT
	var appError = &appError{
		Code:    dummyCode,
		Message: dummyMessage,
	}

	// act
	var code = appError.ErrorCode()

	// assert
	assert.Equal(t, string(dummyCode), code)
}

func TestAppError_HTTPStatusCode(t *testing.T) {
	// arrange
	var dummyCode = errorCode("some error code")
	var dummyMessage = "some error message"

	// SUT
	var appError = &appError{
		Code:    dummyCode,
		Message: dummyMessage,
	}

	// act
	var code = appError.HTTPStatusCode()

	// assert
	assert.Equal(t, dummyCode.httpStatusCode(), code)
}

func TestAppError_HTTPResponseMessage(t *testing.T) {
	// arrange
	var dummyCode = errorCode("some error code")
	var dummyMessage = "some error message"
	var dummyInnerErrors = []*appError{
		{Message: "some inner error 1"},
		{Message: "some inner error 2"},
		{Message: "some inner error 3"},
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
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(json.Marshal).Expects(dummyAppError).Returns(dummyBytes, dummyError).Once()

	// SUT
	var sut = dummyAppError

	// act
	var result = sut.HTTPResponseMessage()

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestEqualsError_SameError(t *testing.T) {
	// arrange
	var dummyError = errors.New("some error")
	var dummyTarget = dummyError

	// SUT + act
	var result = equalsError(
		dummyError,
		dummyTarget,
	)

	// assert
	assert.True(t, result)
}

func TestEqualsError_SameMessage(t *testing.T) {
	// arrange
	var dummyError = errors.New("some error")
	var dummyTarget = errors.New("some error")

	// SUT + act
	var result = equalsError(
		dummyError,
		dummyTarget,
	)

	// assert
	assert.True(t, result)
}

func TestEqualsError_ErrorIs(t *testing.T) {
	// arrange
	var dummyError = errors.New("some error")
	var dummyTarget = errors.New("some target")
	var dummyResult = rand.Intn(100) > 50

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(errors.Is).Expects(dummyError, dummyTarget).Returns(dummyResult).Once()

	// SUT + act
	var result = equalsError(
		dummyError,
		dummyTarget,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestAppErrorContains_HappyPath(t *testing.T) {
	// arrange
	var dummyError = errors.New("some error")
	var dummyappError = &appError{
		Code:    errorCode("some error code"),
		Message: dummyError.Error(),
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(equalsError).Expects(dummyappError, dummyError).Returns(true).Once()

	// SUT + act
	var result = appErrorContains(
		dummyappError,
		dummyError,
	)

	// assert
	assert.True(t, result)
}

func TestInnerErrorContains_NilInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors []*appError
	var dummyError = errors.New("some error")

	// SUT + act
	var result = innerErrorContains(
		dummyInnerErrors,
		dummyError,
	)

	// assert
	assert.False(t, result)
}

func TestInnerErrorContains_EmptyInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors = []*appError{}
	var dummyError = errors.New("some error")

	// SUT + act
	var result = innerErrorContains(
		dummyInnerErrors,
		dummyError,
	)

	// assert
	assert.False(t, result)
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
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(appErrorContains).Expects(dummyInnerError, dummyError).Returns(true).Once()

	// SUT + act
	var result = innerErrorContains(
		dummyInnerErrors,
		dummyError,
	)

	// assert
	assert.True(t, result)
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
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(appErrorContains).Expects(dummyInnerError1, dummyError).Returns(false).Once()
	m.Mock(appErrorContains).Expects(dummyInnerError2, dummyError).Returns(false).Once()

	// SUT + act
	var result = innerErrorContains(
		dummyInnerErrors,
		dummyError,
	)

	// assert
	assert.False(t, result)
}

func TestAppError_Contains_ErrorEqual(t *testing.T) {
	// arrange
	var dummyappError = &appError{
		Code:    errorCode("some error code"),
		Message: "some error message",
	}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(equalsError).Expects(dummyappError, dummyError).Returns(true).Once()

	// SUT
	var sut = dummyappError

	// act
	var result = sut.Contains(
		dummyError,
	)

	// assert
	assert.True(t, result)
}

func TestAppError_Contains_InnerErrorEqual(t *testing.T) {
	// arrange
	var dummyInnerErrors = []*appError{
		{Message: "some inner error 1"},
		{Message: "some inner error 2"},
		{Message: "some inner error 3"},
	}
	var dummyappError = &appError{
		Code:        errorCode("some error code"),
		Message:     "some error message",
		InnerErrors: dummyInnerErrors,
	}
	var dummyError = errors.New("some error")
	var dummyResult = rand.Intn(100) > 50

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(equalsError).Expects(dummyappError, dummyError).Returns(false).Once()
	m.Mock(innerErrorContains).Expects(dummyInnerErrors, dummyError).Returns(dummyResult).Once()

	// SUT
	var sut = dummyappError

	// act
	var result = sut.Contains(
		dummyError,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestCleanupInnerErrors_NilInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors []error

	// SUT + act
	var result = cleanupInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Empty(t, result)
}

func TestCleanupInnerErrors_EmptyInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors = []error{}

	// SUT + act
	var result = cleanupInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Empty(t, result)
}

func TestCleanupInnerErrors_NoValidInnerErrors(t *testing.T) {
	// arrange
	var dummyInnerErrors = []error{
		nil,
		nil,
		nil,
	}

	// SUT + act
	var result = cleanupInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Empty(t, result)
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

	// SUT + act
	var result = cleanupInnerErrors(
		dummyInnerErrors,
	)

	// assert
	assert.Equal(t, 3, len(result))
	assert.Equal(t, dummyInnerError1.Error(), result[0].Message)
	assert.Equal(t, dummyInnerError2.Error(), result[1].Message)
	assert.Equal(t, dummyInnerError3.Error(), result[2].Message)
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
			{Message: dummyInnerMostErrorMessage},
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
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(cleanupInnerErrors).Expects([]error{nil, nil, nil}).Returns(dummyCleanedInnerErrors).Once()

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
			{Message: dummyInnerMostErrorMessage},
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
		{Message: dummyNewInnerError1.Error()},
		{Message: dummyNewInnerError2.Error()},
		{Message: dummyNewInnerError3.Error()},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(cleanupInnerErrors).Expects(dummyNewInnerErrors).Returns(dummyCleanedInnerErrors).Once()

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
}

func TestGetGeneralFailure(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeGeneralFailure,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetUnauthorized(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeUnauthorized,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetInvalidOperation(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeInvalidOperation,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetBadRequest(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeBadRequest,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetNotFound(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeNotFound,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetCircuitBreak(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeCircuitBreak,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetOperationLock(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeOperationLock,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetAccessForbidden(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeAccessForbidden,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetDataCorruption(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeDataCorruption,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestGetNotImplemented(t *testing.T) {
	// arrange
	var dummyErrorMessage = "some error message"
	var dummyInnerError1 = errors.New("dummy inner error 1")
	var dummyInnerError2 = errors.New("dummy inner error 2")
	var dummyInnerError3 = errors.New("dummy inner error 3")
	var dummyResult = &appError{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeNotImplemented,
		dummyErrorMessage,
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyResult).Once()

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
}

func TestWrapError_NormalError(t *testing.T) {
	// arrange
	var dummySourceError = errors.New("some source error")
	var dummyInnerError1 = errors.New("some inner error 1")
	var dummyInnerError2 = errors.New("some inner error 2")
	var dummyInnerError3 = errors.New("some inner error 3")
	var dummyAppError = &appError{Message: "some app error"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(
		errorCodeGeneralFailure,
		dummySourceError.Error(),
		[]error{
			dummyInnerError1,
			dummyInnerError2,
			dummyInnerError3,
		},
	).Returns(dummyAppError).Once()

	// SUT + act
	var err = WrapError(
		dummySourceError,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
}

func TestWrapError_AppError(t *testing.T) {
	// arrange
	var dummySourceError = &appError{Message: "some source error"}
	var dummyInnerError1 = errors.New("some inner error 1")
	var dummyInnerError2 = errors.New("some inner error 2")
	var dummyInnerError3 = errors.New("some inner error 3")
	var dummyAppError = &appError{Message: "some app error"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*appError).Wrap).Expects(dummySourceError, dummyInnerError1, dummyInnerError2, dummyInnerError3).Returns(dummyAppError).Once()

	// SUT + act
	var err = WrapError(
		dummySourceError,
		dummyInnerError1,
		dummyInnerError2,
		dummyInnerError3,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
}
