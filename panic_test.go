package webserver

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func TestHandlePanic_NilRecoverResult(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyRecoverResult interface{}

	// SUT + act
	handlePanic(
		dummySession,
		dummyRecoverResult,
	)
}

func TestHandlePanic_HappyPath(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	var dummyRecoverResult = errors.New("some error").(interface{})
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "RecoverPanic", 1, func(self *DefaultCustomization, session Session, recoverResult interface{}) (interface{}, error) {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRecoverResult, recoverResult)
		return dummyResponseObject, dummyResponseError
	})
	m.ExpectFunc(writeResponse, 1, func(session *session, responseObject interface{}, responseError error) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyResponseError, responseError)
	})

	// SUT + act
	handlePanic(
		dummySession,
		dummyRecoverResult,
	)
}
