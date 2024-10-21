package webserver

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func TestGetRequestBody_NilRequest(t *testing.T) {
	// arrange
	var dummyRequest *http.Request

	// SUT + act
	var result = getRequestBody(
		dummyRequest,
	)

	// assert
	assert.Zero(t, result)
}

func TestGetRequestBody_NilBody(t *testing.T) {
	// arrange
	var dummyRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}

	// SUT + act
	var result = getRequestBody(
		dummyRequest,
	)

	// assert
	assert.Zero(t, result)
}

func TestGetRequestBody_ErrorBody(t *testing.T) {
	// arrange
	var bodyContent = "some body content"
	var dummyRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
		Body:       io.NopCloser(strings.NewReader(bodyContent)),
	}
	var dummyError = errors.New("some error message")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(io.ReadAll, 1, func(r io.Reader) ([]byte, error) {
		assert.Equal(t, dummyRequest.Body, r)
		return nil, dummyError
	})

	// SUT + act
	var result = getRequestBody(
		dummyRequest,
	)

	// assert
	assert.Zero(t, result)
}

func TestGetRequestBody_Success(t *testing.T) {
	// arrange
	var bodyContent = "some body content"
	var dummyRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
		Body:       io.NopCloser(strings.NewReader(bodyContent)),
	}
	var dummyBuffer = &bytes.Buffer{}
	var dummyReadCloser = io.NopCloser(nil)
	var dummyReadAll, dummyError = io.ReadAll(dummyRequest.Body)

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(io.ReadAll, 1, func(r io.Reader) ([]byte, error) {
		assert.Equal(t, dummyRequest.Body, r)
		return dummyReadAll, dummyError
	})
	m.ExpectFunc(bytes.NewBuffer, 1, func(buf []byte) *bytes.Buffer {
		assert.Equal(t, []byte(bodyContent), buf)
		return dummyBuffer
	})
	m.ExpectFunc(io.NopCloser, 1, func(r io.Reader) io.ReadCloser {
		assert.Equal(t, dummyBuffer, r)
		return dummyReadCloser
	})

	// SUT + act
	var result = getRequestBody(
		dummyRequest,
	)

	// assert
	assert.Equal(t, bodyContent, result)
	assert.Equal(t, dummyReadCloser, dummyRequest.Body)
}
