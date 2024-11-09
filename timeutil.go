package webserver

import "time"

const (
	dateFormat     string = "2006-01-02"
	timeFormat     string = "15:04:05.999"
	dateTimeFormat string = dateFormat + "T" + timeFormat
)

// getTimeNowUTC returns the current time in UTC format
func getTimeNowUTC() time.Time {
	return time.Now().UTC()
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
