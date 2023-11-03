package cache

import (
	"os"

	"github.com/go-redis/redis/v8"
	"kego.com/infrastructure/logger"
)

var (
	Client *redis.Client
)

func connectRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		logger.Warning("could not connect to redis", logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return
	}

	Client = redis.NewClient(opt)
	logger.Info("connected to redis successfully")
}
