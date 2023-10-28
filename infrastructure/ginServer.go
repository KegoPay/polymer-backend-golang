package infrastructure

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"kego.com/infrastructure/logger"
	server_response "kego.com/infrastructure/serverResponse"
	startup "kego.com/infrastructure/startUp"
)

type ginServer struct {}

func (s *ginServer)Start(){
	err := godotenv.Load()

	if err != nil {
		logger.Error(errors.New("could not find .env file"))
	}

	startup.StartServices()
	defer startup.CleanUpServices()

	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		server_response.Responder.Respond(ctx, http.StatusOK, "pong!", nil, nil)
	})

	server.NoRoute(func(ctx *gin.Context) {
		server_response.Responder.Respond(ctx, http.StatusNotFound, fmt.Sprintf("%s %s does not exist", ctx.Request.Method, ctx.Request.URL), nil, nil)
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