package logger

import (
	"go.uber.org/zap"
)

type LoggerOptions struct{
	Key string
	Data interface{}
}

// This logs info level messages.
//
// Only the first parameter will be accepted.
func Info(msg string, payload ...LoggerOptions) {
	if (len(payload) > 0){
		Logger.Info(msg, zap.Any(payload[0].Key, payload[0].Data))
	}else {
		Logger.Info(msg)
	}
}

// This logs error messages.
//
// Only the first parameter will be accepted.
func Error(err error, payload ...LoggerOptions) {
	if (len(payload) > 0){
		Logger.Error(err.Error(), zap.Any(payload[0].Key, payload[0].Data))
	}else {
		Logger.Error(err.Error())
	}
}
