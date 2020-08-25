package webserver

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate_PatternError(t *testing.T) {
	// arrange
	var dummyPattern = "a(b"
	var dummyValue = "some value"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+dummyPattern+"$", pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterType(dummyPattern)

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.Error(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Anything_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeAnything+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeAnything

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Anything_Valid(t *testing.T) {
	// arrange
	var dummyValue = "some_thing to:BE!1-2evaluat3d??..."

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeAnything+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeAnything

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_String_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeString+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeString

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_String_Valid(t *testing.T) {
	// arrange
	var dummyValue = "some_value"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeString+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeString

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_String_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "some!value"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeString+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeString

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Integer_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeInteger+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeInteger

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Integer_Valid(t *testing.T) {
	// arrange
	var dummyValue = "1234"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeInteger+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeInteger

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Integer_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "123abc456"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeInteger+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeInteger

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_UUID_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeUUID+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeUUID

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_UUID_Valid(t *testing.T) {
	// arrange
	var dummyValue = "12345678-90ab-cdef-1234-567890abcdef"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeUUID+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeUUID

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_UUID_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "12345678-90ab-cdef-ghij-klmnopqrstuv"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeUUID+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeUUID

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Date_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeDate+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeDate

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Date_Valid(t *testing.T) {
	// arrange
	var dummyValue = "2006-01-02"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeDate+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeDate

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Date_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "2006/01/02"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeDate+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeDate

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Time_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Time_Valid_Full(t *testing.T) {
	// arrange
	var dummyValue = "15:04:05.789012"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Time_Valid_Short(t *testing.T) {
	// arrange
	var dummyValue = "15:04:05"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Time_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "15.04.05.0a"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_DateTime_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeDateTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_DateTime_Valid_WithTimeZone(t *testing.T) {
	// arrange
	var dummyValue = "2006-01-02T15:04:05.789+10:00"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeDateTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_DateTime_Valid_NoTimeZone(t *testing.T) {
	// arrange
	var dummyValue = "2006-01-02T15:04:05.789Z"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeDateTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_DateTime_Valid_Short(t *testing.T) {
	// arrange
	var dummyValue = "2006-01-02T15:04:05Z"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeDateTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_DateTime_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "2006/01/02 15.04.05Z"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeDateTime+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeDateTime

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Boolean_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeBoolean+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeBoolean

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Boolean_Valid(t *testing.T) {
	// arrange
	var dummyValue = "TrUe"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeBoolean+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeBoolean

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Boolean_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "False?"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeBoolean+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeBoolean

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Float_Empty(t *testing.T) {
	// arrange
	var dummyValue string

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeFloat+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeFloat

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Float_Valid(t *testing.T) {
	// arrange
	var dummyValue = "123.456"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeFloat+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeFloat

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluate_Float_Invalid(t *testing.T) {
	// arrange
	var dummyValue = "123.456.789"

	// mock
	createMock(t)

	// expect
	regexpMatchStringExpected = 1
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		assert.Equal(t, "^"+string(ParameterTypeFloat+"$"), pattern)
		assert.Equal(t, dummyValue, s)
		return regexp.MatchString(pattern, s)
	}

	// SUT
	var parameterType = ParameterTypeFloat

	// act
	var result, err = parameterType.Evaludate(
		dummyValue,
	)

	// assert
	assert.False(t, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}
