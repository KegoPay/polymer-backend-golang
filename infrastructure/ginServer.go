package infrastructure

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	apperrors "kego.com/application/appErrors"
	"kego.com/infrastructure/logger"
	middlewares "kego.com/infrastructure/middleware"
	ratelimiter "kego.com/infrastructure/rateLimiter"
	routev1 "kego.com/infrastructure/routes/ginRouter/mobile/v1"
	webroutev1 "kego.com/infrastructure/routes/ginRouter/web/v1"
	server_response "kego.com/infrastructure/serverResponse"
	startup "kego.com/infrastructure/startUp"
)

type ginServer struct {}

func (s *ginServer)Start(){
	err := godotenv.Load()

	startup.StartServices()
	defer startup.CleanUpServices()

	if err != nil {
	}

	server := gin.Default()
	origins := []string{}
	if os.Getenv("GIN_MODE") == "debug" {
		origins = append(origins, "http://localhost:5173")
	}else if os.Getenv("GIN_MODE") == "release" {
		origins = append(origins, "https://usepolymer.co",  "https://www.usepolymer.co",  "www.usepolymer.co", "www.usepolymer.co", "https://www.usepolymer.co/")
	}
	corsConfig := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "web-api-key", "polymer-device-id", "User-Agent"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	server.Use(cors.New(corsConfig))
	server.Use(ratelimiter.LeakyBucket())
	server.MaxMultipartMemory =  15 << 20  // 8 MiB

	server.Use(logger.MetricMonitor.MetricMiddleware().(gin.HandlerFunc))

	v1 := server.Group("/api",)

	{
		routerV1 := v1.Group("/v1")
		routerV1.Use(middlewares.UserAgentMiddleware(true))
		{
			routev1.AuthRouter(routerV1)
			routev1.InfoRouter(routerV1)
			routev1.UserRouter(routerV1)
			routev1.BusinessRouter(routerV1)
			routev1.WalletRouter(routerV1)
			routev1.TransactionRouter(routerV1)
		}

		webRouterV1 := v1.Group("/v1/web")
		webRouterV1.Use(middlewares.UserAgentMiddleware(false))
		{
			webroutev1.EmailSubsRouter(webRouterV1)
			webroutev1.WebhookRouter(webRouterV1)
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
		server.Run(fmt.Sprintf(":%s", port))
	} else {
		panic(fmt.Sprintf("invalid gin mode used - %s", gin_mode))
	}
}