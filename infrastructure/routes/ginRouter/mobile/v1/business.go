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


func BusinessRouter(router *gin.RouterGroup) {
	businessRouter := router.Group("/business")
	{
		businessRouter.POST("/create", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.BusinessDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx,  utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.BusinessDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.CreateBusiness(&appContext)
		})

		businessRouter.PATCH("/:businessID/update", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.UpdateBusinessDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx,  utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.UpdateBusinessDTO]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.UpdateBusiness(&appContext)
		})

		businessRouter.GET("/fetch", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			appContext := interfaces.ApplicationContext[any]{
				Keys: appContextAny.Keys,
				Ctx: appContextAny.Ctx,
			}
			controllers.FetchBusinesses(&appContext)
		})

		businessRouter.DELETE("/:businessID/delete", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			appContext := interfaces.ApplicationContext[any]{
				Keys: appContextAny.Keys,
				Ctx: appContextAny.Ctx,
			}
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.DeleteBusiness(&appContext)
		})

		businessRouter.POST("/cac/search/name", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SearchCACByName
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx,  utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.SearchCACByName]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			controllers.SearchCACByName(&appContext)
		})

		businessRouter.PATCH("/:businessID/cac/set", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.SetCACInfo
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx,  utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.SetCACInfo]{
				Keys: appContextAny.Keys,
				Body: &body,
				Ctx: appContextAny.Ctx,
			}
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.SetCACInfo(&appContext)
		})
	}
}
