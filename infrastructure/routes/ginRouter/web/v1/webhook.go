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
			if ctx.Request.Header.Get("Content-Type") == "application/json" {
				if err := ctx.BindJSON(&body); err != nil {
					apperrors.ErrorProcessingPayload(ctx)
					return
				}
				var transfer dto.FlutterwaveWebhookTransfer
				if err := ctx.BindJSON(&transfer); err != nil {
					apperrors.ErrorProcessingPayload(ctx)
					return
				}
				body.Transfer = &transfer
			}else {
				if err := ctx.Bind(&body); err != nil {
					apperrors.ErrorProcessingPayload(ctx)
					return
				}
				var customer dto.FlutterwaveWebhookCustomer
				if err := ctx.Bind(&customer); err != nil {
					apperrors.ErrorProcessingPayload(ctx)
					return
				}
				body.Customer = &customer
				var entity dto.FlutterwaveWebhookEntity
				if err := ctx.Bind(&entity); err != nil {
					apperrors.ErrorProcessingPayload(ctx)
					return
				}
				body.Entity = &entity
			}
			appContext := interfaces.ApplicationContext[dto.FlutterwaveWebhookDTO]{
				Keys: appContextAny.Keys,
				Ctx: appContextAny.Ctx,
				Body: &body,
			}
			controllers.FlutterwaveWebhook(&appContext)
		})
	}
}
