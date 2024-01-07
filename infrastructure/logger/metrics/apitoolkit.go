package metrics

import (
	"context"
	"net/http"
	"os"

	apitoolkit "github.com/apitoolkit/apitoolkit-go"
)


type APIToolKitMonitor struct {
	Client *apitoolkit.Client
}

func (toolKit *APIToolKitMonitor) Init(){
	tkInstance, err := apitoolkit.NewClient(context.Background(),
		apitoolkit.Config{
			RedactHeaders: []string{"Authorization", "Cookies", "X-API-KEY", "X-Api-Key", "x-api-key"},
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

func (toolkit *APIToolKitMonitor) GetRoundTripper(ctx context.Context) http.RoundTripper {
	return toolkit.Client.WrapRoundTripper(ctx, http.DefaultTransport, apitoolkit.WithRedactHeaders("X-API-KEY", "X-Api-Key", "x-api-key"))
}

func (toolkit *APIToolKitMonitor) LogMessage(key string, err error) {
	toolkit.Client.PublishMessage(context.Background(), apitoolkit.Payload{
		Tags: []string{key},
	})
}


func (toolkit *APIToolKitMonitor) CleanUp() error {
	return toolkit.Client.Close()
}
