package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"usepolymer.co/infrastructure/logger"
)

type RabbitMQ struct {
	QueueChannel *amqp.Channel
}

func (rmq *RabbitMQ) Connect() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		logger.Error(errors.New("set RabbitMQ url"))
		return
	}
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		logger.Error(errors.New("could not connect to rabbitmq"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return
	}
	logger.Info("connected to RabbitMQ")
	defer conn.Close()

	rmq.QueueChannel, err = conn.Channel()
	defer rmq.QueueChannel.Close()
	if err != nil {
		logger.Error(errors.New("could not open channel on rabbitmq"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return
	}
	rmq.CreateQueue(fmt.Sprintf("%s_queue", os.Getenv("ENV")))

	// var wg sync.WaitGroup
	// wg.Add(1)

	returned := make(chan amqp.Return)
	rmq.QueueChannel.NotifyReturn(returned)
	go func() {
		// defer wg.Done()
		for returnedMsg := range returned {
			logger.Error(errors.New("could not publish message"), logger.LoggerOptions{
				Key:  "data",
				Data: returnedMsg,
			})
		}
	}()

	// wg.Wait()
	return
}

func (rmq *RabbitMQ) CreateQueue(name string) error {
	_, err := rmq.QueueChannel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		logger.Error(errors.New("could not open channel on rabbitmq"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return err
	}
	logger.Info("opened channel on rabbitmq", logger.LoggerOptions{
		Key:  "name",
		Data: name,
	})
	return nil
}

func (rmq *RabbitMQ) PublishMessage(key string, payload any) {
	body, err := json.Marshal(map[string]any{
		"action":  key,
		"payload": payload,
	})
	if err != nil {
		logger.Error(errors.New("could not marshal payload"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return
	}
	rmq.QueueChannel.PublishWithContext(context.Background(), "", key, true, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	logger.Info("published to queue", logger.LoggerOptions{
		Key:  "queue_name",
		Data: key,
	})
	return
}
