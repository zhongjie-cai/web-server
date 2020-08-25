package webserver

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestBody_NilRequest(t *testing.T) {
	// arrange
	var dummySessionID *http.Request

	// mock
	createMock(t)

	// SUT + act
	var result = getRequestBody(
		dummySessionID,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetRequestBody_NilBody(t *testing.T) {
	// arrange
	var dummySessionID = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}

	// mock
	createMock(t)

	// SUT + act
	var result = getRequestBody(
		dummySessionID,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetRequestBody_ErrorBody(t *testing.T) {
	// arrange
	var bodyContent = "some body content"
	var dummySessionID = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
		Body:       ioutil.NopCloser(strings.NewReader(bodyContent)),
	}
	var dummyError = errors.New("some error message")

	// mock
	createMock(t)

	// expect
	ioutilReadAllExpected = 1
	ioutilReadAll = func(r io.Reader) ([]byte, error) {
		ioutilReadAllCalled++
		assert.Equal(t, dummySessionID.Body, r)
		return nil, dummyError
	}

	// SUT + act
	var result = getRequestBody(
		dummySessionID,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetRequestBody_Success(t *testing.T) {
	// arrange
	var bodyContent = "some body content"
	var dummySessionID = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
		Body:       ioutil.NopCloser(strings.NewReader(bodyContent)),
	}
	var dummyBuffer = &bytes.Buffer{}
	var dummyReadCloser = ioutil.NopCloser(nil)

	// mock
	createMock(t)

	// expect
	ioutilReadAllExpected = 1
	ioutilReadAll = func(r io.Reader) ([]byte, error) {
		ioutilReadAllCalled++
		return ioutil.ReadAll(r)
	}
	bytesNewBufferExpected = 1
	bytesNewBuffer = func(buf []byte) *bytes.Buffer {
		bytesNewBufferCalled++
		assert.Equal(t, []byte(bodyContent), buf)
		return dummyBuffer
	}
	ioutilNopCloserExpected = 1
	ioutilNopCloser = func(r io.Reader) io.ReadCloser {
		ioutilNopCloserCalled++
		assert.Equal(t, dummyBuffer, r)
		return dummyReadCloser
	}

	// SUT + act
	var result = getRequestBody(
		dummySessionID,
	)

	// assert
	assert.Equal(t, bodyContent, result)
	assert.Equal(t, dummyReadCloser, dummySessionID.Body)

	// verify
	verifyAll(t)
}
