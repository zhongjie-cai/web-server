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
	"github.com/zhongjie-cai/gomocker/v2"
)

func TestGetClientForRequest_UseCustomClient(t *testing.T) {
	// arrange
	var dummyHttpClient0 = &http.Client{Timeout: time.Duration(rand.Int())}
	var dummyHTTPClient1 = &http.Client{Timeout: time.Duration(rand.Int())}
	var dummyHTTPClient2 = &http.Client{Timeout: time.Duration(rand.Int())}

	// stub
	httpClientCustom = dummyHttpClient0
	httpClientWithCert = dummyHTTPClient1
	httpClientNoCert = dummyHTTPClient2

	// SUT + act
	var result = getClientForRequest(true)

	// assert
	assert.Equal(t, dummyHttpClient0, result)
}

func TestGetClientForRequest_SendClientCert(t *testing.T) {
	// arrange
	var dummyHTTPClient1 = &http.Client{Timeout: time.Duration(rand.Int())}
	var dummyHTTPClient2 = &http.Client{Timeout: time.Duration(rand.Int())}

	// stub
	httpClientCustom = nil
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
	httpClientCustom = nil
	httpClientWithCert = dummyHTTPClient1
	httpClientNoCert = dummyHTTPClient2

	// SUT + act
	var result = getClientForRequest(false)

	// assert
	assert.Equal(t, dummyHTTPClient2, result)
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
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject, dummyResponseError).Once()

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
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject, dummyResponseError).Once()
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject, nil).Once()
	m.Mock(time.Sleep).Expects(dummyRetryDelay).Returns().Once()

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
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject, dummyResponseError).Times(3)
	m.Mock(time.Sleep).Expects(dummyRetryDelay).Returns().Twice()

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
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject, nil).Once()

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
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject, nil).Once()

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
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject1, nil).Once()
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject2, nil).Once()
	m.Mock(time.Sleep).Expects(dummyRetryDelay).Returns().Once()

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
	m.Mock((*http.Client).Do).Expects(dummyClient, dummyRequestObject).Returns(dummyResponseObject, nil).Times(3)
	m.Mock(time.Sleep).Expects(dummyRetryDelay).Returns().Twice()

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
	var dummyRoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: dummySkipServerCertVerification,
		},
		Proxy: http.ProxyFromEnvironment,
	}
	var dummyRoundTripperWrapper = func(rt http.RoundTripper) http.RoundTripper { return nil }

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(dummyRoundTripperWrapper).Expects(gomocker.Matches(func(value *http.Transport) bool {
		return value.TLSClientConfig.InsecureSkipVerify == dummySkipServerCertVerification &&
			len(value.TLSClientConfig.Certificates) == 0 &&
			functionPointerEquals(value.Proxy, http.ProxyFromEnvironment)
	})).Returns(dummyRoundTripper).Once()

	// SUT + act
	var result = getHTTPTransport(
		dummySkipServerCertVerification,
		dummyClientCert,
		dummyRoundTripperWrapper,
	)

	// assert
	assert.Equal(t, dummyRoundTripper, result)
}

func TestGetHTTPTransport_WithClientCert(t *testing.T) {
	// arrange
	var dummySkipServerCertVerification = rand.Intn(100) < 50
	var dummyClientCert = &tls.Certificate{}
	var dummyRoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{*dummyClientCert},
			InsecureSkipVerify: dummySkipServerCertVerification,
		},
		Proxy: http.ProxyFromEnvironment,
	}
	var dummyRoundTripperWrapper = func(rt http.RoundTripper) http.RoundTripper { return nil }

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(dummyRoundTripperWrapper).Expects(gomocker.Matches(func(value *http.Transport) bool {
		return value.TLSClientConfig.InsecureSkipVerify == dummySkipServerCertVerification &&
			len(value.TLSClientConfig.Certificates) == 1 &&
			assert.ObjectsAreEqual(value.TLSClientConfig.Certificates[0], *dummyClientCert) &&
			functionPointerEquals(value.Proxy, http.ProxyFromEnvironment)
	})).Returns(dummyRoundTripper).Once()

	// SUT + act
	var result = getHTTPTransport(
		dummySkipServerCertVerification,
		dummyClientCert,
		dummyRoundTripperWrapper,
	)

	// assert
	assert.Equal(t, dummyRoundTripper, result)
}

