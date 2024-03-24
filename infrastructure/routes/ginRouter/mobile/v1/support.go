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


func SupportRouter(router *gin.RouterGroup) {
	supportRouter := router.Group("/support")
	{
		supportRouter.POST("/error/report", middlewares.AuthenticationMiddleware(false, true),  func(ctx *gin.Context) {
			var body dto.ErrorSupportRequestDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx,  utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			controllers.ErrSupportRequest(&interfaces.ApplicationContext[dto.ErrorSupportRequestDTO]{
				Keys: ctx.Keys,
				Ctx: ctx,
				Body: &body,
			})
		})
	}
}
