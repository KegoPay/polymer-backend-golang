package routev1

import (
	"github.com/gin-gonic/gin"
	"usepolymer.co/application/controllers"
	"usepolymer.co/application/interfaces"
	middlewares "usepolymer.co/infrastructure/middleware"
)

func TransactionRouter(router *gin.RouterGroup) {
	transactionRouter := router.Group("/transaction")
	{
		transactionRouter.GET("/:businessID/latest", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContext, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			appContext.Param = map[string]any{
				"businessID": ctx.Param("businessID"),
			}
			controllers.FetchPastBusinessTransactions(appContext)
		})

		transactionRouter.GET("/latest", middlewares.AuthenticationMiddleware(false, true), func(ctx *gin.Context) {
			appContext, _ := ctx.MustGet("AppContext").(*interfaces.ApplicationContext[any])
			controllers.FetchPastPersonalTransactions(appContext)
		})
	}
}
