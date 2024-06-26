package middlewares

import (
	"github.com/gin-gonic/gin"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/middlewares"
)

func UserAgentMiddleware(mobileOnly bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appContext, next := middlewares.UserAgentMiddleware(&interfaces.ApplicationContext[any]{
			Ctx:    ctx,
			Keys:   ctx.Keys,
			Header: ctx.Request.Header,
		}, "0.0.1", ctx.ClientIP())
		if next {
			ctx.Set("AppContext", appContext)
			ctx.Next()
		}
	}
}
