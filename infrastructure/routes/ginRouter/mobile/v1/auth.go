package routev1

import (
	"github.com/gin-gonic/gin"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/controllers"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/utils"
	"usepolymer.co/infrastructure/logger"
	middlewares "usepolymer.co/infrastructure/middleware"
)

func AuthRouter(router *gin.RouterGroup) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/account/create", middlewares.AttestationMiddleware(), func(ctx *gin.Context) {
			var body dto.CreateAccountDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			body.DeviceID = ctx.GetHeader("polymer-device-id")
			body.UserAgent = ctx.Request.UserAgent()
			body.PushNotificationToken = ctx.GetHeader("x-firebase-push-token")
			appVersion := utils.ExtractAppVersionFromUserAgentHeader(ctx.Request.UserAgent())
			if appVersion == nil {
				apperrors.UnsupportedAppVersion(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			body.AppVersion = *appVersion
			controllers.CreateAccount(&interfaces.ApplicationContext[dto.CreateAccountDTO]{
				Ctx:  ctx,
				Body: &body,
			})
		})

		authRouter.POST("/account/login", middlewares.AttestationMiddleware(), func(ctx *gin.Context) {
			var body dto.LoginDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			deviceID := ctx.GetHeader("polymer-device-id")
			if deviceID == "" {
				apperrors.AuthenticationError(ctx, "no client id", utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			body.DeviceID = deviceID
			pushNotificationToken := ctx.GetHeader("x-firebase-push-token")
			if pushNotificationToken == "" {
				apperrors.AuthenticationError(ctx, "no push notification token", utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			body.PushNotificationToken = pushNotificationToken
			controllers.LoginUser(&interfaces.ApplicationContext[dto.LoginDTO]{
				Ctx:    ctx,
				Body:   &body,
				Header: ctx.Request.Header,
			})
		})

		authRouter.POST("/otp/resend", func(ctx *gin.Context) {
			var body dto.ResendOTP
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			controllers.ResendOTP(&interfaces.ApplicationContext[dto.ResendOTP]{
				Ctx:    ctx,
				Body:   &body,
				Header: ctx.Request.Header,
			})
		})

		authRouter.POST("/otp/verify", func(ctx *gin.Context) {
			var body dto.VerifyOTPDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			controllers.VerifyOTP(&interfaces.ApplicationContext[dto.VerifyOTPDTO]{
				Ctx:    ctx,
				Body:   &body,
				Header: ctx.Request.Header,
			})
		})

		authRouter.PATCH("/email/verify", middlewares.OTPTokenMiddleware("verify_account"), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			controllers.VerifyEmail(&interfaces.ApplicationContext[any]{
				Ctx:  ctx,
				Keys: appContextAny.Keys,
			})
		})

		authRouter.PATCH("/phone/verify", middlewares.OTPTokenMiddleware("verify_phone"), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.IsAuthOne
			err := ctx.ShouldBindJSON(&body)
			logger.Info("error parsing json body for PATCH /phone/verify", logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			controllers.VerifyPhone(&interfaces.ApplicationContext[dto.IsAuthOne]{
				Ctx:  ctx,
				Keys: appContextAny.Keys,
				Body: &body,
			})
		})

		authRouter.GET("/account/exits", func(ctx *gin.Context) {
			query := map[string]any{
				"email": ctx.Query("email"),
			}
			controllers.AccountWithEmailExists(&interfaces.ApplicationContext[any]{
				Ctx:   ctx,
				Query: query,
			})
		})

		authRouter.GET("/account/logout", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			controllers.LogOut(&interfaces.ApplicationContext[any]{
				Ctx:  ctx,
				Keys: appContextAny.Keys,
			})
		})

		authRouter.POST("/account/verify", middlewares.AuthenticationMiddleware(false, false), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.VerifyAccountData
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			controllers.VerifyAccount(&interfaces.ApplicationContext[dto.VerifyAccountData]{
				Ctx:  ctx,
				Body: &body,
				Keys: appContextAny.Keys,
			})
		})

		authRouter.POST("/account/password/reset", middlewares.OTPTokenMiddleware("update_password"), func(ctx *gin.Context) {
			var body dto.ResetPasswordDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			controllers.ResetPassword(&interfaces.ApplicationContext[dto.ResetPasswordDTO]{
				Ctx:  ctx,
				Body: &body,
				Keys: appContextAny.Keys,
			})
		})

		authRouter.POST("/account/password/update", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.UpdatePassword
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.UpdatePassword]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.UpdatePassword(&appContext)
		})

		authRouter.POST("/account/deactivate", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.ConfirmPin
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.ConfirmPin]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.DeactivateAccount(&appContext)
		})

		authRouter.POST("/account/transaction-pin/set", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SetTransactionPinDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.SetTransactionPinDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.SetTransactionPin(&appContext)
		})
	}
}
