package logger

import (
	"os"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
)

type SentryMonitor struct {
}

func (sm *SentryMonitor) Init() {
	 _ = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		Debug: true,
		AttachStacktrace: true,
		EnableTracing: true,
		TracesSampleRate: 1.0,
		ProfilesSampleRate: 1.0,
	})
}

func (sm *SentryMonitor) MetricMiddleware() any {
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}

func (sm *SentryMonitor) ReportError(err error, logData []LoggerOptions) {
	sm.Log(err.Error(), logData, ErrorLevel)
	sentry.CaptureException(err)
}

func (sm *SentryMonitor) Log(msg string, logData []LoggerOptions, level LogLevel) {
	var sentryLevel sentry.Level
	if level == ErrorLevel {
		sentryLevel = sentry.LevelError
	}else if level == WarningLevel{
		sentryLevel = sentry.LevelWarning
	}else {
		sentryLevel = sentry.LevelInfo
	}
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Message: msg,
		Data: func () map[string]any {
			data := map[string]any{}
			for _, d := range logData {
				data[d.Key] = d.Data
			}
			return data
		}(),
		Level: sentryLevel,
	})
}

func (sm *SentryMonitor) CleanUp() error {
	return nil
}