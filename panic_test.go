package webserver

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/zhongjie-cai/gomocker/v2"
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
	m.Mock((*DefaultCustomization).RecoverPanic).Expects(dummyCustomization, dummySession, dummyRecoverResult).
		Returns(dummyResponseObject, dummyResponseError).Once()
	m.Mock(writeResponse).Expects(dummySession, dummyResponseObject, dummyResponseError).Returns().Once()

	// SUT + act
	handlePanic(
		dummySession,
		dummyRecoverResult,
	)
}
