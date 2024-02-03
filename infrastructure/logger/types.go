package logger



type MetricType interface {
	MetricMiddleware() any
	ReportError(err error, data []LoggerOptions)
	Log(msg string, data []LoggerOptions, logLevel LogLevel)
	Init() 
	CleanUp() error
}

type LogLevel string

var (
	ErrorLevel LogLevel = "error"
	InfoLevel LogLevel = "info"
	WarningLevel LogLevel = "warning"
)