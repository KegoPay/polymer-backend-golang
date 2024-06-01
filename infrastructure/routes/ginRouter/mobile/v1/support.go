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

func SupportRouter(router *gin.RouterGroup) {
	supportRouter := router.Group("/support")
	{
		supportRouter.POST("/error/report", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			var body dto.ErrorSupportRequestDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			controllers.ErrSupportRequest(&interfaces.ApplicationContext[dto.ErrorSupportRequestDTO]{
				Keys: ctx.Keys,
				Ctx:  ctx,
				Body: &body,
			})
		})
	}
}
