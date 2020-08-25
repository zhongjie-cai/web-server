package webserver

import (
	"crypto/tls"
	"io"
	"net/http"
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
		responseObject, responseError = clientDoFunc(
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
		timeSleep(
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
	webcallTimeout time.Duration,
	skipServerCertVerification bool,
	clientCertificate *tls.Certificate,
	roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper,
) {
	httpClientWithCert = &http.Client{
		Transport: getHTTPTransportFunc(skipServerCertVerification, clientCertificate, roundTripperWrapper),
		Timeout:   webcallTimeout,
	}
	httpClientNoCert = &http.Client{
		Transport: getHTTPTransportFunc(skipServerCertVerification, nil, roundTripperWrapper),
		Timeout:   webcallTimeout,
	}
}

// WebRequest is an interface for easy operating on webcall requests and responses
type WebRequest interface {
	// EnableRetry sets up automatic retry upon error of specific HTTP status codes; each entry maps an HTTP status code to how many times retry should happen if code matches
	EnableRetry(connectivityRetryCount int, httpStatusRetryCount map[int]int, retryDuration time.Duration)
	// Process sends the webcall request over the wire, retrieves and serialize the response to dataTemplate, and provides status code, header and error if applicable
	Process(dataTemplate interface{}) (statusCode int, responseHeader http.Header, responseError error)
	// ProcessRaw sends the webcall request over the wire, retrieves the response, and returns that response and error if applicable
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

// EnableRetry sets up automatic retry upon error of specific HTTP status codes; each entry maps an HTTP status code to how many times retry should happen if code matches; 0 stands for error not mapped to an HTTP status code, e.g. webcall or connectivity issue
func (webRequest *webRequest) EnableRetry(connectivityRetryCount int, httpStatusRetryCount map[int]int, retryDelay time.Duration) {
	webRequest.connRetry = connectivityRetryCount
	webRequest.httpRetry = httpStatusRetryCount
	webRequest.retryDelay = retryDelay
}

func createHTTPRequest(webRequest *webRequest) (*http.Request, error) {
	if webRequest == nil ||
		webRequest.session == nil {
		return nil, ErrWebRequestNil
	}
	var requestBody = stringsNewReader(
		webRequest.payload,
	)
	var requestObject, requestError = httpNewRequest(
		webRequest.method,
		webRequest.url,
		requestBody,
	)
	if requestError != nil {
		return nil, requestError
	}
	logWebcallStartFunc(
		webRequest.session,
		webRequest.method,
		"",
		webRequest.url,
	)
	logWebcallRequestFunc(
		webRequest.session,
		"Payload",
		"",
		webRequest.payload,
	)
	requestObject.Header = make(http.Header)
	for name, value := range webRequest.header {
		requestObject.Header.Add(name, value)
	}
	logWebcallRequestFunc(
		webRequest.session,
		"Header",
		"",
		marshalIgnoreErrorFunc(
			requestObject.Header,
		),
	)
	return webRequest.session.customization.WrapRequest(
		webRequest.session,
		requestObject,
	), nil
}

func logErrorResponse(session *session, responseError error, startTime time.Time) {
	logWebcallResponseFunc(
		session,
		"Message",
		"",
		"%+v",
		responseError,
	)
	logWebcallFinishFunc(
		session,
		"Error",
		"",
		"%s",
		timeSince(startTime),
	)
}

func logSuccessResponse(session *session, response *http.Response, startTime time.Time) {
	if response == nil {
		return
	}
	var (
		responseStatusCode = response.StatusCode
		responseBody, _    = ioutilReadAll(response.Body)
		responseHeaders    = response.Header
	)
	response.Body.Close()
	response.Body = ioutilNopCloser(
		bytesNewBuffer(
			responseBody,
		),
	)
	logWebcallResponseFunc(
		session,
		"Header",
		"",
		marshalIgnoreErrorFunc(
			responseHeaders,
		),
	)
	logWebcallResponseFunc(
		session,
		"Body",
		"",
		string(responseBody),
	)
	logWebcallFinishFunc(
		session,
		httpStatusText(responseStatusCode),
		strconvItoa(responseStatusCode),
		"%s",
		timeSince(startTime),
	)
}

func doRequestProcessing(webRequest *webRequest) (*http.Response, error) {
	if webRequest == nil ||
		webRequest.session == nil {
		return nil, ErrWebRequestNil
	}
	var requestObject, requestError = createHTTPRequestFunc(
		webRequest,
	)
	if requestError != nil {
		return nil, requestError
	}
	var httpClient = getClientForRequestFunc(
		webRequest.sendClientCert,
	)
	var startTime = getTimeNowUTCFunc()
	var responseObject, responseError = clientDoWithRetryFunc(
		httpClient,
		requestObject,
		webRequest.connRetry,
		webRequest.httpRetry,
		webRequest.retryDelay,
	)
	if responseError != nil {
		logErrorResponseFunc(
			webRequest.session,
			responseError,
			startTime,
		)
	} else {
		logSuccessResponseFunc(
			webRequest.session,
			responseObject,
			startTime,
		)
	}
	return responseObject, responseError
}

// ProcessRaw sends the webcall request over the wire, retrieves the response, and returns that response and error if applicable
func (webRequest *webRequest) ProcessRaw() (responseObject *http.Response, responseError error) {
	if webRequest == nil {
		return nil, ErrWebRequestNil
	}
	return doRequestProcessingFunc(
		webRequest,
	)
}

func parseResponse(session *session, body io.ReadCloser, dataTemplate interface{}) error {
	var bodyBytes, bodyError = ioutilReadAll(
		body,
	)
	if bodyError != nil {
		return bodyError
	}
	var unmarshalError = tryUnmarshalFunc(
		string(bodyBytes),
		dataTemplate,
	)
	if unmarshalError != nil {
		logWebcallResponseFunc(
			session,
			"Body",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return ErrResponseInvalid
	}
	return nil
}

// Process sends the webcall request over the wire, retrieves and serialize the response to dataTemplate, and provides status code, header and error if applicable
func (webRequest *webRequest) Process(dataTemplate interface{}) (statusCode int, responseHeader http.Header, responseError error) {
	if webRequest == nil ||
		webRequest.session == nil {
		return http.StatusInternalServerError, http.Header{}, ErrWebRequestNil
	}
	var responseObject *http.Response
	responseObject, responseError = doRequestProcessingFunc(
		webRequest,
	)
	if responseError != nil {
		if responseObject == nil {
			return http.StatusInternalServerError,
				make(http.Header),
				responseError
		}
	} else {
		if responseObject == nil {
			return 0, make(http.Header), nil
		}
		responseError = parseResponseFunc(
			webRequest.session,
			responseObject.Body,
			dataTemplate,
		)
	}
	return responseObject.StatusCode,
		responseObject.Header,
		responseError
}
