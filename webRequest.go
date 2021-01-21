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
	// AddQuery adds a query to the request URL for sending through HTTP
	AddQuery(name string, value string) WebRequest
	// AddHeader adds a header to the request Header for sending through HTTP
	AddHeader(name string, value string) WebRequest
	// SetupRetry sets up automatic retry upon error of specific HTTP status codes; each entry maps an HTTP status code to how many times retry should happen if code matches
	SetupRetry(connectivityRetryCount int, httpStatusRetryCount map[int]int, retryDelay time.Duration) WebRequest
	// Anticipate registers a data template to be deserialized to when the given range of HTTP status codes are returned during the processing of the web request; latter registration overrides former when overlapping
	Anticipate(beginStatusCode int, endStatusCode int, dataTemplate interface{}) WebRequest
	// Process sends the webcall request over the wire, retrieves and serialize the response to registered data templates, and returns status code, header and error accordingly
	Process() (statusCode int, responseHeader http.Header, responseError error)
}

type dataReceiver struct {
	beginStatusCode int
	endStatusCode   int
	dataTemplate    interface{}
}

type webRequest struct {
	session        *session
	method         string
	url            string
	payload        string
	query          map[string][]string
	header         map[string][]string
	connRetry      int
	httpRetry      map[int]int
	sendClientCert bool
	retryDelay     time.Duration
	dataReceivers  []dataReceiver
}

// AddQuery adds a query to the request URL for sending through HTTP
func (webRequest *webRequest) AddQuery(name string, value string) WebRequest {
	if webRequest.query == nil {
		webRequest.query = make(map[string][]string)
	}
	var queryValues = webRequest.query[name]
	queryValues = append(
		queryValues,
		value,
	)
	webRequest.query[name] = queryValues
	return webRequest
}

// AddHeader adds a header to the request Header for sending through HTTP
func (webRequest *webRequest) AddHeader(name string, value string) WebRequest {
	if webRequest.header == nil {
		webRequest.header = make(map[string][]string)
	}
	var headerValues = webRequest.header[name]
	headerValues = append(
		headerValues,
		value,
	)
	webRequest.header[name] = headerValues
	return webRequest
}

// SetupRetry sets up automatic retry upon error of specific HTTP status codes; each entry maps an HTTP status code to how many times retry should happen if code matches; 0 stands for error not mapped to an HTTP status code, e.g. webcall or connectivity issue
func (webRequest *webRequest) SetupRetry(connectivityRetryCount int, httpStatusRetryCount map[int]int, retryDelay time.Duration) WebRequest {
	webRequest.connRetry = connectivityRetryCount
	webRequest.httpRetry = httpStatusRetryCount
	webRequest.retryDelay = retryDelay
	return webRequest
}

// Anticipate registers a data template to be deserialized to when the given range of HTTP status codes are returned during the processing of the web request; latter registration overrides former when overlapping
func (webRequest *webRequest) Anticipate(beginStatusCode int, endStatusCode int, dataTemplate interface{}) WebRequest {
	webRequest.dataReceivers = append(
		webRequest.dataReceivers,
		dataReceiver{
			beginStatusCode: beginStatusCode,
			endStatusCode:   endStatusCode,
			dataTemplate:    dataTemplate,
		},
	)
	return webRequest
}

func createQueryString(
	query map[string][]string,
) string {
	var queryStrings []string
	for name, values := range query {
		if name == "" {
			continue
		}
		for _, value := range values {
			queryStrings = append(
				queryStrings,
				fmtSprintf(
					"%v=%v",
					urlQueryEscape(name),
					urlQueryEscape(value),
				),
			)
		}
	}
	return stringsJoin(
		queryStrings,
		"&",
	)
}

func generateRequestURL(
	baseURL string,
	query map[string][]string,
) string {
	if len(query) == 0 {
		return baseURL
	}
	var queryString = createQueryStringFunc(
		query,
	)
	if queryString == "" {
		return baseURL
	}
	return fmtSprintf(
		"%v?%v",
		baseURL,
		queryString,
	)
}

func createHTTPRequest(webRequest *webRequest) (*http.Request, error) {
	if webRequest == nil ||
		webRequest.session == nil {
		return nil,
			newAppErrorFunc(
				errorCodeGeneralFailure,
				errorMessageWebRequestNil,
				[]error{},
			)
	}
	var requestURL = generateRequestURLFunc(
		webRequest.url,
		webRequest.query,
	)
	var requestBody = stringsNewReader(
		webRequest.payload,
	)
	var requestObject, requestError = httpNewRequest(
		webRequest.method,
		requestURL,
		requestBody,
	)
	if requestError != nil {
		return nil, requestError
	}
	logWebcallStartFunc(
		webRequest.session,
		webRequest.method,
		webRequest.url,
		requestURL,
	)
	logWebcallRequestFunc(
		webRequest.session,
		"Payload",
		"Content",
		webRequest.payload,
	)
	requestObject.Header = make(http.Header)
	for name, values := range webRequest.header {
		for _, value := range values {
			requestObject.Header.Add(name, value)
		}
	}
	logWebcallRequestFunc(
		webRequest.session,
		"Header",
		"Content",
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
		"Error",
		"Content",
		"%+v",
		responseError,
	)
	logWebcallFinishFunc(
		session,
		"Error",
		"-1",
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
		"Content",
		marshalIgnoreErrorFunc(
			responseHeaders,
		),
	)
	logWebcallResponseFunc(
		session,
		"Body",
		"Content",
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
		return nil,
			newAppErrorFunc(
				errorCodeGeneralFailure,
				errorMessageWebRequestNil,
				[]error{},
			)
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

func getDataTemplate(session *session, statusCode int, dataReceivers []dataReceiver) interface{} {
	var dataTemplate interface{}
	for _, dataReceiver := range dataReceivers {
		if dataReceiver.beginStatusCode <= statusCode &&
			dataReceiver.endStatusCode > statusCode {
			dataTemplate = dataReceiver.dataTemplate
		}
	}
	if isInterfaceValueNilFunc(dataTemplate) {
		logWebcallResponseFunc(
			session,
			"Body",
			"Receiver",
			"No data template registered for HTTP status %v",
			statusCode,
		)
	}
	return dataTemplate
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
		return newAppErrorFunc(
			errorCodeGeneralFailure,
			errorMessageResponseInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

// Process sends the webcall request over the wire, retrieves and serialize the response to dataTemplate, and provides status code, header and error if applicable
func (webRequest *webRequest) Process() (statusCode int, responseHeader http.Header, responseError error) {
	if webRequest == nil ||
		webRequest.session == nil {
		return http.StatusInternalServerError,
			http.Header{},
			newAppErrorFunc(
				errorCodeGeneralFailure,
				errorMessageWebRequestNil,
				[]error{},
			)
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
		var dataTemplate = getDataTemplateFunc(
			webRequest.session,
			responseObject.StatusCode,
			webRequest.dataReceivers,
		)
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
