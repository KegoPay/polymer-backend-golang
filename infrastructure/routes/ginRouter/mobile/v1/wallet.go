package routev1

import (
	"github.com/gin-gonic/gin"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	middlewares "kego.com/infrastructure/middleware"
)

func WalletRouter(router *gin.RouterGroup) {
	walletRouter := router.Group("/wallet")
	{
		walletRouter.POST("/:businessID/payment/international/send", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SendPaymentDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			body.IPAddress = ctx.ClientIP()
			appContext := interfaces.ApplicationContext[dto.SendPaymentDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.InitiateBusinessInternationalPayment(&appContext)
		})

		walletRouter.POST("/payment/international/send", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SendPaymentDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			body.IPAddress = ctx.ClientIP()
			appContext := interfaces.ApplicationContext[dto.SendPaymentDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.InitiatePersonalInternationalPayment(&appContext)
		})

		walletRouter.POST("/:businessID/payment/local/send", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SendPaymentDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			body.IPAddress = ctx.ClientIP()
			appContext := interfaces.ApplicationContext[dto.SendPaymentDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.InitiateBusinessLocalPayment(&appContext)
		})

		walletRouter.POST("/:businessID/payment/local/fee", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SendPaymentDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			body.IPAddress = ctx.ClientIP()
			appContext := interfaces.ApplicationContext[dto.SendPaymentDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.BusinessLocalPaymentFee(&appContext)
		})

		walletRouter.POST("/:businessID/payment/international/fee", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SendPaymentDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			body.IPAddress = ctx.ClientIP()
			appContext := interfaces.ApplicationContext[dto.SendPaymentDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.BusinessInternationalPaymentFee(&appContext)
		})
		
		walletRouter.POST("/:businessID/payment/local/verify-name", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.NameVerificationDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			appContext := interfaces.ApplicationContext[dto.NameVerificationDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.VerifyLocalAccountName(&appContext)
		})
	}
}
