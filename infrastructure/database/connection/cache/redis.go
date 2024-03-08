package cache

import (
	"crypto/tls"
	"os"

	"github.com/go-redis/redis"
	"kego.com/infrastructure/logger"
)

var (
	Client *redis.Client
)

func connectRedis() {
	opt := &redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
		TLSConfig: &tls.Config{},
		PoolSize: 50,
	}
	Client = redis.NewClient(opt)
	logger.Info("connected to redis successfully")
}
