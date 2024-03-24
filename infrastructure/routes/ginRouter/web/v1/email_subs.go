package webroutev1

import (
	"github.com/gin-gonic/gin"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	middlewares "kego.com/infrastructure/middleware"
)


func EmailSubsRouter(router *gin.RouterGroup) {
	emailSubRouter := router.Group("/emailsub")
	emailSubRouter.Use(middlewares.UserAgentMiddleware(false))
	{
		emailSubRouter.POST("/subscribe", middlewares.WebAgentMiddleware(), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.EmailSubscriptionDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, nil)
				return
			}
			appContext := interfaces.ApplicationContext[dto.EmailSubscriptionDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.EmailSubscription(&appContext)
		})
	}
}
