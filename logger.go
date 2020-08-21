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
	if isInterfaceValueNil(session.customization) {
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

// appRoot logs the given message as AppRoot category
func appRoot(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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

// apiEnter logs the given message as APIEnter category
func apiEnter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		APIEnter,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// apiRequest logs the given message as APIRequest category
func apiRequest(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		APIRequest,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// methodEnter logs the given message as MethodEnter category
func methodEnter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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

// methodParameter logs the given message as MethodParameter category
func methodParameter(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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

// methodLogic logs the given message as MethodLogic category
func methodLogic(session *session, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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

// networkCall logs the given message as NetworkCall category
func networkCall(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		NetworkCall,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// networkRequest logs the given message as NetworkRequest category
func networkRequest(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		NetworkRequest,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// networkResponse logs the given message as NetworkResponse category
func networkResponse(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		NetworkResponse,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// networkFinish logs the given message as NetworkFinish category
func networkFinish(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		NetworkFinish,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// methodReturn logs the given message as MethodReturn category
func methodReturn(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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

// methodExit logs the given message as MethodExit category
func methodExit(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
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

// apiResponse logs the given message as APIResponse category
func apiResponse(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		APIResponse,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}

// apiExit logs the given message as APIExit category
func apiExit(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	prepareLogging(
		session,
		APIExit,
		Info,
		category,
		subcategory,
		fmt.Sprintf(
			messageFormat,
			parameters...,
		),
	)
}
