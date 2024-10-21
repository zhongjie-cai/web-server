package webserver

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func TestGetClientForRequest_SendClientCert(t *testing.T) {
	// arrange
	var dummyHTTPClient1 = &http.Client{Timeout: time.Duration(rand.Int())}
	var dummyHTTPClient2 = &http.Client{Timeout: time.Duration(rand.Int())}

	// stub
	httpClientWithCert = dummyHTTPClient1
	httpClientNoCert = dummyHTTPClient2

	// SUT + act
	var result = getClientForRequest(true)

	// assert
	assert.Equal(t, dummyHTTPClient1, result)
}

func TestGetClientForRequest_NoSendClientCert(t *testing.T) {
	// arrange
	var dummyHTTPClient1 = &http.Client{Timeout: time.Duration(rand.Int())}
	var dummyHTTPClient2 = &http.Client{Timeout: time.Duration(rand.Int())}

	// stub
	httpClientWithCert = dummyHTTPClient1
	httpClientNoCert = dummyHTTPClient2

	// SUT + act
	var result = getClientForRequest(false)

	// assert
	assert.Equal(t, dummyHTTPClient2, result)
}

func TestClientDo(t *testing.T) {
	// arrange
	var dummyClient = &http.Client{}
	var dummyRequest, _ = http.NewRequest(
		"",
		"",
		nil,
	)

	// assert
	assert.NotPanics(
		t,
		func() {
			// SUT + act
			clientDo(
				dummyClient,
				dummyRequest,
			)
		},
	)
}

