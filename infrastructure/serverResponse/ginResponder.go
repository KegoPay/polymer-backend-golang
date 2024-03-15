package server_response

import (
	"errors"

	"github.com/gin-gonic/gin"
	"kego.com/infrastructure/logger"
)

type ginResponder struct{}

// Sends an encrypted payload to the client
func (gr ginResponder)Respond(ctx interface{}, code int, message string, payload interface{}, errs []error, response_code *uint) {
	ginCtx, ok := (ctx).(*gin.Context)
    if !ok {
		logger.Error(errors.New("could not transform *interface{} to gin.Context in serverResponse package"), logger.LoggerOptions{
			Key: "payload",
			Data: ctx,
		})
        return
    }
	ginCtx.Abort()
	response := map[string]any{
		"message": message,
		"body":    payload,
	}
	if response_code != nil{
		response["response_code"] = response_code
	}
	if errs != nil{
		errMsgs := []string{}
		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error())
		}
		response["errors"] = errMsgs
	}
	ginCtx.JSON(code, response)
}

// Sends a response to the client using plain JSON
func (gr ginResponder) UnEncryptedRespond(ctx interface{}, code int, message string, payload interface{}, errs []error, response_code *uint) {
	ginCtx, ok := (ctx).(*gin.Context)
    if !ok {
		logger.Error(errors.New("could not transform *interface{} to gin.Context in serverResponse package"), logger.LoggerOptions{
			Key: "payload",
			Data: ctx,
		})
        return
    }
	ginCtx.Abort()
	response := map[string]any{
		"message": message,
		"body":    payload,
	}
	if response_code != nil{
		response["response_code"] = response_code
	}
	if errs != nil{
		errMsgs := []string{}
		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error())
		}
		response["errors"] = errMsgs
	}
	ginCtx.JSON(code, response)
}
