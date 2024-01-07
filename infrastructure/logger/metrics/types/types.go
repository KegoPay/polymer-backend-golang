package types

import (
	"context"
	"net/http"
)

type MetricType interface {
	MetricMiddleware() any
	ReportError(any, error)
	GetRoundTripper(ctx context.Context) http.RoundTripper
	Init() 
	CleanUp() error
}