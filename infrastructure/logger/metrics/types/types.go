package types

type MetricType interface {
	MetricMiddleware() any
	ReportError(any, error)
	Init() 
}