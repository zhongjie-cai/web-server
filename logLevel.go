package webserver

// LogLevel is the severity level of logging
type LogLevel int

// These are the enum definitions of log types and presets
const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	maxLogLevel
)

// These are the string representations of log category and preset names
const (
	debugLogLevelName string = "Debug"
	infoLogLevelName  string = "Info"
	warnLogLevelName  string = "Warn"
	errorLogLevelName string = "Error"
	fatalLogLevelName string = "Fatal"
)

var supportedLogLevels = map[LogLevel]string{
	LogLevelDebug: debugLogLevelName,
	LogLevelInfo:  infoLogLevelName,
	LogLevelWarn:  warnLogLevelName,
	LogLevelError: errorLogLevelName,
	LogLevelFatal: fatalLogLevelName,
}

var logLevelNameMapping = map[string]LogLevel{
	debugLogLevelName: LogLevelDebug,
	infoLogLevelName:  LogLevelInfo,
	warnLogLevelName:  LogLevelWarn,
	errorLogLevelName: LogLevelError,
	fatalLogLevelName: LogLevelFatal,
}

// FromString converts a LogLevel flag instance to its string representation
func (logLevel LogLevel) String() string {
	for key, value := range supportedLogLevels {
		if logLevel == key {
			return value
		}
	}
	return debugLogLevelName
}

// NewLogLevel converts a string representation of LogLevel flag to its strongly typed instance
func NewLogLevel(value string) LogLevel {
	var logLevel, found = logLevelNameMapping[value]
	if !found {
		return LogLevelDebug
	}
	return logLevel
}
