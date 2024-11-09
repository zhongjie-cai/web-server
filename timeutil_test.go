package webserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker/v2"
)

func TestGetTimeNowUTC(t *testing.T) {
	// arrange
	var dummyResult = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	var expectedResult = dummyResult.UTC()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(time.Now).Expects().Returns(expectedResult).Once()

	// SUT + act
	var result = getTimeNowUTC()

	// assert
	assert.Equal(t, expectedResult, result)
}

func TestFormatDate(t *testing.T) {
	// arrange
	var dummyTime = time.Date(2345, 6, 7, 8, 9, 10, 11, time.UTC)
	var expectedResult = "2345-06-07"

	// SUT + act
	var result = formatDate(dummyTime)

	// assert
	assert.Equal(t, expectedResult, result)
}

func TestFormatTime(t *testing.T) {
	// arrange
	var dummyTime = time.Date(2345, 6, 7, 8, 9, 10, 11, time.UTC)
	var expectedResult = "08:09:10"

	// SUT + act
	var result = formatTime(dummyTime)

	// assert
	assert.Equal(t, expectedResult, result)
}

func TestFormatDateTime(t *testing.T) {
	// arrange
	var dummyTime = time.Date(2345, 6, 7, 8, 9, 10, 11, time.UTC)
	var expectedResult = "2345-06-07T08:09:10"

	// SUT + act
	var result = formatDateTime(dummyTime)

	// assert
	assert.Equal(t, expectedResult, result)
}
