package webserver

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	httpClientWithCert *http.Client
	httpClientNoCert   *http.Client
)

func getClientForRequest(sendClientCert bool) *http.Client {
	if sendClientCert {
		return httpClientWithCert
	}
	return httpClientNoCert
}

func clientDo(
	httpClient *http.Client,
	httpRequest *http.Request,
) (*http.Response, error) {
	return httpClient.Do(
		httpRequest,
	)
}

func delayForRetry(retryDelay time.Duration) {
	time.Sleep(retryDelay)
}

func clientDoWithRetry(
	httpClient *http.Client,
	httpRequest *http.Request,
	connectivityRetryCount int,
	httpStatusRetryCount map[int]int,
	retryDelay time.Duration,
) (*http.Response, error) {
	var responseObject *http.Response
	var responseError error
	for {
		responseObject, responseError = clientDo(
			httpClient,
			httpRequest,
		)
		if responseError != nil {
			if connectivityRetryCount <= 0 {
				break
			}
			connectivityRetryCount--
		} else if responseObject != nil {
			var retry, found = httpStatusRetryCount[responseObject.StatusCode]
			if !found || retry <= 0 {
				break
			}
			httpStatusRetryCount[responseObject.StatusCode] = retry - 1
		} else {
			break
		}
		delayForRetry(
			retryDelay,
		)
	}
	return responseObject, responseError
}

func getHTTPTransport(
	skipServerCertVerification bool,
	clientCertificate *tls.Certificate,
	roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper,
) http.RoundTripper {
	var httpTransport = http.DefaultTransport
	if clientCertificate != nil {
		var tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{
				*clientCertificate,
			},
			InsecureSkipVerify: skipServerCertVerification,
		}
		httpTransport = &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		}
	}
	return roundTripperWrapper(
		httpTransport,
	)
}

func initializeHTTPClients(
	networkTimeout time.Duration,
	skipServerCertVerification bool,
	clientCertificate *tls.Certificate,
	roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper,
) {
	httpClientWithCert = &http.Client{
		Transport: getHTTPTransport(skipServerCertVerification, clientCertificate, roundTripperWrapper),
		Timeout:   networkTimeout,
	}
	httpClientNoCert = &http.Client{
		Transport: getHTTPTransport(skipServerCertVerification, nil, roundTripperWrapper),
		Timeout:   networkTimeout,
	}
}

// WebRequest is an interface for easy operating on network requests and responses
type WebRequest interface {
	// EnableRetry sets up automatic retry upon error of specific HTTP status codes; each entry maps an HTTP status code to how many times retry should happen if code matches
	EnableRetry(connectivityRetryCount int, httpStatusRetryCount map[int]int, retryDuration time.Duration)
	// Process sends the network request over the wire, retrieves and serialize the response to dataTemplate, and provides status code, header and error if applicable
	Process(dataTemplate interface{}) (statusCode int, responseHeader http.Header, responseError error)
	// ProcessRaw sends the network request over the wire, retrieves the response, and returns that response and error if applicable
	ProcessRaw() (responseObject *http.Response, responseError error)
}

type webRequest struct {
	session        *session
	method         string
	url            string
	payload        string
	header         map[string]string
	connRetry      int
	httpRetry      map[int]int
	sendClientCert bool
	retryDelay     time.Duration
}

// EnableRetry sets up automatic retry upon error of specific HTTP status codes; each entry maps an HTTP status code to how many times retry should happen if code matches; 0 stands for error not mapped to an HTTP status code, e.g. network or connectivity issue
func (webRequest *webRequest) EnableRetry(connectivityRetryCount int, httpStatusRetryCount map[int]int, retryDelay time.Duration) {
	webRequest.connRetry = connectivityRetryCount
	webRequest.httpRetry = httpStatusRetryCount
	webRequest.retryDelay = retryDelay
}

func customizeHTTPRequest(session *session, httpRequest *http.Request) *http.Request {
	return session.customization.WrapRequest(
		session,
		httpRequest,
	)
}

