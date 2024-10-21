package webserver

import "fmt"

func prepareLogging(
	session *session,
	logType LogType,
	logLevel LogLevel,
	category string,
	subcategory string,
	messageFormat string,
	parameters ...interface{},
) {
	if session == nil {
		return
	}
	session.customization.Log(
		session,
		logType,
		logLevel,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logAppRoot logs the given message as AppRoot category
func logAppRoot(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeAppRoot,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logEndpointEnter logs the given message as EndpointEnter category
func logEndpointEnter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeEndpointEnter,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logEndpointRequest logs the given message as EndpointRequest category
func logEndpointRequest(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeEndpointRequest,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logMethodEnter logs the given message as MethodEnter category
func logMethodEnter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeMethodEnter,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logMethodParameter logs the given message as MethodParameter category
func logMethodParameter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeMethodParameter,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logMethodLogic logs the given message as MethodLogic category
func logMethodLogic(session *session, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeMethodLogic,
		logLevel,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logWebcallStart logs the given message as WebcallStart category
func logWebcallStart(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeWebcallStart,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logWebcallRequest logs the given message as WebcallRequest category
func logWebcallRequest(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeWebcallRequest,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logWebcallResponse logs the given message as WebcallResponse category
func logWebcallResponse(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeWebcallResponse,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logWebcallFinish logs the given message as WebcallFinish category
func logWebcallFinish(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeWebcallFinish,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logMethodReturn logs the given message as MethodReturn category
func logMethodReturn(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeMethodReturn,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logMethodExit logs the given message as MethodExit category
func logMethodExit(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeMethodExit,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logEndpointResponse logs the given message as EndpointResponse category
func logEndpointResponse(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeEndpointResponse,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// logEndpointExit logs the given message as EndpointExit category
func logEndpointExit(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		LogTypeEndpointExit,
		LogLevelInfo,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}
