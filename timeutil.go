package webserver

import "time"

const (
	dateFormat     string = "2006-01-02"
	timeFormat     string = "15:04:05"
	dateTimeFormat string = dateFormat + "T" + timeFormat
)

// getTimeNow returns the current system time
func getTimeNow() time.Time {
	return timeNow()
}

// getTimeNowUTC returns the UTC representation of the current system time
func getTimeNowUTC() time.Time {
	return timeNow().UTC()
}

// formatDate returns the time in string format "yyyy-MM-dd"
func formatDate(value time.Time) string {
	return value.Format(dateFormat)
}

// formatTime returns the time in string format "HH:mm:ss"
func formatTime(value time.Time) string {
	return value.Format(timeFormat)
}

// formatDateTime returns the time in string format "yyyy-MM-ddTHH:mm:ss"
func formatDateTime(value time.Time) string {
	return value.Format(dateTimeFormat)
}
