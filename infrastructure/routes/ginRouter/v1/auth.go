package authroutev1

import (
	"github.com/gin-gonic/gin"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/entities"
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
	}
}
