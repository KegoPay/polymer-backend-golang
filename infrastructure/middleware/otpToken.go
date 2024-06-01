package middlewares

import (
	"github.com/gin-gonic/gin"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/middlewares"
)

func OTPTokenMiddleware(intent string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appContext, next := middlewares.OTPTokenMiddleware(&interfaces.ApplicationContext[any]{
			Ctx:    ctx,
			Keys:   ctx.Keys,
			Header: ctx.Request.Header,
		}, ctx.ClientIP(), intent)
		if next {
			ctx.Set("AppContext", appContext)
			ctx.Next()
		}
	}
}
