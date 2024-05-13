package gocraft

import (
	"context"
	"errors"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"kego.com/infrastructure/logger"
)


type GoCraftScheduler struct {
	Emitter 	 *work.Enqueuer
	CacheAddress  string
	CachePassword string
}

func (es *GoCraftScheduler) StartScheduler() {
	redisPool := &redis.Pool{
		MaxActive: 20,
		MaxIdle:   20,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", es.CacheAddress,
				redis.DialPassword(es.CachePassword),
				redis.DialUseTLS(true))
		},
	}
    es.Emitter = work.NewEnqueuer("polymer_jobs", redisPool)

	pool := work.NewWorkerPool(context.Background(), 56, "polymer_jobs", redisPool)
	pool.Job(string("send_email"),  SendEmail)
	pool.Job(string("lock_account"),  LockAccount)
	pool.Job(string("unlock_account"),  UnlockAccount)
	pool.Job(string("generate_account_statement"),  RequestAccountStatement)
   	pool.Start()
}


func (es *GoCraftScheduler) Emit(channel string, payload map[string]any) error {
	_, err := es.Emitter.Enqueue(channel, payload)
	if err != nil {
		logger.Error(errors.New("error scheduling job"), logger.LoggerOptions{
			Key: "channel",
			Data: channel,
		}, logger.LoggerOptions{
			Key: "payload",
			Data: payload,
		})
		return err
	}
	return nil
}