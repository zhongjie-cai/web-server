package webserver

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetClientForRequest_SendClientCert(t *testing.T) {
	// arrange
	var dummyHTTPClient1 = &http.Client{Timeout: time.Duration(rand.Int())}
	var dummyHTTPClient2 = &http.Client{Timeout: time.Duration(rand.Int())}

	// stub
	httpClientWithCert = dummyHTTPClient1
	httpClientNoCert = dummyHTTPClient2

	// mock
	createMock(t)

	// SUT + act
	var result = getClientForRequest(true)

	// assert
	assert.Equal(t, dummyHTTPClient1, result)

	// verify
	verifyAll(t)
}

func TestGetClientForRequest_NoSendClientCert(t *testing.T) {
	// arrange
	var dummyHTTPClient1 = &http.Client{Timeout: time.Duration(rand.Int())}
	var dummyHTTPClient2 = &http.Client{Timeout: time.Duration(rand.Int())}

	// stub
	httpClientWithCert = dummyHTTPClient1
	httpClientNoCert = dummyHTTPClient2

	// mock
	createMock(t)

	// SUT + act
	var result = getClientForRequest(false)

	// assert
	assert.Equal(t, dummyHTTPClient2, result)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	clientDoFuncExpected = 1
	clientDoFunc = func(client *http.Client, request *http.Request) (*http.Response, error) {
		clientDoFuncCalled++
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, dummyResponseError
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	clientDoFuncExpected = 2
	clientDoFunc = func(client *http.Client, request *http.Request) (*http.Response, error) {
		clientDoFuncCalled++
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		if clientDoFuncCalled == 1 {
			return dummyResponseObject, dummyResponseError
		} else if clientDoFuncCalled == 2 {
			return dummyResponseObject, nil
		}
		return nil, nil
	}
	timeSleepExpected = 1
	timeSleep = func(d time.Duration) {
		timeSleepCalled++
		assert.Equal(t, dummyRetryDelay, d)
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	clientDoFuncExpected = 3
	clientDoFunc = func(client *http.Client, request *http.Request) (*http.Response, error) {
		clientDoFuncCalled++
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, dummyResponseError
	}
	timeSleepExpected = 2
	timeSleep = func(d time.Duration) {
		timeSleepCalled++
		assert.Equal(t, dummyRetryDelay, d)
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	clientDoFuncExpected = 1
	clientDoFunc = func(client *http.Client, request *http.Request) (*http.Response, error) {
		clientDoFuncCalled++
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, nil
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	clientDoFuncExpected = 1
	clientDoFunc = func(client *http.Client, request *http.Request) (*http.Response, error) {
		clientDoFuncCalled++
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, nil
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	clientDoFuncExpected = 2
	clientDoFunc = func(client *http.Client, request *http.Request) (*http.Response, error) {
		clientDoFuncCalled++
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		if clientDoFuncCalled == 1 {
			return dummyResponseObject1, nil
		} else if clientDoFuncCalled == 2 {
			return dummyResponseObject2, nil
		}
		return nil, nil
	}
	timeSleepExpected = 1
	timeSleep = func(d time.Duration) {
		timeSleepCalled++
		assert.Equal(t, dummyRetryDelay, d)
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	clientDoFuncExpected = 3
	clientDoFunc = func(client *http.Client, request *http.Request) (*http.Response, error) {
		clientDoFuncCalled++
		assert.Equal(t, dummyClient, client)
		assert.Equal(t, dummyRequestObject, request)
		return dummyResponseObject, nil
	}
	timeSleepExpected = 2
	timeSleep = func(d time.Duration) {
		timeSleepCalled++
		assert.Equal(t, dummyRetryDelay, d)
	}

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

	// verify
	verifyAll(t)
}

func TestGetHTTPTransport_NoClientCert(t *testing.T) {
	// arrange
	var dummySkipServerCertVerification = rand.Intn(100) < 50
	var dummyClientCert *tls.Certificate
	var dummyRoundTripper = &http.Transport{}
	var dummyRoundTripperWrapperExpected int
	var dummyRoundTripperWrapperCalled int
	var dummyRoundTripperWrapper func(http.RoundTripper) http.RoundTripper

	// mock
	createMock(t)

	// expect
	dummyRoundTripperWrapperExpected = 1
	dummyRoundTripperWrapper = func(original http.RoundTripper) http.RoundTripper {
		dummyRoundTripperWrapperCalled++
		assert.Equal(t, http.DefaultTransport, original)
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

	// verify
	verifyAll(t)
	assert.Equal(t, dummyRoundTripperWrapperExpected, dummyRoundTripperWrapperCalled, "Unexpected number of calls to method dummyRoundTripperWrapper")
}

func TestGetHTTPTransport_WithClientCert(t *testing.T) {
	// arrange
	var dummySkipServerCertVerification = rand.Intn(100) < 50
	var dummyClientCert = &tls.Certificate{}
	var dummyRoundTripper = &http.Transport{}
	var dummyRoundTripperWrapperExpected int
	var dummyRoundTripperWrapperCalled int
	var dummyRoundTripperWrapper func(http.RoundTripper) http.RoundTripper

	// mock
	createMock(t)

	// expect
	dummyRoundTripperWrapperExpected = 1
	dummyRoundTripperWrapper = func(original http.RoundTripper) http.RoundTripper {
		dummyRoundTripperWrapperCalled++
		assert.NotEqual(t, http.DefaultTransport, original)
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

	// verify
	verifyAll(t)
	assert.Equal(t, dummyRoundTripperWrapperExpected, dummyRoundTripperWrapperCalled, "Unexpected number of calls to method dummyRoundTripperWrapper")
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
	createMock(t)

	// expect
	getHTTPTransportFuncExpected = 2
	getHTTPTransportFunc = func(skipServerCertVerification bool, clientCertificate *tls.Certificate, roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper) http.RoundTripper {
		getHTTPTransportFuncCalled++
		assert.Equal(t, dummySkipServerCertVerification, skipServerCertVerification)
		if getHTTPTransportFuncCalled == 1 {
			assert.Equal(t, dummyClientCert, clientCertificate)
			return dummyHTTPTransport1
		} else if getHTTPTransportFuncCalled == 2 {
			assert.Nil(t, clientCertificate)
			return dummyHTTPTransport2
		}
		functionPointerEquals(t, dummyRoundTripperWrapper, roundTripperWrapper)
		return nil
	}

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

	// verify
	verifyAll(t)
}

func TestWebREquestAddQuery_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue1 = "some value 1"
	var dummyValue2 = "some value 2"
	var dummyValue3 = "some value 3"

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestWebREquestAddHeader_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue1 = "some value 1"
	var dummyValue2 = "some value 2"
	var dummyValue3 = "some value 3"

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestCreateQueryString_NilQuery(t *testing.T) {
	// arrange
	var dummyQuery map[string][]string
	var dummyResult = "some result"

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		assert.Empty(t, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	}

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestCreateQueryString_EmptyQuery(t *testing.T) {
	// arrange
	var dummyQuery = map[string][]string{}
	var dummyResult = "some result"

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		assert.Empty(t, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	}

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestCreateQueryString_EmptyQueryName(t *testing.T) {
	// arrange
	var dummyQuery = map[string][]string{
		"": {"empty 1", "empty 2"},
	}
	var dummyResult = "some result"

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		assert.Empty(t, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	}

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestCreateQueryString_EmptyQueryValues(t *testing.T) {
	// arrange
	var dummyQuery = map[string][]string{
		"":          {"empty 1", "empty 2"},
		"some name": {},
	}
	var dummyResult = "some result"

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		assert.Empty(t, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	}

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestCreateQueryString_HappyPath(t *testing.T) {
	// arrange
	var dummyNames = []string{
		"some name 1",
		"some name 2",
		"some name 3",
	}
	var dummyValues = [][]string{
		{"some value 1-1", "some value 1-2", "some value 1-3"},
		{"some value 2-1", "some value 2-2", "some value 2-3"},
		{"some value 3-1", "some value 3-2", "some value 3-3"},
	}
	var dummyQuery = map[string][]string{
		"":          {"empty 1", "empty 2"},
		"some name": {},
	}
	var dummyQueryStrings = []string{}
	var dummyResult = "some joined query"

	// stub
	for i := 0; i < len(dummyNames); i++ {
		dummyQuery[dummyNames[i]] = dummyValues[i]
		for j := 0; j < len(dummyValues[i]); j++ {
			dummyQueryStrings = append(
				dummyQueryStrings,
				url.QueryEscape(dummyNames[i])+"="+url.QueryEscape(dummyValues[i][j]),
			)
		}
	}

	// mock
	createMock(t)

	// expect
	urlQueryEscapeExpected = 18
	urlQueryEscape = func(s string) string {
		urlQueryEscapeCalled++
		return url.QueryEscape(s)
	}
	fmtSprintfExpected = 9
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		return fmt.Sprintf(format, a...)
	}
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		assert.ElementsMatch(t, dummyQueryStrings, a)
		assert.Equal(t, "&", sep)
		return dummyResult
	}

	// SUT + act
	var result = createQueryString(
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestGenerateRequestURL_NilQuery(t *testing.T) {
	// arrange
	var dummyBaseURL = "some base URL"
	var dummyQuery map[string][]string = nil

	// mock
	createMock(t)

	// SUT + act
	var result = generateRequestURL(
		dummyBaseURL,
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyBaseURL, result)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	createQueryStringFuncExpected = 1
	createQueryStringFunc = func(query map[string][]string) string {
		createQueryStringFuncCalled++
		assert.Equal(t, dummyQuery, query)
		return dummyQueryString
	}

	// SUT + act
	var result = generateRequestURL(
		dummyBaseURL,
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyBaseURL, result)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	createQueryStringFuncExpected = 1
	createQueryStringFunc = func(query map[string][]string) string {
		createQueryStringFuncCalled++
		assert.Equal(t, dummyQuery, query)
		return dummyQueryString
	}
	fmtSprintfExpected = 1
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		return fmt.Sprintf(format, a...)
	}

	// SUT + act
	var result = generateRequestURL(
		dummyBaseURL,
		dummyQuery,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

type dummyCustomizationWrapRequest struct {
	dummyCustomization
	wrapRequest func(Session, *http.Request) *http.Request
}

func (customization *dummyCustomizationWrapRequest) WrapRequest(session Session, httpRequest *http.Request) *http.Request {
	if customization.wrapRequest != nil {
		return customization.wrapRequest(session, httpRequest)
	}
	assert.Fail(customization.t, "Unexpected call to WrapRequest")
	return nil
}

func TestCreateHTTPRequest_NilWebRequest(t *testing.T) {
	// arrange
	var dummyWebRequest *webRequest
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// SUT + act
	var result, err = createHTTPRequest(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
}

func TestCreateHTTPRequest_NilWebRequestSession(t *testing.T) {
	// arrange
	var dummyWebRequest = &webRequest{}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// SUT + act
	var result, err = createHTTPRequest(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

	// expect
	generateRequestURLFuncExpected = 1
	generateRequestURLFunc = func(baseURL string, query map[string][]string) string {
		generateRequestURLFuncCalled++
		assert.Equal(t, dummyURL, baseURL)
		assert.Equal(t, dummyQuery, query)
		return dummyRequestURL
	}
	stringsNewReaderExpected = 1
	stringsNewReader = func(s string) *strings.Reader {
		stringsNewReaderCalled++
		return strings.NewReader(s)
	}
	httpNewRequestExpected = 1
	httpNewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		httpNewRequestCalled++
		assert.Equal(t, dummyMethod, method)
		assert.Equal(t, dummyRequestURL, url)
		assert.NotNil(t, body)
		return dummyRequest, dummyError
	}

	// SUT + act
	var result, err = createHTTPRequest(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyError, err)

	// verify
	verifyAll(t)
}

func TestCreateHTTPRequest_Success(t *testing.T) {
	// arrange
	var dummyCustomizationWrapRequest = &dummyCustomizationWrapRequest{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationWrapRequest,
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
	var customizationWrapRequestExpected int
	var customizationWrapRequestCalled int
	var dummyCustomized = &http.Request{
		RequestURI: "def",
	}

	// mock
	createMock(t)

	// expect
	generateRequestURLFuncExpected = 1
	generateRequestURLFunc = func(baseURL string, query map[string][]string) string {
		generateRequestURLFuncCalled++
		assert.Equal(t, dummyURL, baseURL)
		assert.Equal(t, dummyQuery, query)
		return dummyRequestURL
	}
	stringsNewReaderExpected = 1
	stringsNewReader = func(s string) *strings.Reader {
		stringsNewReaderCalled++
		return strings.NewReader(s)
	}
	httpNewRequestExpected = 1
	httpNewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		httpNewRequestCalled++
		assert.Equal(t, dummyMethod, method)
		assert.Equal(t, dummyRequestURL, url)
		assert.NotNil(t, body)
		return dummyRequest, nil
	}
	logWebcallStartFuncExpected = 1
	logWebcallStartFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallStartFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethod, category)
		assert.Equal(t, dummyURL, subcategory)
		assert.Equal(t, dummyRequestURL, messageFormat)
		assert.Empty(t, parameters)
	}
	logWebcallRequestFuncExpected = 2
	logWebcallRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Zero(t, subcategory)
		assert.Empty(t, parameters)
		if logWebcallRequestFuncCalled == 1 {
			assert.Equal(t, "Payload", category)
			assert.Equal(t, dummyPayload, messageFormat)
		} else if logWebcallRequestFuncCalled == 2 {
			assert.Equal(t, "Header", category)
			assert.Equal(t, dummyHeaderContent, messageFormat)
		}
	}
	marshalIgnoreErrorFuncExpected = 1
	marshalIgnoreErrorFunc = func(v interface{}) string {
		marshalIgnoreErrorFuncCalled++
		return dummyHeaderContent
	}
	customizationWrapRequestExpected = 1
	dummyCustomizationWrapRequest.wrapRequest = func(session Session, httpRequest *http.Request) *http.Request {
		customizationWrapRequestCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRequest, httpRequest)
		return dummyCustomized
	}

	// SUT + act
	var result, err = createHTTPRequest(
		dummyWebRequest,
	)

	// assert
	assert.Equal(t, dummyCustomized, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationWrapRequestExpected, customizationWrapRequestCalled, "Unexpected number of calls to method customization.WrapRequest")
}

func TestLogErrorResponse(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyError = errors.New("some error")
	var dummyStartTime = time.Now()
	var dummyTimeSince = time.Duration(rand.Intn(1000))

	// mock
	createMock(t)

	// expect
	logWebcallResponseFuncExpected = 1
	logWebcallResponseFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Message", category)
		assert.Zero(t, subcategory)
		assert.Equal(t, "%+v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	}
	timeSinceExpected = 1
	timeSince = func(ts time.Time) time.Duration {
		timeSinceCalled++
		assert.Equal(t, dummyStartTime, ts)
		return dummyTimeSince
	}
	logWebcallFinishFuncExpected = 1
	logWebcallFinishFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallFinishFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Error", category)
		assert.Zero(t, subcategory)
		assert.Equal(t, "%s", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyTimeSince, parameters[0])
	}

	// SUT + act
	logErrorResponse(
		dummySession,
		dummyError,
		dummyStartTime,
	)

	// verify
	verifyAll(t)
}

func TestLogSuccessResponse_NilResponse(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyResponse *http.Response
	var dummyStartTime = time.Now()

	// mock
	createMock(t)

	// SUT + act
	logSuccessResponse(
		dummySession,
		dummyResponse,
		dummyStartTime,
	)

	// verify
	verifyAll(t)
}

func TestLogSuccessResponse_ValidResponse(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyStatus = "some status"
	var dummyStatusCode = rand.Intn(1000)
	var dummyBody = ioutil.NopCloser(bytes.NewBufferString("some body"))
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
	var dummyNewBody = ioutil.NopCloser(bytes.NewBufferString("some new body"))
	var dummyStartTime = time.Now()
	var dummyHeaderContent = "some header content"
	var dummyTimeSince = time.Duration(rand.Intn(1000))

	// mock
	createMock(t)

	// expect
	ioutilReadAllExpected = 1
	ioutilReadAll = func(r io.Reader) ([]byte, error) {
		ioutilReadAllCalled++
		assert.Equal(t, dummyBody, r)
		return dummyResponseBytes, dummyError
	}
	bytesNewBufferExpected = 1
	bytesNewBuffer = func(buf []byte) *bytes.Buffer {
		bytesNewBufferCalled++
		assert.Equal(t, dummyResponseBytes, buf)
		return dummyBuffer
	}
	ioutilNopCloserExpected = 1
	ioutilNopCloser = func(r io.Reader) io.ReadCloser {
		ioutilNopCloserCalled++
		assert.Equal(t, dummyBuffer, r)
		return dummyNewBody
	}
	httpStatusTextExpected = 1
	httpStatusText = func(code int) string {
		httpStatusTextCalled++
		assert.Equal(t, dummyStatusCode, code)
		return dummyStatus
	}
	strconvItoaExpected = 1
	strconvItoa = func(i int) string {
		strconvItoaCalled++
		assert.Equal(t, dummyStatusCode, i)
		return strconv.Itoa(i)
	}
	logWebcallResponseFuncExpected = 2
	logWebcallResponseFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Zero(t, subcategory)
		assert.Empty(t, parameters)
		if logWebcallResponseFuncCalled == 1 {
			assert.Equal(t, "Header", category)
			assert.Equal(t, dummyHeaderContent, messageFormat)
		} else if logWebcallResponseFuncCalled == 2 {
			assert.Equal(t, "Body", category)
			assert.Equal(t, dummyResponseBody, messageFormat)
		}
	}
	marshalIgnoreErrorFuncExpected = 1
	marshalIgnoreErrorFunc = func(v interface{}) string {
		marshalIgnoreErrorFuncCalled++
		assert.Equal(t, dummyHeader, v)
		return dummyHeaderContent
	}
	timeSinceExpected = 1
	timeSince = func(ts time.Time) time.Duration {
		timeSinceCalled++
		assert.Equal(t, dummyStartTime, ts)
		return dummyTimeSince
	}
	logWebcallFinishFuncExpected = 1
	logWebcallFinishFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallFinishFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStatus, category)
		assert.Equal(t, strconv.Itoa(dummyStatusCode), subcategory)
		assert.Equal(t, "%s", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyTimeSince, parameters[0])
	}

	// SUT + act
	logSuccessResponse(
		dummySession,
		dummyResponse,
		dummyStartTime,
	)

	// assert
	assert.Equal(t, dummyNewBody, dummyResponse.Body)

	// verify
	verifyAll(t)
}

func TestDoRequestProcessing_NilWebRequest(t *testing.T) {
	// arrange
	var dummyWebRequest *webRequest
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
}

func TestDoRequestProcessing_NilWebRequestSession(t *testing.T) {
	// arrange
	var dummyWebRequest = &webRequest{}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
}

func TestDoRequestProcessing_RequestError(t *testing.T) {
	// arrange
	var dummyWebRequest = &webRequest{
		session: &session{id: uuid.New()},
	}
	var dummyRequestObject *http.Request
	var dummyRequestError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	createHTTPRequestFuncExpected = 1
	createHTTPRequestFunc = func(webRequest *webRequest) (*http.Request, error) {
		createHTTPRequestFuncCalled++
		assert.Equal(t, dummyWebRequest, webRequest)
		return dummyRequestObject, dummyRequestError
	}

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyRequestError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	createHTTPRequestFuncExpected = 1
	createHTTPRequestFunc = func(webRequest *webRequest) (*http.Request, error) {
		createHTTPRequestFuncCalled++
		assert.Equal(t, dummyWebRequest, webRequest)
		return dummyRequestObject, nil
	}
	getClientForRequestFuncExpected = 1
	getClientForRequestFunc = func(sendClientCert bool) *http.Client {
		getClientForRequestFuncCalled++
		assert.Equal(t, dummySendClientCert, sendClientCert)
		return dummyHTTPClient
	}
	getTimeNowUTCFuncExpected = 1
	getTimeNowUTCFunc = func() time.Time {
		getTimeNowUTCFuncCalled++
		return dummyStartTime
	}
	clientDoWithRetryFuncExpected = 1
	clientDoWithRetryFunc = func(client *http.Client, request *http.Request, connRetry int, httpRetry map[int]int, retryDelay time.Duration) (*http.Response, error) {
		clientDoWithRetryFuncCalled++
		assert.Equal(t, dummyHTTPClient, client)
		assert.Equal(t, dummyRequestObject, request)
		assert.Equal(t, dummyConnRetry, connRetry)
		assert.Equal(t, dummyHTTPRetry, httpRetry)
		assert.Equal(t, dummyRetryDelay, retryDelay)
		return dummyResponseObject, dummyResponseError
	}
	logErrorResponseFuncExpected = 1
	logErrorResponseFunc = func(session *session, responseError error, startTime time.Time) {
		logErrorResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseError, responseError)
		assert.Equal(t, dummyStartTime, startTime)
	}

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.Equal(t, dummyResponseError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	createHTTPRequestFuncExpected = 1
	createHTTPRequestFunc = func(webRequest *webRequest) (*http.Request, error) {
		createHTTPRequestFuncCalled++
		assert.Equal(t, dummyWebRequest, webRequest)
		return dummyRequestObject, nil
	}
	getClientForRequestFuncExpected = 1
	getClientForRequestFunc = func(sendClientCert bool) *http.Client {
		getClientForRequestFuncCalled++
		assert.Equal(t, dummySendClientCert, sendClientCert)
		return dummyHTTPClient
	}
	getTimeNowUTCFuncExpected = 1
	getTimeNowUTCFunc = func() time.Time {
		getTimeNowUTCFuncCalled++
		return dummyStartTime
	}
	clientDoWithRetryFuncExpected = 1
	clientDoWithRetryFunc = func(client *http.Client, request *http.Request, connRetry int, httpRetry map[int]int, retryDelay time.Duration) (*http.Response, error) {
		clientDoWithRetryFuncCalled++
		assert.Equal(t, dummyHTTPClient, client)
		assert.Equal(t, dummyRequestObject, request)
		assert.Equal(t, dummyConnRetry, connRetry)
		assert.Equal(t, dummyHTTPRetry, httpRetry)
		assert.Equal(t, dummyRetryDelay, retryDelay)
		return dummyResponseObject, nil
	}
	logSuccessResponseFuncExpected = 1
	logSuccessResponseFunc = func(session *session, response *http.Response, startTime time.Time) {
		logSuccessResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, response)
		assert.Equal(t, dummyStartTime, startTime)
	}

	// SUT + act
	var result, err = doRequestProcessing(
		dummyWebRequest,
	)

	// assert
	assert.Equal(t, dummyResponseObject, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestGetDataTemplate_EmptyDataReceivers(t *testing.T) {
	// arrange
	var dummyStatusCode = rand.Intn(100)
	var dummyDataReceivers = []dataReceiver{}

	// mock
	createMock(t)

	// SUT + act
	var result = getDataTemplate(
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
}

func TestGetDataTemplate_NoMatch(t *testing.T) {
	// arrange
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
	createMock(t)

	// SUT + act
	var result = getDataTemplate(
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
}

func TestGetDataTemplate_SingleMatch(t *testing.T) {
	// arrange
	var dummyStatusCode = rand.Intn(100)
	var dummyDataTemplate1 string
	var dummyDataTemplate2 int
	var dummyDataReceivers = []dataReceiver{
		{
			beginStatusCode: dummyStatusCode - rand.Intn(10),
			endStatusCode:   dummyStatusCode + rand.Intn(10),
			dataTemplate:    &dummyDataTemplate1,
		},
		{
			beginStatusCode: rand.Intn(100) + 100,
			endStatusCode:   rand.Intn(100) + 200,
			dataTemplate:    &dummyDataTemplate2,
		},
	}

	// mock
	createMock(t)

	// SUT + act
	var result = getDataTemplate(
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Equal(t, &dummyDataTemplate1, result)

	// verify
	verifyAll(t)
}

func TestGetDataTemplate_OverlapMatch(t *testing.T) {
	// arrange
	var dummyStatusCode = rand.Intn(100)
	var dummyDataTemplate1 string
	var dummyDataTemplate2 int
	var dummyDataReceivers = []dataReceiver{
		{
			beginStatusCode: dummyStatusCode - rand.Intn(10),
			endStatusCode:   dummyStatusCode + rand.Intn(10),
			dataTemplate:    &dummyDataTemplate1,
		},
		{
			beginStatusCode: dummyStatusCode - rand.Intn(10),
			endStatusCode:   dummyStatusCode + rand.Intn(10),
			dataTemplate:    &dummyDataTemplate2,
		},
	}

	// mock
	createMock(t)

	// SUT + act
	var result = getDataTemplate(
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Equal(t, &dummyDataTemplate2, result)

	// verify
	verifyAll(t)
}

func TestParseResponse_ReadError(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyBody = ioutil.NopCloser(bytes.NewBufferString("some body"))
	var dummyBytes = []byte("some bytes")
	var dummyError = errors.New("some error")
	var dummyDataTemplate string

	// mock
	createMock(t)

	// expect
	ioutilReadAllExpected = 1
	ioutilReadAll = func(r io.Reader) ([]byte, error) {
		ioutilReadAllCalled++
		assert.Equal(t, dummyBody, r)
		return dummyBytes, dummyError
	}

	// SUT + act
	var err = parseResponse(
		dummySession,
		dummyBody,
		&dummyDataTemplate,
	)

	// assert
	assert.Zero(t, dummyDataTemplate)
	assert.Equal(t, dummyError, err)

	// verify
	verifyAll(t)
}

func TestParseResponse_JSONError(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyBody = ioutil.NopCloser(bytes.NewBufferString("some body"))
	var dummyBytes = []byte("some bytes")
	var dummyError = errors.New("some error")
	var dummyDataTemplate string
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	ioutilReadAllExpected = 1
	ioutilReadAll = func(r io.Reader) ([]byte, error) {
		ioutilReadAllCalled++
		assert.Equal(t, dummyBody, r)
		return dummyBytes, nil
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, string(dummyBytes), value)
		return dummyError
	}
	logWebcallResponseFuncExpected = 1
	logWebcallResponseFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Body", category)
		assert.Equal(t, "UnmarshalError", subcategory)
		assert.Equal(t, "%+v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageResponseInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	}

	// SUT + act
	var err = parseResponse(
		dummySession,
		dummyBody,
		&dummyDataTemplate,
	)

	// assert
	assert.Zero(t, dummyDataTemplate)
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
}

func TestParseResponse_HappyPath(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyBody = ioutil.NopCloser(bytes.NewBufferString("some body"))
	var dummyData = "some data"
	var dummyBytes = []byte("\"" + dummyData + "\"")
	var dummyDataTemplate string

	// mock
	createMock(t)

	// expect
	ioutilReadAllExpected = 1
	ioutilReadAll = func(r io.Reader) ([]byte, error) {
		ioutilReadAllCalled++
		assert.Equal(t, dummyBody, r)
		return dummyBytes, nil
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, string(dummyBytes), value)
		(*(dataTemplate).(*string)) = dummyData
		return nil
	}

	// SUT + act
	var err = parseResponse(
		dummySession,
		dummyBody,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyData, dummyDataTemplate)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestWebRequestProcess_NilWebRequest(t *testing.T) {
	// arrange
	var dummyAppError = &appError{Message: "some error message"}

	// SUT
	var sut *webRequest

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)
	assert.Empty(t, header)
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
}

func TestWebRequestProcess_NilWebRequestSession(t *testing.T) {
	// arrange
	var dummyAppError = &appError{Message: "some error message"}

	// SUT
	var sut = &webRequest{}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageWebRequestNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)
	assert.Empty(t, header)
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	doRequestProcessingFuncExpected = 1
	doRequestProcessingFunc = func(webRequest *webRequest) (*http.Response, error) {
		doRequestProcessingFuncCalled++
		assert.Equal(t, sut, webRequest)
		return dummyResponseObject, dummyResponseError
	}

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)
	assert.Empty(t, header)
	assert.Equal(t, dummyResponseError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	doRequestProcessingFuncExpected = 1
	doRequestProcessingFunc = func(webRequest *webRequest) (*http.Response, error) {
		doRequestProcessingFuncCalled++
		assert.Equal(t, sut, webRequest)
		return dummyResponseObject, dummyResponseError
	}

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, dummyStatusCode, result)
	assert.Equal(t, http.Header(dummyHeader), header)
	assert.Equal(t, dummyResponseError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	doRequestProcessingFuncExpected = 1
	doRequestProcessingFunc = func(webRequest *webRequest) (*http.Response, error) {
		doRequestProcessingFuncCalled++
		assert.Equal(t, sut, webRequest)
		return dummyResponseObject, dummyResponseError
	}

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Zero(t, result)
	assert.Empty(t, header)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestWebRequestProcess_Success_ValidObject(t *testing.T) {
	// arrange
	var dummyStatusCode = rand.Int()
	var dummyHeader = map[string][]string{
		"foo":  {"bar"},
		"test": {"123", "456", "789"},
	}
	var dummyBody = ioutil.NopCloser(bytes.NewBufferString("some body"))
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
	createMock(t)

	// expect
	doRequestProcessingFuncExpected = 1
	doRequestProcessingFunc = func(webRequest *webRequest) (*http.Response, error) {
		doRequestProcessingFuncCalled++
		assert.Equal(t, sut, webRequest)
		return dummyResponseObject, dummyResponseError
	}
	getDataTemplateFuncExpected = 1
	getDataTemplateFunc = func(statusCode int, dataReceivers []dataReceiver) interface{} {
		getDataTemplateFuncCalled++
		assert.Equal(t, dummyStatusCode, statusCode)
		assert.Equal(t, dummyDataReceivers, dataReceivers)
		return &dummyDataTemplate
	}
	parseResponseFuncExpected = 1
	parseResponseFunc = func(session *session, body io.ReadCloser, dataTemplate interface{}) error {
		parseResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyBody, body)
		(*(dataTemplate).(*string)) = dummyData
		return dummyParseError
	}

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, dummyData, dummyDataTemplate)
	assert.Equal(t, dummyStatusCode, result)
	assert.Equal(t, http.Header(dummyHeader), header)
	assert.Equal(t, dummyParseError, err)

	// verify
	verifyAll(t)
}
