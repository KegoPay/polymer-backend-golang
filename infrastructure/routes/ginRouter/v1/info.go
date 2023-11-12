package authroutev1

import (
	"github.com/gin-gonic/gin"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
)


func InfoRouter(router *gin.RouterGroup) {
	infoRouter := router.Group("/info")
	{
		infoRouter.POST("/countries", func(ctx *gin.Context) {
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

		infoRouter.GET("/banks", func(ctx *gin.Context) {
			controllers.FetchBanks(&interfaces.ApplicationContext[any]{
				Ctx: ctx,
			})
		})
	}
}
