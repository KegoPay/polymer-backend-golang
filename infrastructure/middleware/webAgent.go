package middlewares

import (
	"github.com/gin-gonic/gin"
	"kego.com/application/interfaces"
	"kego.com/application/middlewares"
)

func WebAgentMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appContext, next := middlewares.WebAgentMiddleware(&interfaces.ApplicationContext[any]{
			Ctx:    ctx,
			Keys:   ctx.Keys,
			Header: ctx.Request.Header,
		}, ctx.ClientIP())
		if next {
			ctx.Set("AppContext", appContext)
			ctx.Next()
		}
	}
}