func TestInitializeHTTPClients(t *testing.T) {
	// arrange
	var dummyHttpClient = &http.Client{Timeout: time.Duration(rand.Intn(100))}
	var dummyWebcallTimeout = time.Duration(rand.Int())
	var dummySkipServerCertVerification = rand.Intn(100) < 50
	var dummyClientCert = &tls.Certificate{}
	var dummyHTTPTransport1 = &http.Transport{MaxConnsPerHost: rand.Int()}
	var dummyHTTPTransport2 = &http.Transport{MaxConnsPerHost: rand.Int()}
	var dummyRoundTripperWrapper = func(http.RoundTripper) http.RoundTripper { return nil }

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(getHTTPTransport).Expects(dummySkipServerCertVerification, dummyClientCert, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyRoundTripperWrapper, value)
	})).Returns(dummyHTTPTransport1).Once()
	m.Mock(getHTTPTransport).Expects(dummySkipServerCertVerification, nil, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyRoundTripperWrapper, value)
	})).Returns(dummyHTTPTransport2).Once()

	// SUT + act
	initializeHTTPClients(
		dummyHttpClient,
		dummyWebcallTimeout,
		dummySkipServerCertVerification,
		dummyClientCert,
		dummyRoundTripperWrapper,
	)

	// assert
	assert.NotNil(t, httpClientCustom)
	assert.Equal(t, dummyHttpClient, httpClientCustom)
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

func TestWebREquestAddQueries_HappyPath(t *testing.T) {
	// arrange
	var dummyName1 = "some name 1"
	var dummyValue1 = "some value 1"
	var dummyName2 = "some name 2"
	var dummyValue2 = "some value 2"
	var dummyName3 = "some name 3"
	var dummyValue3 = "some value 3"
	var dummyQueries = map[string]string{
		dummyName1: dummyValue1,
		dummyName2: dummyValue2,
		dummyName3: dummyValue3,
	}

	// SUT
	var sut = &webRequest{}

	// act
	var result, ok = sut.AddQueries(
		dummyQueries,
	).(*webRequest)

	// assert
	assert.True(t, ok)
	assert.NotNil(t, result.query)
	assert.Equal(t, 3, len(result.query))
	var values, found = result.query[dummyName1]
	assert.True(t, found)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, dummyValue1, values[0])
	values, found = result.query[dummyName2]
	assert.True(t, found)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, dummyValue2, values[0])
	values, found = result.query[dummyName3]
	assert.True(t, found)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, dummyValue3, values[0])
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

func TestWebREquestAddHeaders_HappyPath(t *testing.T) {
	// arrange
	var dummyName1 = "some name 1"
	var dummyValue1 = "some value 1"
	var dummyName2 = "some name 2"
	var dummyValue2 = "some value 2"
	var dummyName3 = "some name 3"
	var dummyValue3 = "some value 3"
	var dummyHeaders = map[string]string{
		dummyName1: dummyValue1,
		dummyName2: dummyValue2,
		dummyName3: dummyValue3,
	}

	// SUT
	var sut = &webRequest{}

	// act
	var result, ok = sut.AddHeaders(
		dummyHeaders,
	).(*webRequest)

	// assert
	assert.True(t, ok)
	assert.NotNil(t, result.header)
	assert.Equal(t, 3, len(result.header))
	var values, found = result.header[dummyName1]
	assert.True(t, found)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, dummyValue1, values[0])
	values, found = result.header[dummyName2]
	assert.True(t, found)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, dummyValue2, values[0])
	values, found = result.header[dummyName3]
	assert.True(t, found)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, dummyValue3, values[0])
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

func TestWebRequestAnticipate_NoStatusCodes(t *testing.T) {
	// arrange
	var dummyDataTemplate1 string
	var dummyDataTemplate2 int

	// SUT
	var sut = &webRequest{}

	// act
	var result, ok = sut.Anticipate(
		&dummyDataTemplate1,
	).Anticipate(
		&dummyDataTemplate2,
	).(*webRequest)

	// assert
	assert.True(t, ok)
	assert.Equal(t, 2, len(result.dataReceivers))
	assert.Equal(t, &dummyDataTemplate1, result.dataReceivers[0].dataTemplate)
	assert.Equal(t, 0, result.dataReceivers[0].codeRange.Begin)
	assert.Equal(t, 999, result.dataReceivers[0].codeRange.End)
	assert.Equal(t, &dummyDataTemplate2, result.dataReceivers[1].dataTemplate)
	assert.Equal(t, 0, result.dataReceivers[1].codeRange.Begin)
	assert.Equal(t, 999, result.dataReceivers[1].codeRange.End)
}

