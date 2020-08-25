package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMarshalIgnoreError_Empty(t *testing.T) {
	// arrange
	var dummyObject *struct {
		Foo  string
		Test int
	}
	var expectedResult = "null"

	// mock
	createMock(t)

	// expect
	jsonNewEncoderExpected = 1
	jsonNewEncoder = func(w io.Writer) *json.Encoder {
		jsonNewEncoderCalled++
		return json.NewEncoder(w)
	}
	stringsTrimRightExpected = 1
	stringsTrimRight = func(s string, cutset string) string {
		stringsTrimRightCalled++
		return strings.TrimRight(s, cutset)
	}

	// SUT + act
	var result = marshalIgnoreError(
		dummyObject,
	)

	// assert
	assert.Equal(t, expectedResult, result)

	// verify
	verifyAll(t)
}

func TestMarshalIgnoreError_Success(t *testing.T) {
	// arrange
	var dummyObject = struct {
		Foo  string
		Test int
	}{
		"<bar />",
		123,
	}
	var expectedResult = "{\"Foo\":\"<bar />\",\"Test\":123}"

	// mock
	createMock(t)

	// expect
	jsonNewEncoderExpected = 1
	jsonNewEncoder = func(w io.Writer) *json.Encoder {
		jsonNewEncoderCalled++
		return json.NewEncoder(w)
	}
	stringsTrimRightExpected = 1
	stringsTrimRight = func(s string, cutset string) string {
		stringsTrimRightCalled++
		return strings.TrimRight(s, cutset)
	}

	// SUT + act
	var result = marshalIgnoreError(
		dummyObject,
	)

	// assert
	assert.Equal(t, expectedResult, result)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Empty(t *testing.T) {
	// arrange
	var dummyValue string
	var dummyDataTemplate string

	// mock
	createMock(t)

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.True(t, result)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_String(t *testing.T) {
	// arrange
	var dummyValue = "some value"
	var dummyDataTemplate string

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Bool_Error(t *testing.T) {
	// arrange
	var dummyValue = "some invalid bool"
	var dummyDataTemplate bool

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	stringsToLowerExpected = 1
	stringsToLower = func(s string) string {
		stringsToLowerCalled++
		return strings.ToLower(s)
	}
	strconvParseBoolExpected = 1
	strconvParseBool = func(str string) (bool, error) {
		strconvParseBoolCalled++
		return strconv.ParseBool(str)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Bool_NoError(t *testing.T) {
	// arrange
	var dummyValue = rand.Intn(100) < 50
	var dummyValueString = fmt.Sprintf("%v", dummyValue)
	var dummyDataTemplate bool

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	stringsToLowerExpected = 1
	stringsToLower = func(s string) string {
		stringsToLowerCalled++
		return strings.ToLower(s)
	}
	strconvParseBoolExpected = 1
	strconvParseBool = func(str string) (bool, error) {
		strconvParseBoolCalled++
		return strconv.ParseBool(str)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValueString,
		&dummyDataTemplate,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Integer_Error(t *testing.T) {
	// arrange
	var dummyValue = "some invalid integer"
	var dummyDataTemplate int

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	strconvAtoiExpected = 1
	strconvAtoi = func(str string) (int, error) {
		strconvAtoiCalled++
		return strconv.Atoi(str)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Integer_NoError(t *testing.T) {
	// arrange
	var dummyValue = rand.Intn(math.MaxInt32)
	var dummyValueString = fmt.Sprintf("%v", dummyValue)
	var dummyDataTemplate int

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	strconvAtoiExpected = 1
	strconvAtoi = func(str string) (int, error) {
		strconvAtoiCalled++
		return strconv.Atoi(str)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValueString,
		&dummyDataTemplate,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Int64_Error(t *testing.T) {
	// arrange
	var dummyValue = "some invalid int64"
	var dummyDataTemplate int64

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	strconvParseIntExpected = 1
	strconvParseInt = func(s string, base int, bitSize int) (int64, error) {
		strconvParseIntCalled++
		return strconv.ParseInt(s, base, bitSize)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Int64_NoError(t *testing.T) {
	// arrange
	var dummyValue = rand.Int63()
	var dummyValueString = fmt.Sprintf("%v", dummyValue)
	var dummyDataTemplate int64

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	strconvParseIntExpected = 1
	strconvParseInt = func(s string, base int, bitSize int) (int64, error) {
		strconvParseIntCalled++
		return strconv.ParseInt(s, base, bitSize)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValueString,
		&dummyDataTemplate,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Float64_Error(t *testing.T) {
	// arrange
	var dummyValue = "some invalid float64"
	var dummyDataTemplate float64

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	strconvParseFloatExpected = 1
	strconvParseFloat = func(s string, bitSize int) (float64, error) {
		strconvParseFloatCalled++
		return strconv.ParseFloat(s, bitSize)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Float64_NoError(t *testing.T) {
	// arrange
	var dummyValue = rand.Float64()
	var dummyValueString = fmt.Sprintf("%v", dummyValue)
	var dummyDataTemplate float64

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	strconvParseFloatExpected = 1
	strconvParseFloat = func(s string, bitSize int) (float64, error) {
		strconvParseFloatCalled++
		return strconv.ParseFloat(s, bitSize)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValueString,
		&dummyDataTemplate,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Byte_Error(t *testing.T) {
	// arrange
	var dummyValue = "some invalid byte"
	var dummyDataTemplate byte

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	strconvParseUintExpected = 1
	strconvParseUint = func(s string, base int, bitSize int) (uint64, error) {
		strconvParseUintCalled++
		return strconv.ParseUint(s, base, bitSize)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_Byte_NoError(t *testing.T) {
	// arrange
	var dummyValue = byte(rand.Intn(math.MaxUint8))
	var dummyValueString = fmt.Sprintf("%v", dummyValue)
	var dummyDataTemplate byte

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}
	strconvParseUintExpected = 1
	strconvParseUint = func(s string, base int, bitSize int) (uint64, error) {
		strconvParseUintCalled++
		return strconv.ParseUint(s, base, bitSize)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValueString,
		&dummyDataTemplate,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTtryUnmarshalPrimitiveTypes_OtherTypes(t *testing.T) {
	// arrange
	var dummyValue = "some value"
	var dummyDataTemplate error

	// mock
	createMock(t)

	// expect
	reflectTypeOfExpected = 1
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return reflect.TypeOf(i)
	}

	// SUT + act
	var result = tryUnmarshalPrimitiveTypes(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTryUnmarshal_Primitype(t *testing.T) {
	// arrange
	var dummyValue = "some value"
	var dummyDataTemplate string
	var dummyData = "some data"

	// mock
	createMock(t)

	// expect
	tryUnmarshalPrimitiveTypesFuncExpected = 1
	tryUnmarshalPrimitiveTypesFunc = func(value string, dataTemplate interface{}) bool {
		tryUnmarshalPrimitiveTypesFuncCalled++
		assert.Equal(t, dummyValue, value)
		(*(dataTemplate).(*string)) = dummyData
		return true
	}

	// SUT + act
	var err = tryUnmarshal(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyData, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTryUnmarshal_NoQuoteJSONSuccess_Primitive(t *testing.T) {
	// arrange
	var dummyValue = rand.Int()
	var dummyValueString = strconv.Itoa(dummyValue)
	var dummyDataTemplate int

	// mock
	createMock(t)

	// expect
	tryUnmarshalPrimitiveTypesFuncExpected = 1
	tryUnmarshalPrimitiveTypesFunc = func(value string, dataTemplate interface{}) bool {
		tryUnmarshalPrimitiveTypesFuncCalled++
		assert.Equal(t, dummyValueString, value)
		return false
	}
	jsonUnmarshalExpected = 1
	jsonUnmarshal = func(data []byte, v interface{}) error {
		jsonUnmarshalCalled++
		assert.Equal(t, []byte(dummyValueString), data)
		return json.Unmarshal(data, v)
	}

	// SUT + act
	var err = tryUnmarshal(
		dummyValueString,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyValue, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTryUnmarshal_NoQuoteJSONSuccess_Struct(t *testing.T) {
	// arrange
	var dummyValueString = "{\"foo\":\"bar\",\"test\":123}"
	var dummyDataTemplate struct {
		Foo  string `json:"foo"`
		Test int    `json:"test"`
	}

	// mock
	createMock(t)

	// expect
	tryUnmarshalPrimitiveTypesFuncExpected = 1
	tryUnmarshalPrimitiveTypesFunc = func(value string, dataTemplate interface{}) bool {
		tryUnmarshalPrimitiveTypesFuncCalled++
		assert.Equal(t, dummyValueString, value)
		return false
	}
	jsonUnmarshalExpected = 1
	jsonUnmarshal = func(data []byte, v interface{}) error {
		jsonUnmarshalCalled++
		assert.Equal(t, []byte(dummyValueString), data)
		return json.Unmarshal(data, v)
	}

	// SUT + act
	var err = tryUnmarshal(
		dummyValueString,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "bar", dummyDataTemplate.Foo)
	assert.Equal(t, 123, dummyDataTemplate.Test)

	// verify
	verifyAll(t)
}

func TestTryUnmarshal_WithQuoteJSONSuccess(t *testing.T) {
	// arrange
	var dummyValue = "some value"
	var dummyDataTemplate string

	// mock
	createMock(t)

	// expect
	tryUnmarshalPrimitiveTypesFuncExpected = 1
	tryUnmarshalPrimitiveTypesFunc = func(value string, dataTemplate interface{}) bool {
		tryUnmarshalPrimitiveTypesFuncCalled++
		assert.Equal(t, dummyValue, value)
		return false
	}
	jsonUnmarshalExpected = 2
	jsonUnmarshal = func(data []byte, v interface{}) error {
		jsonUnmarshalCalled++
		if jsonUnmarshalCalled == 1 {
			assert.Equal(t, []byte(dummyValue), data)
		} else if jsonUnmarshalCalled == 2 {
			assert.Equal(t, []byte("\""+dummyValue+"\""), data)
		}
		return json.Unmarshal(data, v)
	}

	// SUT + act
	var err = tryUnmarshal(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyValue, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestTryUnmarshal_Failure(t *testing.T) {
	// arrange
	var dummyValue = "some value"
	var dummyDataTemplate uuid.UUID
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	tryUnmarshalPrimitiveTypesFuncExpected = 1
	tryUnmarshalPrimitiveTypesFunc = func(value string, dataTemplate interface{}) bool {
		tryUnmarshalPrimitiveTypesFuncCalled++
		assert.Equal(t, dummyValue, value)
		return false
	}
	jsonUnmarshalExpected = 2
	jsonUnmarshal = func(data []byte, v interface{}) error {
		jsonUnmarshalCalled++
		if jsonUnmarshalCalled == 1 {
			assert.Equal(t, []byte(dummyValue), data)
		} else if jsonUnmarshalCalled == 2 {
			assert.Equal(t, []byte("\""+dummyValue+"\""), data)
		}
		return json.Unmarshal(data, v)
	}
	fmtErrorfExpected = 1
	fmtErrorf = func(format string, a ...interface{}) error {
		fmtErrorfCalled++
		assert.Equal(t, "Unable to unmarshal value [%v] into data template", format)
		assert.Equal(t, 1, len(a))
		assert.Equal(t, dummyValue, a[0])
		return dummyError
	}

	// SUT + act
	var err = tryUnmarshal(
		dummyValue,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyError, err)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}
