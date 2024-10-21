package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	typeOfString  = reflect.TypeOf((*string)(nil))
	typeOfBool    = reflect.TypeOf((*bool)(nil))
	typeOfInteger = reflect.TypeOf((*int)(nil))
	typeOfInt64   = reflect.TypeOf((*int64)(nil))
	typeOfFloat64 = reflect.TypeOf((*float64)(nil))
	typeOfByte    = reflect.TypeOf((*byte)(nil))
)

// marshalIgnoreError returns the string representation of the given object; returns empty string in case of error
func marshalIgnoreError(v interface{}) string {
	var buffer = &bytes.Buffer{}
	var encoder = json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(v)
	var result = buffer.String()
	return strings.TrimRight(result, "\n")
}

func tryUnmarshalPrimitiveTypes(value string, dataTemplate interface{}) bool {
	if value == "" {
		return true
	}
	switch reflect.TypeOf(dataTemplate) {
	case typeOfString:
		(*(dataTemplate).(*string)) = value
		return true
	case typeOfBool:
		var parsedValue, parseError = strconv.ParseBool(
			strings.ToLower(
				value,
			),
		)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*bool)) = parsedValue
		return true
	case typeOfInteger:
		var parsedValue, parseError = strconv.Atoi(value)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*int)) = parsedValue
		return true
	case typeOfInt64:
		var parsedValue, parseError = strconv.ParseInt(value, 0, 64)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*int64)) = parsedValue
		return true
	case typeOfFloat64:
		var parsedValue, parseError = strconv.ParseFloat(value, 64)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*float64)) = parsedValue
		return true
	case typeOfByte:
		var parsedValue, parseError = strconv.ParseUint(value, 0, 8)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*byte)) = byte(parsedValue)
		return true
	}
	return false
}

// tryUnmarshal tries to unmarshal given value to dataTemplate
func tryUnmarshal(value string, dataTemplate interface{}) error {
	if isInterfaceValueNil(dataTemplate) {
		return nil
	}
	if tryUnmarshalPrimitiveTypes(
		value,
		dataTemplate,
	) {
		return nil
	}
	var noQuoteJSONError = json.Unmarshal(
		[]byte(value),
		dataTemplate,
	)
	if noQuoteJSONError == nil {
		return nil
	}
	var withQuoteJSONError = json.Unmarshal(
		[]byte("\""+value+"\""),
		dataTemplate,
	)
	if withQuoteJSONError == nil {
		return nil
	}
	return fmt.Errorf(
		"unable to unmarshal value [%v] into data template",
		value,
	)
}