func TestClientDoWithRetry_ConnError_NoRetry(t *testing.T) {
	// arrange
	var dummyClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyConnRetry = 0
	var dummyHTTPRetry = map[int]int{}
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyResponseObject = &http.Response{}
	var dummyResponseError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(clientDo, 1, func(client *http.Client, request *http.Request) (*http.Response, error) {
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, dummyResponseError
	})

	// SUT + act
	var result, err = clientDoWithRetry(
		dummyClient,
		dummyRequestObject,
		dummyConnRetry,
		dummyHTTPRetry,
		dummyRetryDelay,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.Equal(t, dummyResponseError, err)
}

func TestClientDoWithRetry_ConnError_RetryOK(t *testing.T) {
	// arrange
	var dummyClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyConnRetry = 2
	var dummyHTTPRetry = map[int]int{}
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyResponseObject = &http.Response{}
	var dummyResponseError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(clientDo, 2, func(client *http.Client, request *http.Request) (*http.Response, error) {
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		if m.FuncCalledCount(clientDo) == 1 {
			return dummyResponseObject, dummyResponseError
		} else if m.FuncCalledCount(clientDo) == 2 {
			return dummyResponseObject, nil
		}
		return nil, nil
	})
	m.ExpectFunc(time.Sleep, 1, func(d time.Duration) {
		assert.Equal(t, dummyRetryDelay, d)
	})

	// SUT + act
	var result, err = clientDoWithRetry(
		dummyClient,
		dummyRequestObject,
		dummyConnRetry,
		dummyHTTPRetry,
		dummyRetryDelay,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.NoError(t, err)
}

func TestClientDoWithRetry_ConnError_RetryFail(t *testing.T) {
	// arrange
	var dummyClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyConnRetry = 2
	var dummyHTTPRetry = map[int]int{}
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyResponseObject = &http.Response{}
	var dummyResponseError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(clientDo, 3, func(client *http.Client, request *http.Request) (*http.Response, error) {
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, dummyResponseError
	})
	m.ExpectFunc(time.Sleep, 2, func(d time.Duration) {
		assert.Equal(t, dummyRetryDelay, d)
	})

	// SUT + act
	var result, err = clientDoWithRetry(
		dummyClient,
		dummyRequestObject,
		dummyConnRetry,
		dummyHTTPRetry,
		dummyRetryDelay,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.Equal(t, dummyResponseError, err)
}

func TestClientDoWithRetry_HTTPError_NilResponse(t *testing.T) {
	// arrange
	var dummyClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyConnRetry = rand.Int()
	var dummyHTTPRetry = map[int]int{}
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyResponseObject *http.Response

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(clientDo, 1, func(client *http.Client, request *http.Request) (*http.Response, error) {
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, nil
	})

	// SUT + act
	var result, err = clientDoWithRetry(
		dummyClient,
		dummyRequestObject,
		dummyConnRetry,
		dummyHTTPRetry,
		dummyRetryDelay,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.NoError(t, err)
}

func TestClientDoWithRetry_HTTPError_NoRetry(t *testing.T) {
	// arrange
	var dummyClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyConnRetry = rand.Int()
	var dummyHTTPRetry = map[int]int{}
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyResponseObject = &http.Response{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(clientDo, 1, func(client *http.Client, request *http.Request) (*http.Response, error) {
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, nil
	})

	// SUT + act
	var result, err = clientDoWithRetry(
		dummyClient,
		dummyRequestObject,
		dummyConnRetry,
		dummyHTTPRetry,
		dummyRetryDelay,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.NoError(t, err)
}

func TestClientDoWithRetry_HTTPError_RetryOK(t *testing.T) {
	// arrange
	var dummyClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyConnRetry = rand.Int()
	var dummyStatusCode = rand.Int()
	var dummyHTTPRetry = map[int]int{
		dummyStatusCode: 2,
	}
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyResponseObject1 = &http.Response{
		StatusCode: dummyStatusCode,
	}
	var dummyResponseObject2 = &http.Response{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(clientDo, 2, func(client *http.Client, request *http.Request) (*http.Response, error) {
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		if m.FuncCalledCount(clientDo) == 1 {
			return dummyResponseObject1, nil
		} else if m.FuncCalledCount(clientDo) == 2 {
			return dummyResponseObject2, nil
		}
		return nil, nil
	})
	m.ExpectFunc(time.Sleep, 1, func(d time.Duration) {
		assert.Equal(t, dummyRetryDelay, d)
	})

	// SUT + act
	var result, err = clientDoWithRetry(
		dummyClient,
		dummyRequestObject,
		dummyConnRetry,
		dummyHTTPRetry,
		dummyRetryDelay,
	)

	// assert
	assert.Equal(t, dummyResponseObject2, result)
	assert.NoError(t, err)
}

func TestClientDoWithRetry_HTTPError_RetryFail(t *testing.T) {
	// arrange
	var dummyClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyConnRetry = rand.Int()
	var dummyStatusCode = rand.Int()
	var dummyHTTPRetry = map[int]int{
		dummyStatusCode: 2,
	}
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyResponseObject = &http.Response{
		StatusCode: dummyStatusCode,
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(clientDo, 3, func(client *http.Client, request *http.Request) (*http.Response, error) {
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, nil
	})
	m.ExpectFunc(time.Sleep, 2, func(d time.Duration) {
		assert.Equal(t, dummyRetryDelay, d)
	})

	// SUT + act
	var result, err = clientDoWithRetry(
		dummyClient,
		dummyRequestObject,
		dummyConnRetry,
		dummyHTTPRetry,
		dummyRetryDelay,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.NoError(t, err)
}

func TestGetHTTPTransport_NoClientCert(t *testing.T) {
	// arrange
	var dummySkipServerCertVerification = rand.Intn(100) < 50
	var dummyClientCert *tls.Certificate
	var dummyRoundTripper = &http.Transport{}
	var dummyRoundTripperWrapperExpected = 1
	var dummyRoundTripperWrapperCalled = 0

	// expect
	var dummyRoundTripperWrapper = func(rt http.RoundTripper) http.RoundTripper {
		dummyRoundTripperWrapperCalled++
		assert.Equal(t, http.DefaultTransport, rt)
		return dummyRoundTripper
	}

	// SUT + act
	var result = getHTTPTransport(
		dummySkipServerCertVerification,
		dummyClientCert,
		dummyRoundTripperWrapper,
	)

	// assert
	assert.Equal(t, dummyRoundTripper, result)
	assert.Equal(t, dummyRoundTripperWrapperExpected, dummyRoundTripperWrapperCalled)
}

func TestGetHTTPTransport_WithClientCert(t *testing.T) {
	// arrange
	var dummySkipServerCertVerification = rand.Intn(100) < 50
	var dummyClientCert = &tls.Certificate{}
	var dummyRoundTripper = &http.Transport{}
	var dummyRoundTripperWrapperExpected = 1
	var dummyRoundTripperWrapperCalled = 0

	// expect
	var dummyRoundTripperWrapper = func(rt http.RoundTripper) http.RoundTripper {
		dummyRoundTripperWrapperCalled++
		assert.NotEqual(t, http.DefaultTransport, rt)
		return dummyRoundTripper
	}

	// SUT + act
	var result = getHTTPTransport(
		dummySkipServerCertVerification,
		dummyClientCert,
		dummyRoundTripperWrapper,
	)

	// assert
	assert.Equal(t, dummyRoundTripper, result)
	assert.Equal(t, dummyRoundTripperWrapperExpected, dummyRoundTripperWrapperCalled)
}

func TestInitializeHTTPClients(t *testing.T) {
	// arrange
	var dummyWebcallTimeout = time.Duration(rand.Int())
	var dummySkipServerCertVerification = rand.Intn(100) < 50
	var dummyClientCert = &tls.Certificate{}
	var dummyHTTPTransport1 = &http.Transport{MaxConnsPerHost: rand.Int()}
	var dummyHTTPTransport2 = &http.Transport{MaxConnsPerHost: rand.Int()}
	var dummyRoundTripperWrapper = func(http.RoundTripper) http.RoundTripper { return nil }

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(getHTTPTransport, 2, func(skipServerCertVerification bool, clientCertificate *tls.Certificate, roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper) http.RoundTripper {
		assert.Equal(t, dummySkipServerCertVerification, skipServerCertVerification)
		if m.FuncCalledCount(getHTTPTransport) == 1 {
			assert.Equal(t, dummyClientCert, clientCertificate)
			return dummyHTTPTransport1
		} else if m.FuncCalledCount(getHTTPTransport) == 2 {
			assert.Nil(t, clientCertificate)
			return dummyHTTPTransport2
		}
		functionPointerEquals(t, dummyRoundTripperWrapper, roundTripperWrapper)
		return nil
	})

	// SUT + act
	initializeHTTPClients(
		dummyWebcallTimeout,
		dummySkipServerCertVerification,
		dummyClientCert,
		dummyRoundTripperWrapper,
	)

	// assert
	assert.NotNil(t, httpClientWithCert)
	assert.Equal(t, dummyHTTPTransport1, httpClientWithCert.Transport)
	assert.Equal(t, dummyWebcallTimeout, httpClientWithCert.Timeout)
	assert.NotNil(t, httpClientNoCert)
	assert.Equal(t, dummyHTTPTransport2, httpClientNoCert.Transport)
	assert.Equal(t, dummyWebcallTimeout, httpClientNoCert.Timeout)
}

func TestWebREquestAddQuery_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue1 = "some value 1"
	var dummyValue2 = "some value 2"
	var dummyValue3 = "some value 3"

	// SUT
	var sut = &webRequest{}

	// act
	var result, ok = sut.AddQuery(
		dummyName,
		dummyValue1,
	).AddQuery(
		dummyName,
		dummyValue2,
	).AddQuery(
		dummyName,
		dummyValue3,
	).(*webRequest)

	// assert
	assert.True(t, ok)
	assert.NotNil(t, result.query)
	assert.Equal(t, 1, len(result.query))
	var values, found = result.query[dummyName]
	assert.True(t, found)
	assert.Equal(t, 3, len(values))
	assert.Equal(t, dummyValue1, values[0])
	assert.Equal(t, dummyValue2, values[1])
	assert.Equal(t, dummyValue3, values[2])
}

func TestWebREquestAddHeader_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue1 = "some value 1"
	var dummyValue2 = "some value 2"
	var dummyValue3 = "some value 3"

	// SUT
	var sut = &webRequest{}

	// act
	var result, ok = sut.AddHeader(
		dummyName,
		dummyValue1,
	).AddHeader(
		dummyName,
		dummyValue2,
	).AddHeader(
		dummyName,
		dummyValue3,
	).(*webRequest)

	// assert
	assert.True(t, ok)
	assert.NotNil(t, result.header)
	assert.Equal(t, 1, len(result.header))
	var values, found = result.header[dummyName]
	assert.True(t, found)
	assert.Equal(t, 3, len(values))
	assert.Equal(t, dummyValue1, values[0])
	assert.Equal(t, dummyValue2, values[1])
	assert.Equal(t, dummyValue3, values[2])
}

func TestWebRequestSetupRetry(t *testing.T) {
	// arrange
	var dummyConnRetry = rand.Int()
	var dummyHTTPRetry = map[int]int{
		rand.Int(): rand.Int(),
		rand.Int(): rand.Int(),
	}
	var dummyRetryDelay = time.Duration(rand.Intn(100))

	// SUT
	var sut = &webRequest{}

	// act
	var result, ok = sut.SetupRetry(
		dummyConnRetry,
		dummyHTTPRetry,
		dummyRetryDelay,
	).(*webRequest)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummyConnRetry, result.connRetry)
	assert.Equal(t, dummyHTTPRetry, result.httpRetry)
	assert.Equal(t, dummyRetryDelay, result.retryDelay)
}

func TestWebRequestAnticipate(t *testing.T) {
	// arrange
	var dummyBeginStatusCode1 = rand.Int()
	var dummyEndStatusCode1 = rand.Int()
	var dummyDataTemplate1 string
	var dummyBeginStatusCode2 = rand.Int()
	var dummyEndStatusCode2 = rand.Int()
	var dummyDataTemplate2 int

	// SUT
	var sut = &webRequest{}

	// act
	var result, ok = sut.Anticipate(
		dummyBeginStatusCode1,
		dummyEndStatusCode1,
		&dummyDataTemplate1,
	).Anticipate(
		dummyBeginStatusCode2,
		dummyEndStatusCode2,
		&dummyDataTemplate2,
	).(*webRequest)

	// assert
	assert.True(t, ok)
	assert.Equal(t, 2, len(result.dataReceivers))
	assert.Equal(t, dummyBeginStatusCode1, result.dataReceivers[0].beginStatusCode)
	assert.Equal(t, dummyEndStatusCode1, result.dataReceivers[0].endStatusCode)
	assert.Equal(t, &dummyDataTemplate1, result.dataReceivers[0].dataTemplate)
	assert.Equal(t, dummyBeginStatusCode2, result.dataReceivers[1].beginStatusCode)
	assert.Equal(t, dummyEndStatusCode2, result.dataReceivers[1].endStatusCode)
	assert.Equal(t, &dummyDataTemplate2, result.dataReceivers[1].dataTemplate)
}

func TestCreateQueryString_NilQuery(t *testing.T) {
	// arrange
	var dummyQuery map[string][]string
	var dummyResult = "some result"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(strings.Join, 1, func(a []string, sep string) string {
		assert.Empty(t, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	})

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestCreateQueryString_EmptyQuery(t *testing.T) {
	// arrange
	var dummyQuery = map[string][]string{}
	var dummyResult = "some result"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(strings.Join, 1, func(a []string, sep string) string {
		assert.Empty(t, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	})

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestCreateQueryString_EmptyQueryName(t *testing.T) {
	// arrange
	var dummyQuery = map[string][]string{
		"": {"empty 1", "empty 2"},
	}
	var dummyResult = "some result"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(strings.Join, 1, func(a []string, sep string) string {
		assert.Empty(t, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	})

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestCreateQueryString_EmptyQueryValues(t *testing.T) {
	// arrange
	var dummyQuery = map[string][]string{
		"":          {"empty 1", "empty 2"},
		"some name": {},
	}
	var dummyResult = "some result"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(strings.Join, 1, func(a []string, sep string) string {
		assert.Empty(t, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	})

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestCreateQueryString_HappyPath(t *testing.T) {
	// arrange
	var dummyQuery = map[string][]string{
		"":          {"empty 1", "empty 2"},
		"some name": {"some value 1", "some value 2", "some value 3"},
	}
	var dummyResult = "some+name=some+value+1&some+name=some+value+2&some+name=some+value+3"

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestGenerateRequestURL_NilQuery(t *testing.T) {
	// arrange
	var dummyBaseURL = "some base URL"
	var dummyQuery map[string][]string = nil

	// SUT + act
	var result = generateRequestURL(
		dummyBaseURL,
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyBaseURL, result)
}

func TestGenerateRequestURL_EmptyQuery(t *testing.T) {
	// arrange
	var dummyBaseURL = "some base URL"
	var dummyQuery = map[string][]string{
		"foo":  {"bar 1", "bar 2"},
		"test": {"123", "456", "789"},
	}
	var dummyQueryString string

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(createQueryString, 1, func(query map[string][]string) string {
		assert.Equal(t, dummyQuery, query)
		return dummyQueryString
	})

	// SUT + act
	var result = generateRequestURL(
		dummyBaseURL,
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyBaseURL, result)
}

func TestGenerateRequestURL_Success(t *testing.T) {
	// arrange
	var dummyBaseURL = "some base URL"
	var dummyQuery = map[string][]string{
		"foo":  {"bar 1", "bar 2"},
		"test": {"123", "456", "789"},
	}
	var dummyQueryString = "some query string"
	var dummyResult = dummyBaseURL + "?" + dummyQueryString

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(createQueryString, 1, func(query map[string][]string) string {
		assert.Equal(t, dummyQuery, query)
		return dummyQueryString
	})

	// SUT + act
	var result = generateRequestURL(
		dummyBaseURL,
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)
}

func TestCreateHTTPRequest_NilWebRequest(t *testing.T) {
	// arrange
	var dummyWebRequest *webRequest
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

	// SUT + act
	var result, err = createHTTPRequest(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)
}

func TestCreateHTTPRequest_NilWebRequestSession(t *testing.T) {
	// arrange
	var dummyWebRequest = &webRequest{}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

	// SUT + act
	var result, err = createHTTPRequest(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)
}

func TestCreateHTTPRequest_RequestError(t *testing.T) {
	// arrange
	var dummySession = &session{}
	var dummyMethod = "some method"
	var dummyURL = "some URL"
	var dummyPayload = "some payload"
	var dummyHeader = map[string][]string{
		"foo":  {"bar"},
		"test": {"123", "456"},
	}
	var dummyQuery = map[string][]string{
		"me":   {"god"},
		"what": {"xyz", "abc"},
	}
	var dummyConnRetry = rand.Int()
	var dummyHTTPRetry = map[int]int{
		rand.Int(): rand.Int(),
		rand.Int(): rand.Int(),
	}
	var dummySendClientCert = rand.Intn(100) < 50
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyDataReceivers = []dataReceiver{
		{0, 999, nil},
	}
	var dummyWebRequest = &webRequest{
		dummySession,
		dummyMethod,
		dummyURL,
		dummyPayload,
		dummyQuery,
		dummyHeader,
		dummyConnRetry,
		dummyHTTPRetry,
		dummySendClientCert,
		dummyRetryDelay,
		dummyDataReceivers,
	}
	var dummyRequestURL = "some request url"
	var dummyRequest *http.Request
	var dummyError = errors.New("some error message")

	// stub
	var dummyStingsReader = strings.NewReader(dummyPayload)

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(generateRequestURL, 1, func(baseURL string, query map[string][]string) string {
		assert.Equal(t, dummyURL, baseURL)
		assert.Equal(t, dummyQuery, query)
		return dummyRequestURL
	})
	m.ExpectFunc(strings.NewReader, 1, func(s string) *strings.Reader {
		return dummyStingsReader
	})
	m.ExpectFunc(http.NewRequest, 1, func(method, url string, body io.Reader) (*http.Request, error) {
		assert.Equal(t, dummyMethod, method)
		assert.Equal(t, dummyRequestURL, url)
		assert.NotNil(t, body)
		return dummyRequest, dummyError
	})

	// SUT + act
	var result, err = createHTTPRequest(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyError, err)
}

func TestCreateHTTPRequest_Success(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyMethod = "some method"
	var dummyURL = "some URL"
	var dummyPayload = "some payload"
	var dummyHeader = map[string][]string{
		"foo":  {"bar"},
		"test": {"123", "456"},
	}
	var dummyQuery = map[string][]string{
		"me":   {"god"},
		"what": {"xyz", "abc"},
	}
	var dummyConnRetry = rand.Int()
	var dummyHTTPRetry = map[int]int{
		rand.Int(): rand.Int(),
		rand.Int(): rand.Int(),
	}
	var dummySendClientCert = rand.Intn(100) < 50
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyDataReceivers = []dataReceiver{
		{0, 999, nil},
	}
	var dummyWebRequest = &webRequest{
		dummySession,
		dummyMethod,
		dummyURL,
		dummyPayload,
		dummyQuery,
		dummyHeader,
		dummyConnRetry,
		dummyHTTPRetry,
		dummySendClientCert,
		dummyRetryDelay,
		dummyDataReceivers,
	}
	var dummyRequestURL = "some request url"
	var dummyRequest = &http.Request{
		RequestURI: "abc",
	}
	var dummyHeaderContent = "some header content"
	var dummyCustomized = &http.Request{
		RequestURI: "def",
	}

	// stub
	var dummyStingsReader = strings.NewReader(dummyPayload)

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(generateRequestURL, 1, func(baseURL string, query map[string][]string) string {
		assert.Equal(t, dummyURL, baseURL)
		assert.Equal(t, dummyQuery, query)
		return dummyRequestURL
	})
	m.ExpectFunc(strings.NewReader, 1, func(s string) *strings.Reader {
		return dummyStingsReader
	})
	m.ExpectFunc(http.NewRequest, 1, func(method, url string, body io.Reader) (*http.Request, error) {
		assert.Equal(t, dummyMethod, method)
		assert.Equal(t, dummyRequestURL, url)
		assert.NotNil(t, body)
		return dummyRequest, nil
	})
	m.ExpectFunc(logWebcallStart, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethod, category)
		assert.Equal(t, dummyURL, subcategory)
		assert.Equal(t, dummyRequestURL, messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(logWebcallRequest, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Content", subcategory)
		assert.Empty(t, parameters)
		if m.FuncCalledCount(logWebcallRequest) == 1 {
			assert.Equal(t, "Payload", category)
			assert.Equal(t, dummyPayload, messageFormat)
		} else if m.FuncCalledCount(logWebcallRequest) == 2 {
			assert.Equal(t, "Header", category)
			assert.Equal(t, dummyHeaderContent, messageFormat)
		}
	})
	m.ExpectFunc(marshalIgnoreError, 1, func(v interface{}) string {
		return dummyHeaderContent
	})
	m.ExpectMethod(dummyCustomization, "WrapRequest", 1, func(self *DefaultCustomization, session Session, httpRequest *http.Request) *http.Request {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRequest, httpRequest)
		return dummyCustomized
	})

	// SUT + act
	var result, err = createHTTPRequest(
		dummyWebRequest,
	)

	// assert
	assert.Equal(t, dummyCustomized, result)
	assert.NoError(t, err)
}

func TestLogErrorResponse(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyError = errors.New("some error")
	var dummyStartTime = time.Now()
	var dummyTimeSince = time.Duration(rand.Intn(1000))

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(logWebcallResponse, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Error", category)
		assert.Equal(t, "Content", subcategory)
		assert.Equal(t, "%+v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	})
	m.ExpectFunc(time.Since, 1, func(ts time.Time) time.Duration {
		assert.Equal(t, dummyStartTime, ts)
		return dummyTimeSince
	})
	m.ExpectFunc(logWebcallFinish, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Error", category)
		assert.Equal(t, "-1", subcategory)
		assert.Equal(t, "%s", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyTimeSince, parameters[0])
	})

	// SUT + act
	logErrorResponse(
		dummySession,
		dummyError,
		dummyStartTime,
	)
}

func TestLogSuccessResponse_NilResponse(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyResponse *http.Response
	var dummyStartTime = time.Now()

	// SUT + act
	logSuccessResponse(
		dummySession,
		dummyResponse,
		dummyStartTime,
	)
}

func TestLogSuccessResponse_ValidResponse(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyStatus = "some status"
	var dummyStatusCode = rand.Intn(1000)
	var dummyBody = io.NopCloser(bytes.NewBufferString("some body"))
	var dummyHeader = http.Header{
		"foo":  []string{"bar"},
		"test": []string{"123", "456", "789"},
	}
	var dummyResponse = &http.Response{
		StatusCode: dummyStatusCode,
		Body:       dummyBody,
		Header:     dummyHeader,
	}
	var dummyResponseBytes = []byte("some response bytes")
	var dummyResponseBody = string(dummyResponseBytes)
	var dummyError = errors.New("some error")
	var dummyBuffer = &bytes.Buffer{}
	var dummyNewBody = io.NopCloser(bytes.NewBufferString("some new body"))
	var dummyStartTime = time.Now()
	var dummyHeaderContent = "some header content"
	var dummyTimeSince = time.Duration(rand.Intn(1000))

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(io.ReadAll, 1, func(r io.Reader) ([]byte, error) {
		assert.Equal(t, dummyBody, r)
		return dummyResponseBytes, dummyError
	})
	m.ExpectFunc(bytes.NewBuffer, 1, func(buf []byte) *bytes.Buffer {
		assert.Equal(t, dummyResponseBytes, buf)
		return dummyBuffer
	})
	m.ExpectFunc(io.NopCloser, 1, func(r io.Reader) io.ReadCloser {
		assert.Equal(t, dummyBuffer, r)
		return dummyNewBody
	})
	m.ExpectFunc(http.StatusText, 1, func(code int) string {
		assert.Equal(t, dummyStatusCode, code)
		return dummyStatus
	})
	m.ExpectFunc(logWebcallResponse, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Content", subcategory)
		assert.Empty(t, parameters)
		if m.FuncCalledCount(logWebcallResponse) == 1 {
			assert.Equal(t, "Header", category)
			assert.Equal(t, dummyHeaderContent, messageFormat)
		} else if m.FuncCalledCount(logWebcallResponse) == 2 {
			assert.Equal(t, "Body", category)
			assert.Equal(t, dummyResponseBody, messageFormat)
		}
	})
	m.ExpectFunc(marshalIgnoreError, 1, func(v interface{}) string {
		assert.Equal(t, dummyHeader, v)
		return dummyHeaderContent
	})
	m.ExpectFunc(time.Since, 1, func(ts time.Time) time.Duration {
		assert.Equal(t, dummyStartTime, ts)
		return dummyTimeSince
	})
	m.ExpectFunc(logWebcallFinish, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStatus, category)
		assert.Equal(t, strconv.Itoa(dummyStatusCode), subcategory)
		assert.Equal(t, "%s", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyTimeSince, parameters[0])
	})

	// SUT + act
	logSuccessResponse(
		dummySession,
		dummyResponse,
		dummyStartTime,
	)

	// assert
	assert.Equal(t, dummyNewBody, dummyResponse.Body)
}

func TestDoRequestProcessing_NilWebRequest(t *testing.T) {
	// arrange
	var dummyWebRequest *webRequest
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)
}

func TestDoRequestProcessing_NilWebRequestSession(t *testing.T) {
	// arrange
	var dummyWebRequest = &webRequest{}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)
}

func TestDoRequestProcessing_RequestError(t *testing.T) {
	// arrange
	var dummyWebRequest = &webRequest{
		session: &session{id: uuid.New()},
	}
	var dummyRequestObject *http.Request
	var dummyRequestError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(createHTTPRequest, 1, func(webRequest *webRequest) (*http.Request, error) {
		assert.Equal(t, dummyWebRequest, webRequest)
		return dummyRequestObject, dummyRequestError
	})

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyRequestError, err)
}

func TestDoRequestProcessing_ResponseError(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyConnRetry = rand.Int()
	var dummyHTTPRetry = map[int]int{
		rand.Int(): rand.Int(),
		rand.Int(): rand.Int(),
	}
	var dummySendClientCert = rand.Intn(100) < 50
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyWebRequest = &webRequest{
		session:        dummySession,
		connRetry:      dummyConnRetry,
		httpRetry:      dummyHTTPRetry,
		sendClientCert: dummySendClientCert,
		retryDelay:     dummyRetryDelay,
	}
	var dummyHTTPClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyResponseObject *http.Response
	var dummyResponseError = errors.New("some error")
	var dummyStartTime = time.Now()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(createHTTPRequest, 1, func(webRequest *webRequest) (*http.Request, error) {
		assert.Equal(t, dummyWebRequest, webRequest)
		return dummyRequestObject, nil
	})
	m.ExpectFunc(getClientForRequest, 1, func(sendClientCert bool) *http.Client {
		assert.Equal(t, dummySendClientCert, sendClientCert)
		return dummyHTTPClient
	})
	m.ExpectFunc(getTimeNowUTC, 1, func() time.Time {
		return dummyStartTime
	})
	m.ExpectFunc(clientDoWithRetry, 1, func(client *http.Client, request *http.Request, connRetry int, httpRetry map[int]int, retryDelay time.Duration) (*http.Response, error) {
		assert.Equal(t, dummyHTTPClient, client)
		assert.Equal(t, dummyRequestObject, request)
		assert.Equal(t, dummyConnRetry, connRetry)
		assert.Equal(t, dummyHTTPRetry, httpRetry)
		assert.Equal(t, dummyRetryDelay, retryDelay)
		return dummyResponseObject, dummyResponseError
	})
	m.ExpectFunc(logErrorResponse, 1, func(session *session, responseError error, startTime time.Time) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseError, responseError)
		assert.Equal(t, dummyStartTime, startTime)
	})

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.Equal(t, dummyResponseError, err)
}

func TestDoRequestProcessing_ResponseSuccess(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyConnRetry = rand.Int()
	var dummyHTTPRetry = map[int]int{
		rand.Int(): rand.Int(),
		rand.Int(): rand.Int(),
	}
	var dummySendClientCert = rand.Intn(100) < 50
	var dummyRetryDelay = time.Duration(rand.Intn(100))
	var dummyWebRequest = &webRequest{
		session:        dummySession,
		connRetry:      dummyConnRetry,
		httpRetry:      dummyHTTPRetry,
		sendClientCert: dummySendClientCert,
		retryDelay:     dummyRetryDelay,
	}
	var dummyHTTPClient = &http.Client{}
	var dummyRequestObject = &http.Request{}
	var dummyResponseObject = &http.Response{}
	var dummyStartTime = time.Now()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(createHTTPRequest, 1, func(webRequest *webRequest) (*http.Request, error) {
		assert.Equal(t, dummyWebRequest, webRequest)
		return dummyRequestObject, nil
	})
	m.ExpectFunc(getClientForRequest, 1, func(sendClientCert bool) *http.Client {
		assert.Equal(t, dummySendClientCert, sendClientCert)
		return dummyHTTPClient
	})
	m.ExpectFunc(getTimeNowUTC, 1, func() time.Time {
		return dummyStartTime
	})
	m.ExpectFunc(clientDoWithRetry, 1, func(client *http.Client, request *http.Request, connRetry int, httpRetry map[int]int, retryDelay time.Duration) (*http.Response, error) {
		assert.Equal(t, dummyHTTPClient, client)
		assert.Equal(t, dummyRequestObject, request)
		assert.Equal(t, dummyConnRetry, connRetry)
		assert.Equal(t, dummyHTTPRetry, httpRetry)
		assert.Equal(t, dummyRetryDelay, retryDelay)
		return dummyResponseObject, nil
	})
	m.ExpectFunc(logSuccessResponse, 1, func(session *session, response *http.Response, startTime time.Time) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, response)
		assert.Equal(t, dummyStartTime, startTime)
	})

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.NoError(t, err)
}

func TestGetDataTemplate_EmptyDataReceivers(t *testing.T) {
	// arrange
	var dummySession = &session{}
	var dummyStatusCode = rand.Intn(100)
	var dummyDataReceivers = []dataReceiver{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Nil(t, i)
		return true
	})
	m.ExpectFunc(logWebcallResponse, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Body", category)
		assert.Equal(t, "Receiver", subcategory)
		assert.Equal(t, "No data template registered for HTTP status %v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyStatusCode, parameters[0])
	})

	// SUT + act
	var result = getDataTemplate(
		dummySession,
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Nil(t, result)
}

func TestGetDataTemplate_NoMatch(t *testing.T) {
	// arrange
	var dummySession = &session{}
	var dummyStatusCode = rand.Intn(100)
	var dummyDataTemplate string
	var dummyDataReceivers = []dataReceiver{
		{
			beginStatusCode: rand.Intn(100) + 100,
			endStatusCode:   rand.Intn(100) + 200,
			dataTemplate:    &dummyDataTemplate,
		},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Nil(t, i)
		return true
	})
	m.ExpectFunc(logWebcallResponse, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Body", category)
		assert.Equal(t, "Receiver", subcategory)
		assert.Equal(t, "No data template registered for HTTP status %v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyStatusCode, parameters[0])
	})

	// SUT + act
	var result = getDataTemplate(
		dummySession,
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Nil(t, result)
}

func TestGetDataTemplate_SingleMatch(t *testing.T) {
	// arrange
	var dummySession = &session{}
	var dummyStatusCode = rand.Intn(100)
	var dummyDataTemplate1 string
	var dummyDataTemplate2 int
	var dummyDataReceivers = []dataReceiver{
		{
			beginStatusCode: dummyStatusCode - rand.Intn(10),
			endStatusCode:   dummyStatusCode + 1 + rand.Intn(10),
			dataTemplate:    &dummyDataTemplate1,
		},
		{
			beginStatusCode: rand.Intn(100) + 100,
			endStatusCode:   rand.Intn(100) + 200,
			dataTemplate:    &dummyDataTemplate2,
		},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, &dummyDataTemplate1, i)
		return false
	})

	// SUT + act
	var result = getDataTemplate(
		dummySession,
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Equal(t, &dummyDataTemplate1, result)
}

func TestGetDataTemplate_OverlapMatch(t *testing.T) {
	// arrange
	var dummySession = &session{}
	var dummyStatusCode = rand.Intn(100)
	var dummyDataTemplate1 string
	var dummyDataTemplate2 int
	var dummyDataReceivers = []dataReceiver{
		{
			beginStatusCode: dummyStatusCode - rand.Intn(10),
			endStatusCode:   dummyStatusCode + 1 + rand.Intn(10),
			dataTemplate:    &dummyDataTemplate1,
		},
		{
			beginStatusCode: dummyStatusCode - rand.Intn(10),
			endStatusCode:   dummyStatusCode + 1 + rand.Intn(10),
			dataTemplate:    &dummyDataTemplate2,
		},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, &dummyDataTemplate2, i)
		return false
	})

	// SUT + act
	var result = getDataTemplate(
		dummySession,
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Equal(t, &dummyDataTemplate2, result)
}

func TestParseResponse_ReadError(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyBody = io.NopCloser(bytes.NewBufferString("some body"))
	var dummyBytes = []byte("some bytes")
	var dummyError = errors.New("some error")
	var dummyDataTemplate string

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(io.ReadAll, 1, func(r io.Reader) ([]byte, error) {
		assert.Equal(t, dummyBody, r)
		return dummyBytes, dummyError
	})

	// SUT + act
	var err = parseResponse(
		dummySession,
		dummyBody,
		&dummyDataTemplate,
	)

	// assert
	assert.Zero(t, dummyDataTemplate)
	assert.Equal(t, dummyError, err)
}

func TestParseResponse_JSONError(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyBody = io.NopCloser(bytes.NewBufferString("some body"))
	var dummyBytes = []byte("some bytes")
	var dummyError = errors.New("some error")
	var dummyDataTemplate string
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(io.ReadAll, 1, func(r io.Reader) ([]byte, error) {
		assert.Equal(t, dummyBody, r)
		return dummyBytes, nil
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, string(dummyBytes), value)
		return dummyError
	})
	m.ExpectFunc(logWebcallResponse, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Body", category)
		assert.Equal(t, "UnmarshalError", subcategory)
		assert.Equal(t, "%+v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageResponseInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	})

	// SUT + act
	var err = parseResponse(
		dummySession,
		dummyBody,
		&dummyDataTemplate,
	)

	// assert
	assert.Zero(t, dummyDataTemplate)
	assert.Equal(t, dummyAppError, err)
}

func TestParseResponse_HappyPath(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyBody = io.NopCloser(bytes.NewBufferString("some body"))
	var dummyData = "some data"
	var dummyBytes = []byte("\"" + dummyData + "\"")
	var dummyDataTemplate string

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(io.ReadAll, 1, func(r io.Reader) ([]byte, error) {
		assert.Equal(t, dummyBody, r)
		return dummyBytes, nil
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, string(dummyBytes), value)
		(*(dataTemplate).(*string)) = dummyData
		return nil
	})

	// SUT + act
	var err = parseResponse(
		dummySession,
		dummyBody,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyData, dummyDataTemplate)
	assert.NoError(t, err)
}

func TestWebRequestProcess_NilWebRequest(t *testing.T) {
	// arrange
	var dummyAppError = &appError{Message: "some error message"}

	// SUT
	var sut *webRequest

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)
	assert.Empty(t, header)
	assert.Equal(t, dummyAppError, err)
}

func TestWebRequestProcess_NilWebRequestSession(t *testing.T) {
	// arrange
	var dummyAppError = &appError{Message: "some error message"}

	// SUT
	var sut = &webRequest{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)
	assert.Empty(t, header)
	assert.Equal(t, dummyAppError, err)
}

func TestWebRequestProcess_Error_NilObject(t *testing.T) {
	// arrange
	var dummyResponseObject *http.Response
	var dummyResponseError = errors.New("some error")

	// SUT
	var sut = &webRequest{
		session: &session{id: uuid.New()},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(doRequestProcessing, 1, func(webRequest *webRequest) (*http.Response, error) {
		assert.Equal(t, sut, webRequest)
		return dummyResponseObject, dummyResponseError
	})

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)
	assert.Empty(t, header)
	assert.Equal(t, dummyResponseError, err)
}

func TestWebRequestProcess_Error_ValidObject(t *testing.T) {
	// arrange
	var dummyStatusCode = rand.Int()
	var dummyHeader = map[string][]string{
		"foo":  {"bar"},
		"test": {"123", "456", "789"},
	}
	var dummyResponseObject = &http.Response{
		StatusCode: dummyStatusCode,
		Header:     dummyHeader,
	}
	var dummyResponseError = errors.New("some error")

	// SUT
	var sut = &webRequest{
		session: &session{id: uuid.New()},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(doRequestProcessing, 1, func(webRequest *webRequest) (*http.Response, error) {
		assert.Equal(t, sut, webRequest)
		return dummyResponseObject, dummyResponseError
	})

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, dummyStatusCode, result)
	assert.Equal(t, http.Header(dummyHeader), header)
	assert.Equal(t, dummyResponseError, err)
}

func TestWebRequestProcess_Success_NilObject(t *testing.T) {
	// arrange
	var dummyResponseObject *http.Response
	var dummyResponseError error

	// SUT
	var sut = &webRequest{
		session: &session{id: uuid.New()},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(doRequestProcessing, 1, func(webRequest *webRequest) (*http.Response, error) {
		assert.Equal(t, sut, webRequest)
		return dummyResponseObject, dummyResponseError
	})

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Zero(t, result)
	assert.Empty(t, header)
	assert.NoError(t, err)
}

func TestWebRequestProcess_Success_ValidObject(t *testing.T) {
	// arrange
	var dummyStatusCode = rand.Int()
	var dummyHeader = map[string][]string{
		"foo":  {"bar"},
		"test": {"123", "456", "789"},
	}
	var dummyBody = io.NopCloser(bytes.NewBufferString("some body"))
	var dummyResponseObject = &http.Response{
		StatusCode: dummyStatusCode,
		Header:     dummyHeader,
		Body:       dummyBody,
	}
	var dummyResponseError error
	var dummyParseError = errors.New("some parse error")
	var dummyDataTemplate string
	var dummyData = "some data"
	var dummySession = &session{id: uuid.New()}
	var dummyDataReceivers = []dataReceiver{
		{0, 999, &dummyDataTemplate},
	}

	// SUT
	var sut = &webRequest{
		session:       dummySession,
		dataReceivers: dummyDataReceivers,
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(doRequestProcessing, 1, func(webRequest *webRequest) (*http.Response, error) {
		assert.Equal(t, sut, webRequest)
		return dummyResponseObject, dummyResponseError
	})
	m.ExpectFunc(getDataTemplate, 1, func(session *session, statusCode int, dataReceivers []dataReceiver) interface{} {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStatusCode, statusCode)
		assert.Equal(t, dummyDataReceivers, dataReceivers)
		return &dummyDataTemplate
	})
	m.ExpectFunc(parseResponse, 1, func(session *session, body io.ReadCloser, dataTemplate interface{}) error {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyBody, body)
		(*(dataTemplate).(*string)) = dummyData
		return dummyParseError
	})

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, dummyData, dummyDataTemplate)
	assert.Equal(t, dummyStatusCode, result)
	assert.Equal(t, http.Header(dummyHeader), header)
	assert.Equal(t, dummyParseError, err)
}
