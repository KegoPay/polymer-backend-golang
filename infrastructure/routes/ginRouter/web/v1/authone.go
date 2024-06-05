package webroutev1

import (
	"os"

	"github.com/gin-gonic/gin"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/controllers"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/utils"
	middlewares "usepolymer.co/infrastructure/middleware"
)

func AuthOneRouter(router *gin.RouterGroup) {
	authOneRouter := router.Group("/authone")
	authOneRouter.Use(middlewares.InhouseAuthMiddleware(os.Getenv("AUTHONE_JWT_SIGNING_KEY"), os.Getenv("AUTHONE_JWT_ISSUER")))
	{
		authOneRouter.GET("/user/:email", func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			appContext := interfaces.ApplicationContext[any]{
				Ctx: appContextAny.Ctx,
				Param: map[string]any{
					"email": ctx.Param("email"),
				},
			}
			controllers.AuthOneFetchUserDetails(&appContext)
		})

		authOneRouter.POST("/email/send", func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.AuthOneSendEmail
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.AuthOneSendEmail]{
				Ctx: appContextAny.Ctx,
				Param: map[string]any{
					"email": ctx.Param("email"),
				},
				Body: &body,
			}
			controllers.AuthOneSendEmail(&appContext)
		})

		authOneRouter.POST("/email/status", func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body string
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[string]{
				Ctx:  appContextAny.Ctx,
				Body: &body,
			}
			controllers.AuthOneVerifyEmailStatus(&appContext)
		})

		authOneRouter.POST("/account/create", func(ctx *gin.Context) {
			appContextAny, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			var body dto.AuthOneCreateUserDTO
			if err := ctx.ShouldBindJSON(&body); err != nil {
				apperrors.ErrorProcessingPayload(ctx, utils.GetStringPointer(ctx.GetHeader("Polymer-Device-Id")))
				return
			}
			appContext := interfaces.ApplicationContext[dto.AuthOneCreateUserDTO]{
				Ctx:  appContextAny.Ctx,
				Body: &body,
			}
			controllers.AuthOneCreateUser(&appContext)
		})
	}
}
