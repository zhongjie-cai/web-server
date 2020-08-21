package webserver

import (
	"sort"
	"strings"
)

// LogType is the entry type of logging
type LogType int

// These are the enum definitions of log types and presets
const (
	AppRoot  LogType = 0
	APIEnter LogType = 1 << iota
	APIRequest
	MethodEnter
	MethodParameter
	MethodLogic
	NetworkCall
	NetworkRequest
	NetworkResponse
	NetworkFinish
	MethodReturn
	MethodExit
	APIResponse
	APIExit

	BasicTracing   LogType = MethodLogic
	GeneralTracing LogType = BasicTracing | APIEnter | APIExit
	VerboseTracing LogType = GeneralTracing | NetworkCall | NetworkFinish
	FullTracing    LogType = VerboseTracing | MethodEnter | MethodExit

	BasicDebugging   LogType = MethodLogic
	GeneralDebugging LogType = BasicDebugging | APIRequest | APIResponse
	VerboseDebugging LogType = GeneralDebugging | NetworkRequest | NetworkResponse
	FullDebugging    LogType = VerboseDebugging | MethodParameter | MethodReturn

	BasicLogging   LogType = BasicTracing | BasicDebugging
	GeneralLogging LogType = BasicLogging | GeneralTracing | GeneralDebugging
	VerboseLogging LogType = GeneralLogging | VerboseTracing | VerboseDebugging
	FullLogging    LogType = VerboseLogging | FullTracing | FullDebugging
)

// These are the string representations of log category and preset names
const (
	appRootName         string = "AppRoot"
	apiEnterName        string = "APIEnter"
	apiRequestName      string = "APIRequest"
	methodEnterName     string = "MethodEnter"
	methodParameterName string = "MethodParameter"
	methodLogicName     string = "MethodLogic"
	networkCallName     string = "NetworkCall"
	networkRequestName  string = "NetworkRequest"
	networkResponseName string = "NetworkResponse"
	networkFinishName   string = "NetworkFinish"
	methodReturnName    string = "MethodReturn"
	methodExitName      string = "MethodExit"
	apiResponseName     string = "APIResponse"
	apiExitName         string = "APIExit"

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
	APIEnter:        apiEnterName,
	APIRequest:      apiRequestName,
	MethodEnter:     methodEnterName,
	MethodParameter: methodParameterName,
	MethodLogic:     methodLogicName,
	NetworkCall:     networkCallName,
	NetworkRequest:  networkRequestName,
	NetworkResponse: networkResponseName,
	NetworkFinish:   networkFinishName,
	MethodReturn:    methodReturnName,
	MethodExit:      methodExitName,
	APIResponse:     apiResponseName,
	APIExit:         apiExitName,
}

var logTypeNameMapping = map[string]LogType{
	appRootName:          AppRoot,
	apiEnterName:         APIEnter,
	apiRequestName:       APIRequest,
	methodEnterName:      MethodEnter,
	methodParameterName:  MethodParameter,
	methodLogicName:      MethodLogic,
	networkCallName:      NetworkCall,
	networkRequestName:   NetworkRequest,
	networkResponseName:  NetworkResponse,
	networkFinishName:    NetworkFinish,
	methodReturnName:     MethodReturn,
	methodExitName:       MethodExit,
	apiResponseName:      APIResponse,
	apiExitName:          APIExit,
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
