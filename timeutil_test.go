package webserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeNow(t *testing.T) {
	// arrange
	var expectedResult = time.Now()

	// mock
	createMock(t)

	// expect
	timeNowExpected = 1
	timeNow = func() time.Time {
		timeNowCalled++
		return expectedResult
	}

	// SUT + act
	var result = getTimeNow()

	// assert
	assert.Equal(t, expectedResult, result)

	// verify
	verifyAll(t)
}

func TestGetTimeNowUTC(t *testing.T) {
	// arrange
	var dummyResult = time.Now()
	var expectedResult = dummyResult.UTC()

	// mock
	createMock(t)

	// expect
	timeNowExpected = 1
	timeNow = func() time.Time {
		timeNowCalled++
		return expectedResult
	}

	// SUT + act
	var result = getTimeNowUTC()

	// assert
	assert.Equal(t, expectedResult, result)

	// verify
	verifyAll(t)
}

func TestFormatDate(t *testing.T) {
	// arrange
	var dummyTime = time.Date(2345, 6, 7, 8, 9, 10, 11, time.UTC)
	var expectedResult = "2345-06-07"

	// mock
	createMock(t)

	// SUT + act
	var result = formatDate(dummyTime)

	// assert
	assert.Equal(t, expectedResult, result)

	// verify
	verifyAll(t)
}

func TestFormatTime(t *testing.T) {
	// arrange
	var dummyTime = time.Date(2345, 6, 7, 8, 9, 10, 11, time.UTC)
	var expectedResult = "08:09:10"

	// mock
	createMock(t)

	// SUT + act
	var result = formatTime(dummyTime)

	// assert
	assert.Equal(t, expectedResult, result)

	// verify
	verifyAll(t)
}

func TestFormatDateTime(t *testing.T) {
	// arrange
	var dummyTime = time.Date(2345, 6, 7, 8, 9, 10, 11, time.UTC)
	var expectedResult = "2345-06-07T08:09:10"

	// mock
	createMock(t)

	// SUT + act
	var result = formatDateTime(dummyTime)

	// assert
	assert.Equal(t, expectedResult, result)

	// verify
	verifyAll(t)
}
