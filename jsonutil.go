package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// These are constants to configure the JSON encoder
const (
	escapeHTML bool = false
)

// marshalIgnoreError returns the string representation of the given object; returns empty string in case of error
func marshalIgnoreError(v interface{}) string {
	var buffer = &bytes.Buffer{}
	var encoder = json.NewEncoder(buffer)
	encoder.SetEscapeHTML(escapeHTML)
	encoder.Encode(v)
	var result = string(buffer.Bytes())
	return strings.TrimRight(result, "\n")
}

func tryUnmarshalPrimitiveTypes(value string, dataTemplate interface{}) bool {
	if value == "" {
		return true
	}
	if reflect.TypeOf(dataTemplate) == reflect.TypeOf((*string)(nil)) {
		(*(dataTemplate).(*string)) = value
		return true
	}
	if reflect.TypeOf(dataTemplate) == reflect.TypeOf((*bool)(nil)) {
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
	}
	if reflect.TypeOf(dataTemplate) == reflect.TypeOf((*int)(nil)) {
		var parsedValue, parseError = strconv.Atoi(value)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*int)) = parsedValue
		return true
	}
	if reflect.TypeOf(dataTemplate) == reflect.TypeOf((*int64)(nil)) {
		var parsedValue, parseError = strconv.ParseInt(value, 0, 64)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*int64)) = parsedValue
		return true
	}
	if reflect.TypeOf(dataTemplate) == reflect.TypeOf((*float64)(nil)) {
		var parsedValue, parseError = strconv.ParseFloat(value, 64)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*float64)) = parsedValue
		return true
	}
	if reflect.TypeOf(dataTemplate) == reflect.TypeOf((*byte)(nil)) {
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
		"Unable to unmarshal value [%v] into data template",
		value,
	)
}
