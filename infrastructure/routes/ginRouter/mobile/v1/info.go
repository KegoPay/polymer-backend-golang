package routev1

import (
	"github.com/gin-gonic/gin"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/controllers"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/utils"
	middlewares "usepolymer.co/infrastructure/middleware"
)

func InfoRouter(router *gin.RouterGroup) {
	infoRouter := router.Group("/info")
	{
		infoRouter.POST("/countries", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			var body dto.CountryFilter
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			controllers.FilterCountries(&interfaces.ApplicationContext[dto.CountryFilter]{
				Keys: ctx.Keys,
				Ctx:  ctx,
				Body: &body,
			})
		})

		infoRouter.GET("/banks/local", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			controllers.FetchLocalBanks(&interfaces.ApplicationContext[any]{
				Ctx: ctx,
			})
		})

		infoRouter.GET("/states", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			controllers.FetchStateData(&interfaces.ApplicationContext[any]{
				Ctx: ctx,
			})
		})

		infoRouter.POST("/banks/international", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.CountryCode
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.CountryCode]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.FetchInternationalBanks(&appContext)
		})

		infoRouter.POST("/exchange-rates", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.FXRateDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.FXRateDTO]{
				Body: &body,
				Ctx:  appContextAny.Ctx,
			}
			controllers.FetchExchangeRates(&interfaces.ApplicationContext[dto.FXRateDTO]{
				Body: &body,
				Ctx:  appContext.Ctx,
			})
		})
	}
}
