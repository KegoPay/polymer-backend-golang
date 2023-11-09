package authroutev1

import (
	"github.com/gin-gonic/gin"
	"kego.com/application/controllers"
	"kego.com/application/interfaces"
	middlewares "kego.com/infrastructure/middleware"
)


func UserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/profile", middlewares.AuthenticationMiddleware(false), func(ctx *gin.Context) {
			appContext, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			controllers.FetchUserProfile(appContext)
		})
	}
}