func createHTTPRequest(webRequest *webRequest) (*http.Request, error) {
	if webRequest == nil ||
		webRequest.session == nil {
		return nil, errWebRequestNil
	}
	var requestBody = strings.NewReader(
		webRequest.payload,
	)
	var requestObject, requestError = http.NewRequest(
		webRequest.method,
		webRequest.url,
		requestBody,
	)
	if requestError != nil {
		return nil, requestError
	}
	logWebcallStart(
		webRequest.session,
		webRequest.method,
		"",
		webRequest.url,
	)
	logWebcallRequest(
		webRequest.session,
		"Payload",
		"",
		webRequest.payload,
	)
	requestObject.Header = make(http.Header)
	for name, value := range webRequest.header {
		requestObject.Header.Add(name, value)
	}
	logWebcallRequest(
		webRequest.session,
		"Header",
		"",
		marshalIgnoreError(
			requestObject.Header,
		),
	)
	return customizeHTTPRequest(
		webRequest.session,
		requestObject,
	), nil
}

func logErrorResponse(session *session, responseError error, startTime time.Time) {
	logWebcallResponse(
		session,
		"Message",
		"",
		"%+v",
		responseError,
	)
	logWebcallFinish(
		session,
		"Error",
		"",
		"%s",
		time.Since(startTime),
	)
}

func logHTTPResponse(session *session, response *http.Response, startTime time.Time) {
	if response == nil {
		return
	}
	var (
		responseStatusCode = response.StatusCode
		responseBody, _    = ioutil.ReadAll(response.Body)
		responseHeaders    = response.Header
	)
	response.Body.Close()
	response.Body = ioutil.NopCloser(
		bytes.NewBuffer(
			responseBody,
		),
	)
	logWebcallRequest(
		session,
		"Header",
		"",
		marshalIgnoreError(
			responseHeaders,
		),
	)
	logWebcallResponse(
		session,
		"Body",
		"",
		string(responseBody),
	)
	logWebcallFinish(
		session,
		http.StatusText(responseStatusCode),
		strconv.Itoa(responseStatusCode),
		"%s",
		time.Since(startTime),
	)
}

func doRequestProcessing(webRequest *webRequest) (*http.Response, error) {
	if webRequest == nil ||
		webRequest.session == nil {
		return nil, errWebRequestNil
	}
	var requestObject, requestError = createHTTPRequest(
		webRequest,
	)
	if requestError != nil {
		return nil, requestError
	}
	var httpClient = getClientForRequest(
		webRequest.sendClientCert,
	)
	var startTime = getTimeNowUTC()
	var responseObject, responseError = clientDoWithRetry(
		httpClient,
		requestObject,
		webRequest.connRetry,
		webRequest.httpRetry,
		webRequest.retryDelay,
	)
	if responseError != nil {
		logErrorResponse(
			webRequest.session,
			responseError,
			startTime,
		)
	} else {
		logHTTPResponse(
			webRequest.session,
			responseObject,
			startTime,
		)
	}
	return responseObject, responseError
}

// ProcessRaw sends the network request over the wire, retrieves the response, and returns that response and error if applicable
func (webRequest *webRequest) ProcessRaw() (responseObject *http.Response, responseError error) {
	if webRequest == nil {
		return nil, errWebRequestNil
	}
	return doRequestProcessing(
		webRequest,
	)
}

func parseResponse(session *session, body io.ReadCloser, dataTemplate interface{}) error {
	var bodyBytes, bodyError = ioutil.ReadAll(
		body,
	)
	if bodyError != nil {
		return bodyError
	}
	var unmarshalError = tryUnmarshal(
		string(bodyBytes),
		dataTemplate,
	)
	if unmarshalError != nil {
		logWebcallResponse(
			session,
			"Body",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return errResponseInvalid
	}
	return nil
}

// Process sends the network request over the wire, retrieves and serialize the response to dataTemplate, and provides status code, header and error if applicable
func (webRequest *webRequest) Process(dataTemplate interface{}) (statusCode int, responseHeader http.Header, responseError error) {
	if webRequest == nil {
		return http.StatusInternalServerError, http.Header{}, errWebRequestNil
	}
	var responseObject *http.Response
	responseObject, responseError = doRequestProcessing(
		webRequest,
	)
	if responseError != nil {
		if responseObject == nil {
			return http.StatusInternalServerError, make(http.Header), responseError
		}
	} else {
		if responseObject == nil {
			return 0, make(http.Header), nil
		}
		responseError = parseResponse(
			webRequest.session,
			responseObject.Body,
			dataTemplate,
		)
	}
	return responseObject.StatusCode, responseObject.Header, responseError
}
