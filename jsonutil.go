package webserver

import (
	"bytes"
	"reflect"
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
	var encoder = jsonNewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(v)
	var result = string(buffer.Bytes())
	return stringsTrimRight(result, "\n")
}

func tryUnmarshalPrimitiveTypes(value string, dataTemplate interface{}) bool {
	if value == "" {
		return true
	}
	switch reflectTypeOf(dataTemplate) {
	case typeOfString:
		(*(dataTemplate).(*string)) = value
		return true
	case typeOfBool:
		var parsedValue, parseError = strconvParseBool(
			stringsToLower(
				value,
			),
		)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*bool)) = parsedValue
		return true
	case typeOfInteger:
		var parsedValue, parseError = strconvAtoi(value)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*int)) = parsedValue
		return true
	case typeOfInt64:
		var parsedValue, parseError = strconvParseInt(value, 0, 64)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*int64)) = parsedValue
		return true
	case typeOfFloat64:
		var parsedValue, parseError = strconvParseFloat(value, 64)
		if parseError != nil {
			return false
		}
		(*(dataTemplate).(*float64)) = parsedValue
		return true
	case typeOfByte:
		var parsedValue, parseError = strconvParseUint(value, 0, 8)
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
	if isInterfaceValueNilFunc(dataTemplate) {
		return nil
	}
	if tryUnmarshalPrimitiveTypesFunc(
		value,
		dataTemplate,
	) {
		return nil
	}
	var noQuoteJSONError = jsonUnmarshal(
		[]byte(value),
		dataTemplate,
	)
	if noQuoteJSONError == nil {
		return nil
	}
	var withQuoteJSONError = jsonUnmarshal(
		[]byte("\""+value+"\""),
		dataTemplate,
	)
	if withQuoteJSONError == nil {
		return nil
	}
	return fmtErrorf(
		"Unable to unmarshal value [%v] into data template",
		value,
	)
}
