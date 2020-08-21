package webserver

// LogLevel is the severity level of logging
type LogLevel int

// These are the enum definitions of log types and presets
const (
	Debug LogLevel = 0
	Info  LogLevel = iota
	Warn
	Error
	Fatal
	maxLogLevel
)

// These are the string representations of log category and preset names
const (
	debugName string = "Debug"
	infoName  string = "Info"
	warnName  string = "Warn"
	errorName string = "Error"
	fatalName string = "Fatal"
)

var supportedLogLevels = map[LogLevel]string{
	Debug: debugName,
	Info:  infoName,
	Warn:  warnName,
	Error: errorName,
	Fatal: fatalName,
}

var logLevelNameMapping = map[string]LogLevel{
	debugName: Debug,
	infoName:  Info,
	warnName:  Warn,
	errorName: Error,
	fatalName: Fatal,
}

// FromString converts a LogLevel flag instance to its string representation
func (logLevel LogLevel) String() string {
	for key, value := range supportedLogLevels {
		if logLevel == key {
			return value
		}
	}
	return debugName
}

// NewLogLevel converts a string representation of LogLevel flag to its strongly typed instance
func NewLogLevel(value string) LogLevel {
	var logLevel, found = logLevelNameMapping[value]
	if !found {
		return Debug
	}
	return logLevel
}
