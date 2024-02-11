package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOptions struct{
	Key string
	Data interface{}
}

var MetricMonitor MetricType = (&SentryMonitor{})
var RequestMetricMonitor  = (&APIToolKitMonitor{})


// This logs info level messages.
func Info(msg string, payload ...LoggerOptions) {
	zapFields := []zapcore.Field{}
	for _, data := range payload{
		zapFields = append(zapFields, zap.Any(data.Key, data.Data))
	}
	MetricMonitor.Log(msg, payload, InfoLevel)
	Logger.Info(msg, zapFields...)
}

// This logs error messages.
func Error(err error, payload ...LoggerOptions) {
	zapFields := []zapcore.Field{}
	for _, data := range payload{
		zapFields = append(zapFields, zap.Any(data.Key, data.Data))
	}
	MetricMonitor.ReportError(err, payload)
	Logger.Error(err.Error(), zapFields...)
}

// This logs warning messages.
func Warning(msg string, payload ...LoggerOptions) {
	zapFields := []zapcore.Field{}
	for _, data := range payload{
		zapFields = append(zapFields, zap.Any(data.Key, data.Data))
	}
	MetricMonitor.Log(msg, payload, InfoLevel)
	Logger.Warn(msg, zapFields...)
}
