package infrastructure

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	apperrors "kego.com/application/appErrors"
	"kego.com/infrastructure/logger"
	middlewares "kego.com/infrastructure/middleware"
	authroutev1 "kego.com/infrastructure/routes/ginRouter/v1"
	server_response "kego.com/infrastructure/serverResponse"
	startup "kego.com/infrastructure/startUp"
)

type ginServer struct {}

func (s *ginServer)Start(){
	err := godotenv.Load()

	if err != nil {
		logger.Warning("could not find .env file")
	}

	startup.StartServices()
	defer startup.CleanUpServices()

	server := gin.Default()

	server.Use(middlewares.UserAgentMiddleware())

	v1 := server.Group("/api",)

	{
		routerV1 := v1.Group("/v1")
		{
			authroutev1.AuthRouter(routerV1)
			authroutev1.InfoRouter(routerV1)
			authroutev1.UserRouter(routerV1)
			authroutev1.BusinessRouter(routerV1)
			authroutev1.WalletRouter(routerV1)
		}
	}

	server.GET("/ping", func(ctx *gin.Context) {
		server_response.Responder.Respond(ctx, http.StatusOK, "pong!", nil, nil)
	})

	server.NoRoute(func(ctx *gin.Context) {
		apperrors.NotFoundError(ctx, fmt.Sprintf("%s %s does not exist", ctx.Request.Method, ctx.Request.URL))
	})

	gin_mode := os.Getenv("GIN_MODE")
	port := os.Getenv("PORT")
	if gin_mode == "debug" || gin_mode == "release"{
		logger.Info(fmt.Sprintf("Server starting on PORT %s", port))
		server.Run(port)
	} else {
		panic("invalid gin mode used")
	}
}