func TestWebRequestAnticipate_WithStatusCodes(t *testing.T) {
	// arrange
	var dummyStatusCode1 = rand.Int()
	var dummyStatusCode2 = rand.Int()
	var dummyDataTemplate1 string
	var dummyStatusCodeRange1 = rand.Int()
	var dummyStatusCodeRange2 = rand.Int()
	var dummyStatusCodeRange3 = rand.Int()
	var dummyDataTemplate2 int
	var dummyStatusCodeInvalid = "invalid"
	var dummyDataTemplate3 bool

	// SUT
	var sut = &webRequest{}

	// act
	var result, ok = sut.Anticipate(
		&dummyDataTemplate1,
		dummyStatusCode1,
		dummyStatusCode2,
	).Anticipate(
		&dummyDataTemplate2,
		StatusCodeRange{
			Begin: dummyStatusCodeRange1,
			End:   dummyStatusCodeRange2,
		},
		StatusCodeRange{
			Begin: dummyStatusCodeRange3,
		},
	).Anticipate(
		&dummyDataTemplate3,
		dummyStatusCodeInvalid,
	).(*webRequest)

	// assert
	assert.True(t, ok)
	assert.Equal(t, 4, len(result.dataReceivers))
	assert.Equal(t, &dummyDataTemplate1, result.dataReceivers[0].dataTemplate)
	assert.Equal(t, dummyStatusCode1, result.dataReceivers[0].codeRange.Begin)
	assert.Equal(t, dummyStatusCode1+1, result.dataReceivers[0].codeRange.End)
	assert.Equal(t, &dummyDataTemplate1, result.dataReceivers[1].dataTemplate)
	assert.Equal(t, dummyStatusCode2, result.dataReceivers[1].codeRange.Begin)
	assert.Equal(t, dummyStatusCode2+1, result.dataReceivers[1].codeRange.End)
	assert.Equal(t, &dummyDataTemplate2, result.dataReceivers[2].dataTemplate)
	assert.Equal(t, dummyStatusCodeRange1, result.dataReceivers[2].codeRange.Begin)
	assert.Equal(t, dummyStatusCodeRange2, result.dataReceivers[2].codeRange.End)
	assert.Equal(t, &dummyDataTemplate2, result.dataReceivers[3].dataTemplate)
	assert.Equal(t, dummyStatusCodeRange3, result.dataReceivers[3].codeRange.Begin)
	assert.Equal(t, 999, result.dataReceivers[3].codeRange.End)
}

