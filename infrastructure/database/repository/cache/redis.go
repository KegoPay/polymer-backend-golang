package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"usepolymer.co/infrastructure/logger"

	redisClient "usepolymer.co/infrastructure/database/connection/cache"
)

var (
	redisRepo RedisRepository
)

type RedisRepository struct {
	Clinet *redis.Client
}

func generateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}

func (redisRepo *RedisRepository) preRequest() {
	if redisRepo.Clinet == nil {
		redisRepo.Clinet = redisClient.Client
		logger.Info("redis repository initialisation complete")
	}
}

func (redisRepo *RedisRepository) CreateEntry(key string, payload interface{}, ttl time.Duration) bool {
	redisRepo.preRequest()
	_, err := redisRepo.Clinet.Set(key, payload, ttl).Result()
	if err != nil {
		logger.Error(errors.New("redis error occured while running CreateEntry"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "key",
			Data: key,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: payload,
		})
		return false
	}

	logger.Info("redis CreateEntry completed")
	return true
}

func (redisRepo *RedisRepository) FindOne(key string) *string {
	redisRepo.preRequest()

	result, err := redisRepo.Clinet.Get(key).Result()

	if err != nil {
		if err.Error() == "redis: nil" {
			return nil
		}
		logger.Error(errors.New("redis error occured while running FindOne"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "key",
			Data: key,
		})
		return nil
	}

	logger.Info("redis FindOne completed")
	return &result
}

func (redisRepo *RedisRepository) FindOneByteArray(key string) *[]byte {
	redisRepo.preRequest()

	result, err := redisRepo.Clinet.Get(key).Bytes()

	if err != nil {
		if err.Error() == "redis: nil" {
			return nil
		}
		logger.Error(errors.New("redis error occured while running FindOneByteArray"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "key",
			Data: key,
		})
		return nil
	}

	logger.Info("redis FindOneByteArray completed")
	return &result
}

func (redisRepo *RedisRepository) DeleteOne(key string) bool {
	redisRepo.preRequest()

	result, err := redisRepo.Clinet.Del(key).Result()

	if err != nil {
		logger.Error(errors.New("redis error occured while running DeleteOne"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "key",
			Data: key,
		})
		return false
	}
	if int(result) != 1 {
		return false
	}

	logger.Info("redis DeleteOne completed")
	return true
}

func (redisRepo *RedisRepository) CreateInSet(key string, score float64, member interface{}) bool {
	redisRepo.preRequest()
	added := redisRepo.Clinet.ZAdd(key, redis.Z{
		Score: score, Member: member,
	})

	if err := added.Err(); err != nil {
		logger.Error(errors.New("redis error occured while running CreateInSet"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "key",
			Data: key,
		}, logger.LoggerOptions{
			Key:  "socre",
			Data: score,
		}, logger.LoggerOptions{
			Key:  "member",
			Data: member,
		})
		return false
	}

	logger.Info("redis CreateInSet completed")
	return added != nil
}

func (redisRepo *RedisRepository) FindSet(key string) *[]string {
	redisRepo.preRequest()

	result := redisRepo.Clinet.ZRange(key, 0, -1)
	if err := result.Err(); err != nil {
		logger.Error(errors.New("redis error occured while running FindSet"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "key",
			Data: key,
		})
		return nil
	}
	if result == nil {
		return nil
	}

	logger.Info("redis FindSet completed")
	val := result.Val()
	return &val
}
