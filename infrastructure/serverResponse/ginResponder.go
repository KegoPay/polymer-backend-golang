package server_response

import (
	"errors"

	"github.com/gin-gonic/gin"
	"kego.com/infrastructure/logger"
)

type ginResponder struct{}

func (gr ginResponder)Respond(ctx interface{}, code int, message string, payload interface{}, errs []error) {
	ginCtx, ok := (ctx).(*gin.Context)
    if !ok {
		logger.Error(errors.New("could not transform *interface{} to gin.Context in serverResponse package"), logger.LoggerOptions{
			Key: "payload",
			Data: ctx,
		})
        return
    }
	ginCtx.Abort()
	ginCtx.JSON(code, gin.H{
		"message": message,
		"body":    payload,
		"errors": func() interface{} {
			if len(errs) == 0 {
				return nil
			}
			errMsgs := []string{}
			for _, err := range errs {
				errMsgs = append(errMsgs, err.Error())
			}
			return errMsgs
		}(),
	})
}
