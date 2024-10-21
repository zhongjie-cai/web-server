package webserver

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilResponseWriter(t *testing.T) {
	// arrange
	var dummyBody = []byte("some body")
	var dummyStatus = rand.Int()

	// SUT
	var nilResponseWriter = &nilResponseWriter{}

	// act
	var header = nilResponseWriter.Header()
	var result, err = nilResponseWriter.Write(dummyBody)
	nilResponseWriter.WriteHeader(dummyStatus)

	// assert
	assert.Empty(t, header)
	assert.Zero(t, result)
	assert.NoError(t, err)
}
