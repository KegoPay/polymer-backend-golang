package authroutev1

import (
	"github.com/gin-gonic/gin"
	"kego.com/application/controllers"
	"kego.com/application/interfaces"
)


func UserRouter(router *gin.RouterGroup) {
	infoRouter := router.Group("/user")
	{
		infoRouter.GET("/profile", func(ctx *gin.Context) {
			controllers.FetchUserProfile(&interfaces.ApplicationContext[any]{
				Keys: ctx.Keys,
				Ctx: ctx,
			})
		})
	}
}
