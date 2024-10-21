package webserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInterfaceValueNil_NilInterface(t *testing.T) {
	// arrange
	var dummyInterface http.ResponseWriter

	// SUT + act
	var result = isInterfaceValueNil(
		dummyInterface,
	)

	// assert
	assert.True(t, result)
}

func TestIsInterfaceValueNil_NilValue(t *testing.T) {
	// arrange
	var dummyInterface *dummyResponseWriter

	// SUT + act
	var result = isInterfaceValueNil(
		dummyInterface,
	)

	// assert
	assert.True(t, result)
}

func TestIsInterfaceValueNil_EmptyValue(t *testing.T) {
	// arrange
	var dummyInterface *int

	// SUT + act
	var result = isInterfaceValueNil(
		dummyInterface,
	)

	// assert
	assert.True(t, result)
}

func TestIsInterfaceValueNil_ValidValue(t *testing.T) {
	// arrange
	var dummyInterface = 0

	// SUT + act
	var result = isInterfaceValueNil(
		dummyInterface,
	)

	// assert
	assert.False(t, result)
}
