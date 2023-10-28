package cache

import (
	"errors"
	"os"

	"github.com/go-redis/redis/v8"
	"kego.com/infrastructure/logger"
)

var (
	client *redis.Client
)

func connectRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		logger.Error(errors.New("could not connect to redis"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return
	}

	client = redis.NewClient(opt)
	logger.Info("connected to redis successfully")
}
