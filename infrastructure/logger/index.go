package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOptions struct{
	Key string
	Data interface{}
}

// This logs info level messages.
//
// Only the first parameter will be accepted.
func Info(msg string, payload ...LoggerOptions) {
	zapFields := []zapcore.Field{}
	for _, data := range payload{
		zapFields = append(zapFields, zap.Any(data.Key, data.Data))
	}
	Logger.Info(msg, zapFields...)
}

// This logs error messages.
//
// Only the first parameter will be accepted.
func Error(err error, payload ...LoggerOptions) {
	zapFields := []zapcore.Field{}
	for _, data := range payload{
		zapFields = append(zapFields, zap.Any(data.Key, data.Data))
	}
	Logger.Error(err.Error(), zapFields...)
}
