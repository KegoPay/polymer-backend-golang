package webroutev1

import (
	"github.com/gin-gonic/gin"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	middlewares "kego.com/infrastructure/middleware"
)


func WebhookRouter(router *gin.RouterGroup) {
	webhookRouter := router.Group("/webhook")
	{
		webhookRouter.POST("/flutterwave/transfer", middlewares.WebAgentMiddleware(), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.FlutterwaveWebhookDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			appContext := interfaces.ApplicationContext[dto.FlutterwaveWebhookDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.FlutterwaveWebhook(&appContext)
		})
	}
}
