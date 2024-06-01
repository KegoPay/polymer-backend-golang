package middlewares

import (
	"github.com/gin-gonic/gin"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/middlewares"
)

func AuthenticationMiddleware(business_route bool, restricted bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appContext, next := middlewares.AuthenticationMiddleware(&interfaces.ApplicationContext[any]{
			Ctx:    ctx,
			Keys:   ctx.Keys,
			Header: ctx.Request.Header,
			Param: map[string]any{
				"businessID": ctx.Param("businessID"),
			},
		}, restricted, business_route)
		if next {
			ctx.Set("AppContext", appContext)
			ctx.Next()
		}
	}
}
