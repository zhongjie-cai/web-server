package webserver

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

type dummyCustomizationRecoverPanic struct {
	dummyCustomization
	recoverPanic func(Session, interface{}) (interface{}, error)
}

func (dummyCustomizationRecoverPanic *dummyCustomizationRecoverPanic) RecoverPanic(session Session, recoverResult interface{}) (interface{}, error) {
	if dummyCustomizationRecoverPanic.recoverPanic != nil {
		return dummyCustomizationRecoverPanic.recoverPanic(session, recoverResult)
	}
	assert.Fail(dummyCustomizationRecoverPanic.t, "Unexpected call to RecoverPanic")
	return nil, nil
}

func TestHandlePanic_HappyPath(t *testing.T) {
	// arrange
	var dummyCustomizationRecoverPanic = &dummyCustomizationRecoverPanic{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationRecoverPanic,
	}
	var dummyRecoverResult = errors.New("some error").(interface{})
	var customizationRecoverPanicExpected int
	var customizationRecoverPanicCalled int
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")

	// mock
	createMock(t)

	// expect
	customizationRecoverPanicExpected = 1
	dummyCustomizationRecoverPanic.recoverPanic = func(session Session, recoverResult interface{}) (interface{}, error) {
		customizationRecoverPanicCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRecoverResult, recoverResult)
		return dummyResponseObject, dummyResponseError
	}
	writeResponseFuncExpected = 1
	writeResponseFunc = func(session *session, responseObject interface{}, responseError error) {
		writeResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyResponseError, responseError)
	}

	// SUT + act
	handlePanic(
		dummySession,
		dummyRecoverResult,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationRecoverPanicExpected, customizationRecoverPanicCalled, "Unexpected number of calls to method customization.RecoverPanic")
}
