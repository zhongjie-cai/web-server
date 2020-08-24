package webserver

import (
	"sort"
	"strings"
)

// LogType is the entry type of logging
type LogType int

// These are the enum definitions of log types and presets
const (
	AppRoot       LogType = 0
	EndpointEnter LogType = 1 << iota
	EndpointRequest
	MethodEnter
	MethodParameter
	MethodLogic
	WebcallStart
	WebcallRequest
	WebcallResponse
	WebcallFinish
	MethodReturn
	MethodExit
	EndpointResponse
	EndpointExit

	BasicTracing   LogType = MethodLogic
	GeneralTracing LogType = BasicTracing | EndpointEnter | EndpointExit
	VerboseTracing LogType = GeneralTracing | WebcallStart | WebcallFinish
	FullTracing    LogType = VerboseTracing | MethodEnter | MethodExit

	BasicDebugging   LogType = MethodLogic
	GeneralDebugging LogType = BasicDebugging | EndpointRequest | EndpointResponse
	VerboseDebugging LogType = GeneralDebugging | WebcallRequest | WebcallResponse
	FullDebugging    LogType = VerboseDebugging | MethodParameter | MethodReturn

	BasicLogging   LogType = BasicTracing | BasicDebugging
	GeneralLogging LogType = BasicLogging | GeneralTracing | GeneralDebugging
	VerboseLogging LogType = GeneralLogging | VerboseTracing | VerboseDebugging
	FullLogging    LogType = VerboseLogging | FullTracing | FullDebugging
)

// These are the string representations of log category and preset names
const (
	appRootName         string = "AppRoot"
	apiEnterName        string = "EndpointEnter"
	apiRequestName      string = "EndpointRequest"
	methodEnterName     string = "MethodEnter"
	methodParameterName string = "MethodParameter"
	methodLogicName     string = "MethodLogic"
	webcallCallName     string = "WebcallStart"
	webcallRequestName  string = "WebcallRequest"
	webcallResponseName string = "WebcallResponse"
	webcallFinishName   string = "WebcallFinish"
	methodReturnName    string = "MethodReturn"
	methodExitName      string = "MethodExit"
	apiResponseName     string = "EndpointResponse"
	apiExitName         string = "EndpointExit"

	basicTracingName   string = "BasicTracing"
	generalTracingName string = "GeneralTracing"
	verboseTracingName string = "VerboseTracing"
	fullTracingName    string = "FullTracing"

	basicDebuggingName   string = "BasicDebugging"
	generalDebuggingName string = "GeneralDebugging"
	verboseDebuggingName string = "VerboseDebugging"
	fullDebuggingName    string = "FullDebugging"

	basicLoggingName   string = "BasicLogging"
	generalLoggingName string = "GeneralLogging"
	verboseLoggingName string = "VerboseLogging"
	fullLoggingName    string = "FullLogging"
)

var supportedLogTypes = map[LogType]string{
	EndpointEnter:    apiEnterName,
	EndpointRequest:  apiRequestName,
	MethodEnter:      methodEnterName,
	MethodParameter:  methodParameterName,
	MethodLogic:      methodLogicName,
	WebcallStart:     webcallCallName,
	WebcallRequest:   webcallRequestName,
	WebcallResponse:  webcallResponseName,
	WebcallFinish:    webcallFinishName,
	MethodReturn:     methodReturnName,
	MethodExit:       methodExitName,
	EndpointResponse: apiResponseName,
	EndpointExit:     apiExitName,
}

var logTypeNameMapping = map[string]LogType{
	appRootName:          AppRoot,
	apiEnterName:         EndpointEnter,
	apiRequestName:       EndpointRequest,
	methodEnterName:      MethodEnter,
	methodParameterName:  MethodParameter,
	methodLogicName:      MethodLogic,
	webcallCallName:      WebcallStart,
	webcallRequestName:   WebcallRequest,
	webcallResponseName:  WebcallResponse,
	webcallFinishName:    WebcallFinish,
	methodReturnName:     MethodReturn,
	methodExitName:       MethodExit,
	apiResponseName:      EndpointResponse,
	apiExitName:          EndpointExit,
	basicTracingName:     BasicTracing,
	generalTracingName:   GeneralTracing,
	verboseTracingName:   VerboseTracing,
	fullTracingName:      FullTracing,
	basicDebuggingName:   BasicDebugging,
	generalDebuggingName: GeneralDebugging,
	verboseDebuggingName: VerboseDebugging,
	fullDebuggingName:    FullDebugging,
	basicLoggingName:     BasicLogging,
	generalLoggingName:   GeneralLogging,
	verboseLoggingName:   VerboseLogging,
	fullLoggingName:      FullLogging,
}

// FromString converts a LogType flag instance to its string representation
func (logtype LogType) String() string {
	if logtype == AppRoot {
		return appRootName
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
	if flag == AppRoot {
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
