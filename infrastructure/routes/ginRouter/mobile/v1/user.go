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

func UserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/profile", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContext, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			controllers.FetchUserProfile(appContext)
		})

		// userRouter.PATCH("/profile/update", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
		// 	appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
		// 	var body dto.UpdateUserDTO
		// 	if err := ctx.ShouldBindJSON(&body); err != nil {
		// 		apperrors.ErrorProcessingPayload(ctx)
		// 		return
		// 	}
		// 	appContext := interfaces.ApplicationContext[dto.UpdateUserDTO]{
		// 		Keys: appContextAny.Keys,
		// 		Body: &body,
		// 		Ctx: appContextAny.Ctx,
		// 	}
		// 	controllers.UpdateUserProfile(&appContext)
		// })

		userRouter.PATCH("/address/update", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.UpdateAddressDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.UpdateAddressDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.UpdateAddress(&appContext)
		})

		userRouter.POST("/nin/update", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.LinkNINDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.LinkNINDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.LinkNIN(&appContext)
		})

		userRouter.PATCH("/phone/update", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.UpdatePhoneDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.UpdatePhoneDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.UpdatePhone(&appContext)
		})

		userRouter.PATCH("/address/verify", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContext, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.IsAuthOne
			err := ctx.ShouldBindJSON(&body)
			logger.Info("error parsing json body for PATCH /address/verify", logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})

			controllers.VerifyCurrentAddress(&interfaces.ApplicationContext[dto.IsAuthOne]{
				Ctx:  ctx,
				Body: &body,
				Keys: appContext.Keys,
			})
		})

		userRouter.PATCH("/phone/current/verify", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.IsAuthOne
			err := ctx.ShouldBindJSON(&body)
			logger.Info("error parsing json body for PATCH /phone/current/verify", logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			controllers.VerifyCurrentPhone(&interfaces.ApplicationContext[dto.IsAuthOne]{
				Ctx:  ctx,
				Body: &body,
				Keys: appContextAny.Keys,
			})
		})

		userRouter.PATCH("/profile/payment-tag", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SetPaymentTagDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.SetPaymentTagDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.SetPaymentTag(&appContext)
		})

		userRouter.PATCH("/notification/toggle", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.ToggleNotificationOptionsDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.ToggleNotificationOptionsDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.ToggleNotificationOptions(&appContext)
		})

		userRouter.POST("/file/generate-url", middlewares.AuthenticationMiddleware(false, false), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.FileUploadOptions
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.FileUploadOptions]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.GenerateFileURL(&appContext)
		})

		userRouter.POST("/nextofkin/set", middlewares.AuthenticationMiddleware(false, false), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SetNextOfKin
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.SetNextOfKin]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.SetNextOfKin(&appContext)
		})
	}
}
