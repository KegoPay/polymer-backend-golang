package middlewares

import (
	"github.com/gin-gonic/gin"
	"kego.com/application/interfaces"
	"kego.com/application/middlewares"
)

func AuthenticationMiddleware(business_route bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appContext, next := middlewares.AuthenticationMiddleware(&interfaces.ApplicationContext[any]{
			Ctx:    ctx,
			Keys:   ctx.Keys,
			Header: ctx.Request.Header,
		})
		if next {
			ctx.Set("AppContext", appContext)
			ctx.Next()
		}
	}
}
