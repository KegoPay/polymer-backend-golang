package webroutev1

import (
	"github.com/gin-gonic/gin"
)


func WebhookRouter(router *gin.RouterGroup) {
	// webhookRouter := router.Group("/webhook")
	{
		// webhookRouter.POST("/flutterwave/transfer", middlewares.WebAgentMiddleware(), func(ctx *gin.Context) {
		// 	appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
		// 	var body dto.FlutterwaveWebhookDTO
		// 	if ctx.Request.Header.Get("Content-Type") == "application/json" {
		// 		if err := ctx.BindJSON(&body); err != nil {
		// 			apperrors.ErrorProcessingPayload(ctx, nil)
		// 			return
		// 		}
		// 	}else {
		// 		if err := ctx.Bind(&body); err != nil {
		// 			apperrors.ErrorProcessingPayload(ctx, nil)
		// 			return
		// 		}
		// 	}
		// 	appContext := interfaces.ApplicationContext[dto.FlutterwaveWebhookDTO]{
		// 		Keys: appContextAny.Keys,
		// 		Ctx: appContextAny.Ctx,
		// 		Body: &body,
		// 	}
		// 	controllers.FlutterwaveWebhook(&appContext)
		// })
	}
}
