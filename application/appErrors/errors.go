package apperrors

import (
	"fmt"
	"net/http"

	"kego.com/infrastructure/logger"
	server_response "kego.com/infrastructure/serverResponse"
)

func NotFoundError(ctx interface{}, message string, device_id *string){
	server_response.Responder.Respond(ctx, http.StatusNotFound, message, nil, nil, nil,  device_id)
}

func ValidationFailedError(ctx interface{}, errMessages *[]error, device_id *string){
	server_response.Responder.Respond(ctx, http.StatusUnprocessableEntity, "Payload validation failed 🙄", nil, *errMessages, nil, device_id)
}

func EntityAlreadyExistsError(ctx interface{}, message string, device_id *string){
	server_response.Responder.Respond(ctx, http.StatusConflict, message, nil, nil, nil, device_id)
}

func AuthenticationError(ctx interface{}, message string, device_id *string){
	server_response.Responder.Respond(ctx, http.StatusUnauthorized, message, nil, nil, nil, device_id)
}

func ExternalDependencyError(ctx interface{}, serviceName string, statusCode string, err error, device_id *string) {
	logger.Error(err, logger.LoggerOptions{
		Key: fmt.Sprintf("error with %s. status code %s", serviceName, statusCode),
	})
	// logger.MetricMonitor.ReportError(fmt.Errorf(fmt.Sprintf("error with %s", serviceName)), []logger.LoggerOptions{
	// 	{
	// 		Key: "statusCode",
	// 		Data: statusCode,
	// 	},
	// })
	// logger.MetricMonitor.ReportError(err, nil)
	server_response.Responder.Respond(ctx, http.StatusServiceUnavailable,
		"Omo! Our service is temporarily down 😢. Our team is working to fix it. Please check back later.", nil, nil, nil, device_id)
}

func ErrorProcessingPayload(ctx interface{}, device_id *string){
	server_response.Responder.Respond(ctx, http.StatusBadRequest, "Abnormal payload passed 🤨", nil, nil, nil, device_id)
}

func FatalServerError(ctx interface{}, err error, device_id *string){
	// logger.MetricMonitor.ReportError(err, nil)
	server_response.Responder.Respond(ctx, http.StatusInternalServerError,
		"Omo! Our service is temporarily down 😢. Our team is working to fix it. Please check back later.", nil, nil, nil, device_id)
}

func UnknownError(ctx interface{}, err error, device_id *string){
	// logger.MetricMonitor.ReportError(err, nil)
	server_response.Responder.Respond(ctx, http.StatusBadRequest,
		"Omo! Something went wrong somewhere 😭. Please check back later.", nil, nil, nil, device_id)
}

func CustomError(ctx interface{}, msg string, device_id *string){
	server_response.Responder.Respond(ctx, http.StatusBadRequest, msg, nil, nil, nil, device_id)
}

func UnsupportedAppVersion(ctx interface{}, device_id *string){
	server_response.Responder.Respond(ctx, http.StatusBadRequest,
		"Uh oh! Seems you're using an old version of the app. 🤦🏻‍♂️\n Upgrade to the latest version to continue enjoying our blazing fast services! 🚀", nil, nil, nil, device_id)
}

func UnsupportedUserAgent(ctx interface{}, device_id *string){
	// logger.MetricMonitor.ReportError(errors.New("unspported user agent"), []logger.LoggerOptions{
	// 	{Key: "ctx",
	// 	Data: ctx,},
	// })
	server_response.Responder.Respond(ctx, http.StatusBadRequest,
		"Unsupported user agent 👮🏻‍♂️", nil, nil, nil, device_id)
}

func ClientError(ctx interface{}, msg string, errs []error, response_code *uint, device_id *string){
	server_response.Responder.Respond(ctx, http.StatusBadRequest, msg, nil, errs, response_code, device_id)
}