func TestCreateQueryString_NilQuery(t *testing.T) {
	// arrange
	var dummyQuery map[string][]string
	var dummyResult = "some result"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(strings.Join).Expects(gomocker.Anything(), "&").Returns(dummyResult).Once()

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
	m.Mock(strings.Join).Expects(gomocker.Anything(), "&").Returns(dummyResult).Once()

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
	m.Mock(strings.Join).Expects(gomocker.Anything(), "&").Returns(dummyResult).Once()

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
	m.Mock(strings.Join).Expects(gomocker.Anything(), "&").Returns(dummyResult).Once()

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
	m.Mock(createQueryString).Expects(dummyQuery).Returns(dummyQueryString).Once()

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
	m.Mock(createQueryString).Expects(dummyQuery).Returns(dummyQueryString).Once()

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
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageWebRequestNil).Returns(dummyAppError).Once()

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
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageWebRequestNil).Returns(dummyAppError).Once()

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
		{nil, StatusCodeRange{0, 999}},
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
	m.Mock(generateRequestURL).Expects(dummyURL, dummyQuery).Returns(dummyRequestURL).Once()
	m.Mock(strings.NewReader).Expects(dummyPayload).Returns(dummyStingsReader).Once()
	m.Mock(http.NewRequest).Expects(dummyMethod, dummyRequestURL, gomocker.Anything()).Returns(dummyRequest, dummyError).Once()

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
		{nil, StatusCodeRange{0, 999}},
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
	m.Mock(generateRequestURL).Expects(dummyURL, dummyQuery).Returns(dummyRequestURL).Once()
	m.Mock(strings.NewReader).Expects(dummyPayload).Returns(dummyStingsReader).Once()
	m.Mock(http.NewRequest).Expects(dummyMethod, dummyRequestURL, gomocker.Anything()).Returns(dummyRequest, nil).Once()
	m.Mock(logWebcallStart).Expects(dummySession, dummyMethod, dummyURL, "%s", dummyRequestURL).Returns().Once()
	m.Mock(logWebcallRequest).Expects(dummySession, "Payload", "Content", "%s", dummyPayload).Returns().Once()
	m.Mock(logWebcallRequest).Expects(dummySession, "Header", "Content", "%s", dummyHeaderContent).Returns().Once()
	m.Mock(marshalIgnoreError).Expects(gomocker.Anything()).Returns(dummyHeaderContent).Once()
	m.Mock((*DefaultCustomization).WrapRequest).Expects(dummyCustomization, dummySession, dummyRequest).Returns(dummyCustomized).Once()

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
	m.Mock(logWebcallResponse).Expects(dummySession, "Error", "Content", "%+v", dummyError).Returns().Once()
	m.Mock(time.Since).Expects(dummyStartTime).Returns(dummyTimeSince).Once()
	m.Mock(logWebcallFinish).Expects(dummySession, "Error", "-1", "%s", dummyTimeSince).Returns().Once()

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
	m.Mock(io.ReadAll).Expects(dummyBody).Returns(dummyResponseBytes, dummyError).Once()
	m.Mock(bytes.NewBuffer).Expects(dummyResponseBytes).Returns(dummyBuffer).Once()
	m.Mock(io.NopCloser).Expects(dummyBuffer).Returns(dummyNewBody).Once()
	m.Mock(http.StatusText).Expects(dummyStatusCode).Returns(dummyStatus).Once()
	m.Mock(logWebcallResponse).Expects(dummySession, "Header", "Content", "%s", dummyHeaderContent).Returns().Once()
	m.Mock(logWebcallResponse).Expects(dummySession, "Body", "Content", "%s", dummyResponseBody).Returns().Once()
	m.Mock(marshalIgnoreError).Expects(gomocker.Anything()).Returns(dummyHeaderContent).Once()
	m.Mock(time.Since).Expects(dummyStartTime).Returns(dummyTimeSince).Once()
	m.Mock(logWebcallFinish).Expects(dummySession, dummyStatus, strconv.Itoa(dummyStatusCode), "%s", dummyTimeSince).Returns().Once()

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
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageWebRequestNil).Returns(dummyAppError).Once()

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
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageWebRequestNil).Returns(dummyAppError).Once()

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
	m.Mock(createHTTPRequest).Expects(dummyWebRequest).Returns(dummyRequestObject, dummyRequestError).Once()

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
	m.Mock(createHTTPRequest).Expects(dummyWebRequest).Returns(dummyRequestObject, nil).Once()
	m.Mock(getClientForRequest).Expects(dummySendClientCert).Returns(dummyHTTPClient).Once()
	m.Mock(getTimeNowUTC).Expects().Returns(dummyStartTime).Once()
	m.Mock(clientDoWithRetry).Expects(dummyHTTPClient, dummyRequestObject, dummyConnRetry, dummyHTTPRetry, dummyRetryDelay).Returns(dummyResponseObject, dummyResponseError).Once()
	m.Mock(logErrorResponse).Expects(dummySession, dummyResponseError, dummyStartTime).Returns().Once()

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
	m.Mock(createHTTPRequest).Expects(dummyWebRequest).Returns(dummyRequestObject, nil).Once()
	m.Mock(getClientForRequest).Expects(dummySendClientCert).Returns(dummyHTTPClient).Once()
	m.Mock(getTimeNowUTC).Expects().Returns(dummyStartTime).Once()
	m.Mock(clientDoWithRetry).Expects(dummyHTTPClient, dummyRequestObject, dummyConnRetry, dummyHTTPRetry, dummyRetryDelay).Returns(dummyResponseObject, nil).Once()
	m.Mock(logSuccessResponse).Expects(dummySession, dummyResponseObject, dummyStartTime).Returns().Once()

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
	m.Mock(isInterfaceValueNil).Expects(nil).Returns(true).Once()
	m.Mock(logWebcallResponse).Expects(dummySession, "Body", "Receiver", "No data template registered for HTTP status %v", dummyStatusCode).Returns().Once()

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
			codeRange: StatusCodeRange{
				Begin: rand.Intn(100) + 100,
				End:   rand.Intn(100) + 200,
			},
			dataTemplate: &dummyDataTemplate,
		},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(isInterfaceValueNil).Expects(nil).Returns(true).Once()
	m.Mock(logWebcallResponse).Expects(dummySession, "Body", "Receiver", "No data template registered for HTTP status %v", dummyStatusCode).Returns().Once()

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
			codeRange: StatusCodeRange{
				Begin: dummyStatusCode - rand.Intn(10),
				End:   dummyStatusCode + 1 + rand.Intn(10),
			},
			dataTemplate: &dummyDataTemplate1,
		},
		{
			codeRange: StatusCodeRange{
				Begin: rand.Intn(100) + 100,
				End:   rand.Intn(100) + 200,
			},
			dataTemplate: &dummyDataTemplate2,
		},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(isInterfaceValueNil).Expects(&dummyDataTemplate1).Returns(false).Once()

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
			codeRange: StatusCodeRange{
				Begin: dummyStatusCode - rand.Intn(10),
				End:   dummyStatusCode + 1 + rand.Intn(10),
			},
			dataTemplate: &dummyDataTemplate1,
		},
		{
			codeRange: StatusCodeRange{
				Begin: dummyStatusCode - rand.Intn(10),
				End:   dummyStatusCode + 1 + rand.Intn(10),
			},
			dataTemplate: &dummyDataTemplate2,
		},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(isInterfaceValueNil).Expects(&dummyDataTemplate2).Returns(false).Once()

	// SUT + act
	var result = getDataTemplate(
		dummySession,
		dummyStatusCode,
		dummyDataReceivers,
	)

	// assert
	assert.Equal(t, &dummyDataTemplate2, result)
}

