package webserver

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInterfaceValueNil_NilInterface(t *testing.T) {
	// arrange
	var dummyInterface http.ResponseWriter

	// mock
	createMock(t)

	// expect

	// SUT + act
	var result = isInterfaceValueNil(
		dummyInterface,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestIsInterfaceValueNil_NilValue(t *testing.T) {
	// arrange
	var dummyInterface *dummyResponseWriter

	// mock
	createMock(t)

	// expect
	reflectValueOfExpected = 1
	reflectValueOf = func(i interface{}) reflect.Value {
		reflectValueOfCalled++
		assert.Equal(t, dummyInterface, i)
		return reflect.ValueOf(i)
	}

	// SUT + act
	var result = isInterfaceValueNil(
		dummyInterface,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestIsInterfaceValueNil_EmptyValue(t *testing.T) {
	// arrange
	var dummyInterface = 0

	// mock
	createMock(t)

	// expect
	reflectValueOfExpected = 1
	reflectValueOf = func(i interface{}) reflect.Value {
		reflectValueOfCalled++
		assert.Equal(t, dummyInterface, i)
		return reflect.Value{}
	}

	// SUT + act
	var result = isInterfaceValueNil(
		dummyInterface,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestIsInterfaceValueNil_ValidValue(t *testing.T) {
	// arrange
	var dummyInterface = 0

	// mock
	createMock(t)

	// expect
	reflectValueOfExpected = 1
	reflectValueOf = func(i interface{}) reflect.Value {
		reflectValueOfCalled++
		assert.Equal(t, dummyInterface, i)
		return reflect.ValueOf(dummyInterface)
	}

	// SUT + act
	var result = isInterfaceValueNil(
		dummyInterface,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}
