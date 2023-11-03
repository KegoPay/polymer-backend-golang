package apperrors

import (
	"fmt"
	"net/http"

	"kego.com/infrastructure/logger"
	server_response "kego.com/infrastructure/serverResponse"
)

func NotFoundError(ctx interface{}, message string){
	server_response.Responder.Respond(ctx, http.StatusNotFound, message, nil, nil)
}

func ValidationFailedError(ctx interface{}, errMessages *[]error){
	server_response.Responder.Respond(ctx, http.StatusUnprocessableEntity, "payload validation failed", nil, *errMessages)
}

func EntityAlreadyExistsError(ctx interface{}, message string){
	server_response.Responder.Respond(ctx, http. StatusConflict, message, nil, nil)
}

func ExternalDependencyError(ctx interface{}, serviceName string, statusCode string, err error) {
	logger.Error(err, logger.LoggerOptions{
		Key: fmt.Sprintf("error with %s. status code %s", serviceName, statusCode),
	})
	server_response.Responder.Respond(ctx, http.StatusServiceUnavailable,
		"Oops! Our service is temporarily down. Our team working to fix it. Please check back later.", nil, nil)
}

func ErrorProcessingPayload(ctx interface{}){
	server_response.Responder.Respond(ctx, http.StatusBadRequest, "abnormal payload passed", nil, nil)
}