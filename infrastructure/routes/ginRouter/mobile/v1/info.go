package routev1

import (
	"github.com/gin-gonic/gin"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	middlewares "kego.com/infrastructure/middleware"
)


func InfoRouter(router *gin.RouterGroup) {
	infoRouter := router.Group("/info")
	{
		infoRouter.POST("/countries", middlewares.AuthenticationMiddleware(false, true),  func(ctx *gin.Context) {
			var body dto.CountryFilter
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			controllers.FilterCountries(&interfaces.ApplicationContext[dto.CountryFilter]{
				Keys: ctx.Keys,
				Ctx: ctx,
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
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			appContext := interfaces.ApplicationContext[dto.CountryCode]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.FetchInternationalBanks(&appContext)
		})

		infoRouter.GET("/exchange-rates", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.FXRateDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx)
				return
			}
			appContext := interfaces.ApplicationContext[dto.FXRateDTO]{
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.FetchExchangeRates(&interfaces.ApplicationContext[dto.FXRateDTO]{
				Body: &body,
				Ctx: appContext.Ctx,
			})
		})
	}
}
