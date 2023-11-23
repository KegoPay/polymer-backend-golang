package authroutev1

import (
	"github.com/gin-gonic/gin"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/entities"
	middlewares "kego.com/infrastructure/middleware"
)


func AuthRouter(router *gin.RouterGroup) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/account/create", func(ctx *gin.Context) {
			var body dto.CreateAccountDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			body.DeviceID = ctx.GetHeader("KEGO_DEVICE_ID")
			body.DeviceType = entities.DeviceType(ctx.GetHeader("KEGO_DEVICE_TYPE"))
			controllers.CreateAccount(&interfaces.ApplicationContext[dto.CreateAccountDTO]{
				Ctx: ctx,
				Body: &body,
			})
		})

		authRouter.POST("/account/login", func(ctx *gin.Context) {
			var body dto.LoginDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			controllers.LoginUser(&interfaces.ApplicationContext[dto.LoginDTO]{
				Ctx: ctx,
				Body: &body,
			})
		})

		authRouter.GET("/otp/resend", func(ctx *gin.Context) {
			query := map[string]any{
				"email": ctx.Query("email"),
			}
			controllers.ResendOTP(&interfaces.ApplicationContext[any]{
				Ctx: ctx,
				Query: query,
			})
		})

		authRouter.POST("/account/verify", func(ctx *gin.Context) {
			var body dto.VerifyData
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			controllers.VerifyAccount(&interfaces.ApplicationContext[dto.VerifyData]{
				Ctx: ctx,
				Body: &body,
			})
		})

		authRouter.GET("/account/exits",  func(ctx *gin.Context) {
			query := map[string]any{
				"email": ctx.Query("email"),
			}
			controllers.AccountWithEmailExists(&interfaces.ApplicationContext[any]{
				Ctx: ctx,
				Query: query,
			})
		})

		authRouter.GET("/account/kyc/retry", middlewares.AuthenticationMiddleware(false), func(ctx *gin.Context) {
			appContext, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			controllers.RetryIdentityVerification(appContext)
		})

		authRouter.POST("/account/password/reset",  func(ctx *gin.Context) {
			var body dto.ResetPasswordDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			controllers.ResetPassword(&interfaces.ApplicationContext[dto.ResetPasswordDTO]{
				Ctx: ctx,
				Body: &body,
			})
		})

		authRouter.POST("/account/password/update", middlewares.AuthenticationMiddleware(false), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.UpdatePassword
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			appContext := interfaces.ApplicationContext[dto.UpdatePassword]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.UpdatePassword(&appContext)
		})
	}
}
