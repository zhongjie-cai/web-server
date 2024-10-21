package webserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate_PatternError(t *testing.T) {
	// arrange
	var dummyPattern = "a(b"
	var dummyValue = "some value"

	// SUT
	var parameterType = ParameterType(dummyPattern)

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.Error(t, err)
}

func TestEvaluate_Anything_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeAnything

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Anything_Valid(t *testing.T) {
	// arrange
	var dummyValue = "some_thing to:BE!1-2evaluat3d??..."

	// SUT
	var parameterType = ParameterTypeAnything

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_String_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeString

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_String_Valid(t *testing.T) {
	// arrange
	var dummyValue = "some_value"

	// SUT
	var parameterType = ParameterTypeString

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_String_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "some!value"

	// SUT
	var parameterType = ParameterTypeString

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Integer_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeInteger

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Integer_Valid(t *testing.T) {
	// arrange
	var dummyValue = "1234"

	// SUT
	var parameterType = ParameterTypeInteger

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Integer_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "123abc456"

	// SUT
	var parameterType = ParameterTypeInteger

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_UUID_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeUUID

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_UUID_Valid(t *testing.T) {
	// arrange
	var dummyValue = "12345678-90ab-cdef-1234-567890abcdef"

	// SUT
	var parameterType = ParameterTypeUUID

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_UUID_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "12345678-90ab-cdef-ghij-klmnopqrstuv"

	// SUT
	var parameterType = ParameterTypeUUID

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Date_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeDate

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Date_Valid(t *testing.T) {
	// arrange
	var dummyValue = "2006-01-02"

	// SUT
	var parameterType = ParameterTypeDate

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Date_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "2006/01/02"

	// SUT
	var parameterType = ParameterTypeDate

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Time_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Time_Valid_Full(t *testing.T) {
	// arrange
	var dummyValue = "15:04:05.789012"

	// SUT
	var parameterType = ParameterTypeTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Time_Valid_Short(t *testing.T) {
	// arrange
	var dummyValue = "15:04:05"

	// SUT
	var parameterType = ParameterTypeTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Time_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "15.04.05.0a"

	// SUT
	var parameterType = ParameterTypeTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_DateTime_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_DateTime_Valid_WithTimeZone(t *testing.T) {
	// arrange
	var dummyValue = "2006-01-02T15:04:05.789+10:00"

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_DateTime_Valid_NoTimeZone(t *testing.T) {
	// arrange
	var dummyValue = "2006-01-02T15:04:05.789Z"

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_DateTime_Valid_Short(t *testing.T) {
	// arrange
	var dummyValue = "2006-01-02T15:04:05Z"

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_DateTime_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "2006/01/02 15.04.05Z"

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Boolean_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeBoolean

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Boolean_Valid(t *testing.T) {
	// arrange
	var dummyValue = "TrUe"

	// SUT
	var parameterType = ParameterTypeBoolean

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Boolean_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "False?"

	// SUT
	var parameterType = ParameterTypeBoolean

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Float_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// SUT
	var parameterType = ParameterTypeFloat

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Float_Valid(t *testing.T) {
	// arrange
	var dummyValue = "123.456"

	// SUT
	var parameterType = ParameterTypeFloat

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEvaluate_Float_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "123.456.789"

	// SUT
	var parameterType = ParameterTypeFloat

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)
}
