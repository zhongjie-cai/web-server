package webserver

// AppError is the error wrapper interface for all WebServiceTemplate service generated errors
type AppError interface {
	// AppBaseError is the error extending interface that is compatible with Golang error type
	AppBaseError
	// AppHTTPError is the error extending interface that can be translated into HTTP status code and response body
	AppHTTPError
	// AppContainerError is the error extending interface that can be used to wrap inner errors
	AppContainerError
}

// AppBaseError is the error extending interface that is compatible with Golang error type
type AppBaseError interface {
	// Error refers to the Golang built-in error interface method
	Error() string
	// ErrorCode returns the string representation of the error code enum
	ErrorCode() string
}

// AppHTTPError is the error extending interface that can be translated into HTTP status code and response body
type AppHTTPError interface {
	// HTTPStatusCode returns the corresponding HTTP status code mapped to the error code value
	HTTPStatusCode() int
	// HTTPResponseMessage returns the JSON representation of this error for HTTP response
	HTTPResponseMessage() string
}

// AppContainerError is the error extending interface that can be used to wrap inner errors
type AppContainerError interface {
	// Contains checks if the current error object or any of its inner errors contains the given error object
	Contains(err error) bool
	// Wrap wraps the given list of inner errors into the current app error object
	Wrap(innerErrors ...error) AppError
}

// These are print formatting related constants
const (
	errorMessageorFormat string = "(%v) %v"    // (ErrorCode) Message
	errorJoiningFormat   string = " -> [ %v ]" // -> [ joined contents ]
	errorSeparator       string = " | "        // {content 1} | {content 2} | {content 3}
)

type appError struct {
	Code        errorCode   `json:"code"`
	Message     string      `json:"message"`
	InnerErrors []*appError `json:"innerErrors,omitempty"`
}

func newAppError(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
	return &appError{
		Code:    errorCode,
		Message: errorMessageor,
		InnerErrors: cleanupInnerErrorsFunc(
			innerErrors,
		),
	}
}

func getErrorMessage(err error) string {
	return err.Error()
}

func printInnerErrors(innerErrors []*appError) string {
	if len(innerErrors) == 0 {
		return ""
	}
	var innerErrorMessages = []string{}
	for _, innerError := range innerErrors {
		if innerError != nil {
			innerErrorMessages = append(
				innerErrorMessages,
				getErrorMessageFunc(innerError),
			)
		}
	}
	return fmtSprintf(
		errorJoiningFormat,
		stringsJoin(
			innerErrorMessages,
			errorSeparator,
		),
	)
}

func (appError *appError) Error() string {
	var baseErrorMessage = fmtSprintf(
		errorMessageorFormat,
		appError.Code,
		appError.Message,
	)
	var innerErrorMessage = printInnerErrorsFunc(
		appError.InnerErrors,
	)
	return fmtSprint(
		baseErrorMessage,
		innerErrorMessage,
	)
}

// ErrorCode returns string representation of the error code of the app error
func (appError *appError) ErrorCode() string {
	return string(appError.Code)
}

// HTTPStatusCode returns HTTP status code according to the error code of the app error
func (appError *appError) HTTPStatusCode() int {
	return appError.Code.httpStatusCode()
}

// HTTPResponseMessage returns the JSON representation of this error for HTTP response
func (appError *appError) HTTPResponseMessage() string {
	var bytes, _ = jsonMarshal(appError)
	return string(bytes)
}

func equalsError(err, target error) bool {
	return err == target ||
		err.Error() == target.Error() ||
		errorsIs(err, target)
}

func appErrorContains(appError AppError, targetError error) bool {
	return appError.Contains(targetError)
}

func innerErrorContains(innerErrors []*appError, targetError error) bool {
	for _, innerError := range innerErrors {
		if appErrorContainsFunc(
			innerError,
			targetError,
		) {
			return true
		}
	}
	return false
}

// Contains checks if the current error object or any of its inner errors contains the given error object
func (appError *appError) Contains(targetError error) bool {
	return equalsErrorFunc(
		appError,
		targetError,
	) || innerErrorContainsFunc(
		appError.InnerErrors,
		targetError,
	)
}

func cleanupInnerErrors(innerErrors []error) []*appError {
	var cleanedInnerErrors = []*appError{}
	for _, innerError := range innerErrors {
		if innerError != nil {
			var typedError, isTyped = innerError.(*appError)
			if !isTyped {
				typedError = &appError{
					Code:    errorCodeGeneralFailure,
					Message: innerError.Error(),
				}
			}
			cleanedInnerErrors = append(
				cleanedInnerErrors,
				typedError,
			)
		}
	}
	return cleanedInnerErrors
}

// Wrap wraps the given list of inner errors into the current app error object
func (appError *appError) Wrap(innerErrors ...error) AppError {
	var cleanedInnerErrors = cleanupInnerErrorsFunc(
		innerErrors,
	)
	appError.InnerErrors = append(
		appError.InnerErrors,
		cleanedInnerErrors...,
	)
	return appError
}

// GetGeneralFailure creates a generic error based on GeneralFailure
func GetGeneralFailure(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeGeneralFailure,
		"A general error occurred during execution",
		innerErrors,
	)
}

// GetUnauthorized creates an error related to Unauthorized
func GetUnauthorized(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeUnauthorized,
		"Access denied due to authorization error",
		innerErrors,
	)
}

// GetInvalidOperation creates an error related to InvalidOperation
func GetInvalidOperation(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeInvalidOperation,
		"Operation (method) not allowed",
		innerErrors,
	)
}

// GetBadRequest creates an error related to BadRequest
func GetBadRequest(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeBadRequest,
		"Request URI or body is invalid",
		innerErrors,
	)
}

// GetNotFound creates an error related to NotFound
func GetNotFound(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeNotFound,
		"Requested resource is not found in the storage",
		innerErrors,
	)
}

// GetCircuitBreak creates an error related to CircuitBreak
func GetCircuitBreak(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeCircuitBreak,
		"Operation refused due to internal circuit break on correlation ID",
		innerErrors,
	)
}

// GetOperationLock creates an error related to OperationLock
func GetOperationLock(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeOperationLock,
		"Operation refused due to mutex lock on correlation ID or trip ID",
		innerErrors,
	)
}

// GetAccessForbidden creates an error related to AccessForbidden
func GetAccessForbidden(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeAccessForbidden,
		"Operation failed due to access forbidden",
		innerErrors,
	)
}

// GetDataCorruption creates an error related to DataCorruption
func GetDataCorruption(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeDataCorruption,
		"Operation failed due to internal storage data corruption",
		innerErrors,
	)
}

// GetNotImplemented creates an error related to NotImplemented
func GetNotImplemented(innerErrors ...error) AppError {
	return newAppErrorFunc(
		errorCodeNotImplemented,
		"Operation failed due to internal business logic not implemented",
		innerErrors,
	)
}
