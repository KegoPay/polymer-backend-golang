package webroutev1

import (
	"github.com/gin-gonic/gin"
	"usepolymer.co/application/controllers"
	"usepolymer.co/application/interfaces"
	middlewares "usepolymer.co/infrastructure/middleware"
)

func BusinessRouter(router *gin.RouterGroup) {
	businessRouter := router.Group("/business")
	businessRouter.Use(middlewares.UserAgentMiddleware(false))
	{
		businessRouter.POST("/verify/manual", middlewares.AttestationMiddleware(), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			appContext := interfaces.ApplicationContext[any]{
				Ctx: appContextAny.Ctx,
				Query: map[string]any{
					"token": ctx.Query("token"),
				},
			}
			controllers.VerifyBusinessManual(&appContext)
		})
	}
}
