package webserver

import (
	"fmt"
)

func prepareLogging(
	session *session,
	logType LogType,
	logLevel LogLevel,
	category string,
	subCategory string,
	description string,
) {
	if session == nil {
		return
	}
	session.customization.Log(
		session,
		logType,
		logLevel,
		category,
		subCategory,
		description,
	)
}

// logAppRoot logs the given message as AppRoot category
func logAppRoot(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		AppRoot,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logEndpointEnter logs the given message as EndpointEnter category
func logEndpointEnter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		EndpointEnter,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logEndpointRequest logs the given message as EndpointRequest category
func logEndpointRequest(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		EndpointRequest,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logMethodEnter logs the given message as MethodEnter category
func logMethodEnter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		MethodEnter,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logMethodParameter logs the given message as MethodParameter category
func logMethodParameter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		MethodParameter,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logMethodLogic logs the given message as MethodLogic category
func logMethodLogic(session *session, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		MethodLogic,
		logLevel,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logWebcallStart logs the given message as WebcallStart category
func logWebcallStart(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		WebcallStart,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logWebcallRequest logs the given message as WebcallRequest category
func logWebcallRequest(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		WebcallRequest,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logWebcallResponse logs the given message as WebcallResponse category
func logWebcallResponse(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		WebcallResponse,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logWebcallFinish logs the given message as WebcallFinish category
func logWebcallFinish(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		WebcallFinish,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logMethodReturn logs the given message as MethodReturn category
func logMethodReturn(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		MethodReturn,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logMethodExit logs the given message as MethodExit category
func logMethodExit(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		MethodExit,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logEndpointResponse logs the given message as EndpointResponse category
func logEndpointResponse(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		EndpointResponse,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// logEndpointExit logs the given message as EndpointExit category
func logEndpointExit(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		EndpointExit,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}
