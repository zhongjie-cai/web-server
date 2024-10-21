package webserver

import (
	"sort"
	"strings"
)

// LogType is the entry type of logging
type LogType int

// These are the enum definitions of log types and presets
const (
	LogTypeEndpointEnter LogType = 1 << iota
	LogTypeEndpointRequest
	LogTypeMethodEnter
	LogTypeMethodParameter
	LogTypeMethodLogic
	LogTypeWebcallStart
	LogTypeWebcallRequest
	LogTypeWebcallResponse
	LogTypeWebcallFinish
	LogTypeMethodReturn
	LogTypeMethodExit
	LogTypeEndpointResponse
	LogTypeEndpointExit

	LogTypeBasicTracing   LogType = LogTypeMethodLogic
	LogTypeGeneralTracing LogType = LogTypeBasicTracing | LogTypeEndpointEnter | LogTypeEndpointExit
	LogTypeVerboseTracing LogType = LogTypeGeneralTracing | LogTypeWebcallStart | LogTypeWebcallFinish
	LogTypeFullTracing    LogType = LogTypeVerboseTracing | LogTypeMethodEnter | LogTypeMethodExit

	LogTypeBasicDebugging   LogType = LogTypeMethodLogic
	LogTypeGeneralDebugging LogType = LogTypeBasicDebugging | LogTypeEndpointRequest | LogTypeEndpointResponse
	LogTypeVerboseDebugging LogType = LogTypeGeneralDebugging | LogTypeWebcallRequest | LogTypeWebcallResponse
	LogTypeFullDebugging    LogType = LogTypeVerboseDebugging | LogTypeMethodParameter | LogTypeMethodReturn

	LogTypeBasicLogging   LogType = LogTypeBasicTracing | LogTypeBasicDebugging
	LogTypeGeneralLogging LogType = LogTypeBasicLogging | LogTypeGeneralTracing | LogTypeGeneralDebugging
	LogTypeVerboseLogging LogType = LogTypeGeneralLogging | LogTypeVerboseTracing | LogTypeVerboseDebugging
	LogTypeFullLogging    LogType = LogTypeVerboseLogging | LogTypeFullTracing | LogTypeFullDebugging

	LogTypeAppRoot LogType = 0
)

// These are the string representations of log category and preset names
const (
	apiEnterLogTypeName        string = "EndpointEnter"
	apiRequestLogTypeName      string = "EndpointRequest"
	methodEnterLogTypeName     string = "MethodEnter"
	methodParameterLogTypeName string = "MethodParameter"
	methodLogicLogTypeName     string = "MethodLogic"
	webcallCallLogTypeName     string = "WebcallStart"
	webcallRequestLogTypeName  string = "WebcallRequest"
	webcallResponseLogTypeName string = "WebcallResponse"
	webcallFinishLogTypeName   string = "WebcallFinish"
	methodReturnLogTypeName    string = "MethodReturn"
	methodExitLogTypeName      string = "MethodExit"
	apiResponseLogTypeName     string = "EndpointResponse"
	apiExitLogTypeName         string = "EndpointExit"

	basicTracingLogTypeName   string = "BasicTracing"
	generalTracingLogTypeName string = "GeneralTracing"
	verboseTracingLogTypeName string = "VerboseTracing"
	fullTracingLogTypeName    string = "FullTracing"

	basicDebuggingLogTypeName   string = "BasicDebugging"
	generalDebuggingLogTypeName string = "GeneralDebugging"
	verboseDebuggingLogTypeName string = "VerboseDebugging"
	fullDebuggingLogTypeName    string = "FullDebugging"

	basicLoggingLogTypeName   string = "BasicLogging"
	generalLoggingLogTypeName string = "GeneralLogging"
	verboseLoggingLogTypeName string = "VerboseLogging"
	fullLoggingLogTypeName    string = "FullLogging"

	appRootLogTypeName string = "AppRoot"
)

var supportedLogTypes = map[LogType]string{
	LogTypeEndpointEnter:    apiEnterLogTypeName,
	LogTypeEndpointRequest:  apiRequestLogTypeName,
	LogTypeMethodEnter:      methodEnterLogTypeName,
	LogTypeMethodParameter:  methodParameterLogTypeName,
	LogTypeMethodLogic:      methodLogicLogTypeName,
	LogTypeWebcallStart:     webcallCallLogTypeName,
	LogTypeWebcallRequest:   webcallRequestLogTypeName,
	LogTypeWebcallResponse:  webcallResponseLogTypeName,
	LogTypeWebcallFinish:    webcallFinishLogTypeName,
	LogTypeMethodReturn:     methodReturnLogTypeName,
	LogTypeMethodExit:       methodExitLogTypeName,
	LogTypeEndpointResponse: apiResponseLogTypeName,
	LogTypeEndpointExit:     apiExitLogTypeName,
}

var logTypeNameMapping = map[string]LogType{
	apiEnterLogTypeName:         LogTypeEndpointEnter,
	apiRequestLogTypeName:       LogTypeEndpointRequest,
	methodEnterLogTypeName:      LogTypeMethodEnter,
	methodParameterLogTypeName:  LogTypeMethodParameter,
	methodLogicLogTypeName:      LogTypeMethodLogic,
	webcallCallLogTypeName:      LogTypeWebcallStart,
	webcallRequestLogTypeName:   LogTypeWebcallRequest,
	webcallResponseLogTypeName:  LogTypeWebcallResponse,
	webcallFinishLogTypeName:    LogTypeWebcallFinish,
	methodReturnLogTypeName:     LogTypeMethodReturn,
	methodExitLogTypeName:       LogTypeMethodExit,
	apiResponseLogTypeName:      LogTypeEndpointResponse,
	apiExitLogTypeName:          LogTypeEndpointExit,
	basicTracingLogTypeName:     LogTypeBasicTracing,
	generalTracingLogTypeName:   LogTypeGeneralTracing,
	verboseTracingLogTypeName:   LogTypeVerboseTracing,
	fullTracingLogTypeName:      LogTypeFullTracing,
	basicDebuggingLogTypeName:   LogTypeBasicDebugging,
	generalDebuggingLogTypeName: LogTypeGeneralDebugging,
	verboseDebuggingLogTypeName: LogTypeVerboseDebugging,
	fullDebuggingLogTypeName:    LogTypeFullDebugging,
	basicLoggingLogTypeName:     LogTypeBasicLogging,
	generalLoggingLogTypeName:   LogTypeGeneralLogging,
	verboseLoggingLogTypeName:   LogTypeVerboseLogging,
	fullLoggingLogTypeName:      LogTypeFullLogging,
	appRootLogTypeName:          LogTypeAppRoot,
}

// FromString converts a LogType flag instance to its string representation
func (logtype LogType) String() string {
	if logtype == LogTypeAppRoot {
		return appRootLogTypeName
	}
	var result []string
	for key, value := range supportedLogTypes {
		if logtype&key == key {
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return strings.Join(result, "|")
}

// HasFlag checks whether this log category has the flag set or not
func (logtype LogType) HasFlag(flag LogType) bool {
	if flag == LogTypeAppRoot {
		return true
	}
	if logtype&flag == flag {
		return true
	}
	return false
}

// NewLogType converts a string representation of LogType flag to its strongly typed instance
func NewLogType(value string) LogType {
	var splitValues = strings.Split(
		value,
		"|",
	)
	var combinedLogType LogType
	for _, splitValue := range splitValues {
		var logType, found = logTypeNameMapping[splitValue]
		if found {
			combinedLogType = combinedLogType | logType
		}
	}
	return combinedLogType
}