func TestParseResponse_NilDataTemplate(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyBody = io.NopCloser(bytes.NewBufferString("some body"))
	var dummyDataTemplate string

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(isInterfaceValueNil).Expects(&dummyDataTemplate).Returns(true).Once()
	m.Mock(logWebcallResponse).Expects(dummySession, "Body", "Skipped", "No data receiver defined for this body").Returns().Once()

	// SUT + act
	var err = parseResponse(
		dummySession,
		dummyBody,
		&dummyDataTemplate,
	)

	// assert
	assert.Zero(t, dummyDataTemplate)
	assert.NoError(t, err)
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
	m.Mock(isInterfaceValueNil).Expects(&dummyDataTemplate).Returns(false).Once()
	m.Mock(io.ReadAll).Expects(dummyBody).Returns(dummyBytes, dummyError).Once()

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
	m.Mock(isInterfaceValueNil).Expects(&dummyDataTemplate).Returns(false).Once()
	m.Mock(io.ReadAll).Expects(dummyBody).Returns(dummyBytes, nil).Once()
	m.Mock(tryUnmarshal).Expects(string(dummyBytes), gomocker.Anything()).Returns(dummyError).Once()
	m.Mock(logWebcallResponse).Expects(dummySession, "Body", "UnmarshalError", "%+v", dummyError).Returns().Once()
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageResponseInvalid, dummyError).Returns(dummyAppError).Once()

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
	m.Mock(isInterfaceValueNil).Expects(&dummyDataTemplate).Returns(false).Once()
	m.Mock(io.ReadAll).Expects(dummyBody).Returns(dummyBytes, nil).Once()
	m.Mock(tryUnmarshal).Expects(string(dummyBytes), gomocker.Anything()).Returns(nil).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *string) { *value = dummyData })).Once()

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
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageWebRequestNil).Returns(dummyAppError).Once()

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
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageWebRequestNil).Returns(dummyAppError).Once()

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
	m.Mock(doRequestProcessing).Expects(sut).Returns(dummyResponseObject, dummyResponseError).Once()

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
	m.Mock(doRequestProcessing).Expects(sut).Returns(dummyResponseObject, dummyResponseError).Once()

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
	m.Mock(doRequestProcessing).Expects(sut).Returns(dummyResponseObject, dummyResponseError).Once()
	m.Mock(logWebcallResponse).Expects(sut.session, "webRequest", "Process", "Nil response object received").Returns().Once()

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
		{&dummyDataTemplate, StatusCodeRange{0, 999}},
	}

	// SUT
	var sut = &webRequest{
		session:       dummySession,
		dataReceivers: dummyDataReceivers,
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(doRequestProcessing).Expects(sut).Returns(dummyResponseObject, dummyResponseError).Once()
	m.Mock(getDataTemplate).Expects(dummySession, dummyStatusCode, dummyDataReceivers).Returns(&dummyDataTemplate).Once()
	m.Mock(parseResponse).Expects(dummySession, dummyBody, gomocker.Anything()).Returns(dummyParseError).SideEffects(
		gomocker.ParamSideEffect(1, 3, func(value *string) { *value = dummyData })).Once()

	// act
	var result, header, err = sut.Process()

	// assert
	assert.Equal(t, dummyData, dummyDataTemplate)
	assert.Equal(t, dummyStatusCode, result)
	assert.Equal(t, http.Header(dummyHeader), header)
	assert.Equal(t, dummyParseError, err)
}
