package routev1

import (
	"github.com/gin-gonic/gin"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/utils"
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
			body.DeviceID = ctx.GetHeader("polymer-device-id")
			body.UserAgent = ctx.Request.UserAgent()
			appVersion := utils.ExtractAppVersionFromUserAgentHeader(ctx.Request.UserAgent())
			if appVersion == nil {
				apperrors.UnsupportedAppVersion(ctx)
				return
			}
			body.AppVersion = *appVersion
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
			deviceID := ctx.GetHeader("polymer-device-id")
			if deviceID == "" {
				apperrors.AuthenticationError(ctx, "no client id")
				return
			}
			body.DeviceID = deviceID
			controllers.LoginUser(&interfaces.ApplicationContext[dto.LoginDTO]{
				Ctx: ctx,
				Body: &body,
				Header: ctx.Request.Header,
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

		authRouter.POST("/email/verify", func(ctx *gin.Context) {
			var body dto.VerifyEmailData
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			controllers.VerifyEmail(&interfaces.ApplicationContext[dto.VerifyEmailData]{
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

		authRouter.POST("/account/verify", middlewares.WebAgentMiddleware(), middlewares.AuthenticationMiddleware(false, false) ,func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.VerifyAccountData
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			controllers.VerifyAccount(&interfaces.ApplicationContext[dto.VerifyAccountData]{
				Ctx: ctx,
				Body: &body,
				Keys: appContextAny.Keys,
			})
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

		authRouter.POST("/account/password/update", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
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

		authRouter.POST("/account/deactivate", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.ConfirmPin
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			appContext := interfaces.ApplicationContext[dto.ConfirmPin]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.DeactivateAccount(&appContext)
		})

		authRouter.POST("/file/generate-url", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.FileUploadOptions
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			appContext := interfaces.ApplicationContext[dto.FileUploadOptions]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.GenerateFileURL(&appContext)
		})

		authRouter.POST("/account/transaction-pin/set", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SetTransactionPinDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			appContext := interfaces.ApplicationContext[dto.SetTransactionPinDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.SetTransactionPin(&appContext)
		})
	}
}
