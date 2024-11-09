package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker/v2"
)

func TestSkipResponseHandling(t *testing.T) {
	// SUT + act
	var result, err = SkipResponseHandling()

	// assert
	assert.IsType(t, skipResponseHandlingDummy{}, result)
	assert.NoError(t, err)
}

func TestShouldSkipHandling_HasError(t *testing.T) {
	// arrange
	var dummyResponseObject skipResponseHandlingDummy
	var dummyResponseError = errors.New("some error")

	// SUT + act
	var result = shouldSkipHandling(
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.False(t, result)
}

func TestShouldSkipHandling_Yes(t *testing.T) {
	// arrange
	var dummyResponseObject skipResponseHandlingDummy
	var dummyResponseError error

	// SUT + act
	var result = shouldSkipHandling(
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.True(t, result)
}

func TestShouldSkipHandling_No(t *testing.T) {
	// arrange
	var dummyResponseObject = rand.Int()
	var dummyResponseError error

	// SUT + act
	var result = shouldSkipHandling(
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.False(t, result)
}

func TestConstructResponse_ResponseError(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")
	var dummyCode = rand.Int()
	var dummyMessage = "some message"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).InterpretError).Expects(dummyCustomization, dummyResponseError).Returns(dummyCode, dummyMessage).Once()

	// SUT + act
	var code, message = constructResponse(
		dummySession,
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.Equal(t, dummyCode, code)
	assert.Equal(t, dummyMessage, message)
}

func TestConstructResponse_Success(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyResponseObject = rand.Int()
	var dummyCode = rand.Int()
	var dummyMessage = "some message"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).InterpretSuccess).Expects(dummyCustomization, dummyResponseObject).Returns(dummyCode, dummyMessage).Once()

	// SUT + act
	var code, message = constructResponse(
		dummySession,
		dummyResponseObject,
		nil,
	)

	// assert
	assert.Equal(t, dummyCode, code)
	assert.Equal(t, dummyMessage, message)
}

func TestWriteResponse_SkipHandling(t *testing.T) {
	// arrange
	var dummyResponseWriter = &dummyResponseWriter{}
	var dummySession = &session{
		responseWriter: dummyResponseWriter,
	}
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(shouldSkipHandling).Expects(dummyResponseObject, dummyResponseError).Returns(true).Once()
	m.Mock(logEndpointResponse).Expects(dummySession, "None", "-1", "Skipped response handling").Returns().Once()

	// SUT + act
	writeResponse(
		dummySession,
		dummyResponseObject,
		dummyResponseError,
	)
}

func TestWriteResponse_HappyPath(t *testing.T) {
	// arrange
	var dummyResponseWriterInstance = &dummyResponseWriter{}
	var dummySession = &session{
		responseWriter: dummyResponseWriterInstance,
	}
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")
	var dummyCode = rand.Int()
	var dummyMessage = "some message"
	var dummyStatusText = "some status text"
	var dummyCodeString = strconv.Itoa(dummyCode)
	var dummyHeader = make(http.Header)

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(shouldSkipHandling).Expects(dummyResponseObject, dummyResponseError).Returns(false).Once()
	m.Mock(constructResponse).Expects(dummySession, dummyResponseObject, dummyResponseError).Returns(dummyCode, dummyMessage).Once()
	m.Mock(http.StatusText).Expects(dummyCode).Returns(dummyStatusText).Once()
	m.Mock(strconv.Itoa).Expects(dummyCode).Returns(dummyCodeString).Once()
	m.Mock(logEndpointResponse).Expects(dummySession, dummyStatusText, dummyCodeString, dummyMessage).Returns().Once()
	m.Mock(isInterfaceValueNil).Expects(dummyResponseWriterInstance).Returns(false).Once()
	m.Mock((*dummyResponseWriter).Header).Expects(dummyResponseWriterInstance).Returns(dummyHeader).Once()
	m.Mock((*dummyResponseWriter).WriteHeader).Expects(dummyResponseWriterInstance, dummyCode).Returns().Once()
	m.Mock((*dummyResponseWriter).Write).Expects(dummyResponseWriterInstance, []byte(dummyMessage)).Returns(rand.Int(), errors.New("some error")).Once()

	// SUT + act
	writeResponse(
		dummySession,
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.Equal(t, 1, len(dummyHeader))
	assert.Equal(t, ContentTypeJSON, dummyHeader.Get("Content-Type"))
}
