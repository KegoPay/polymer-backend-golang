package metrics

import (
	"context"
	"os"

	apitoolkit "github.com/apitoolkit/apitoolkit-go"
)


type APIToolKitMonitor struct {
	Client *apitoolkit.Client
}

func (toolKit *APIToolKitMonitor) Init(){
	tkInstance, err := apitoolkit.NewClient(context.Background(),
		apitoolkit.Config{
			RedactHeaders: []string{"Authorization", "Cookies"},
			RedactRequestBody: []string{"$.password", "$.transactionPin", "$.bvn", "$.otp"},
			RedactResponseBody: []string{"$.body.token"},
			APIKey: os.Getenv("APITOOLKIT_API_KEY")})
	if err != nil {
		panic(err)
	}
	toolKit.Client = tkInstance
}

func (toolkit *APIToolKitMonitor) MetricMiddleware()  any {
	return toolkit.Client.GinMiddleware
}

func (toolkit *APIToolKitMonitor) ReportError(ctx any,err error) {
	apitoolkit.ReportError(ctx.(context.Context), err)
}
