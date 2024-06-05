package middlewares

import (
	"github.com/gin-gonic/gin"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/middlewares"
)

func InhouseAuthMiddleware(signingKey string, issuer string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appContext, next := middlewares.InhouseAuthMiddleware(&interfaces.ApplicationContext[any]{
			Ctx:    ctx,
			Header: ctx.Request.Header,
		}, signingKey, issuer)
		if next {
			ctx.Set("AppContext", appContext)
			ctx.Next()
		}
	}
}
