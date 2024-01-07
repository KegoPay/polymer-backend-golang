package rabbitmq

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/streadway/amqp"
	"kego.com/infrastructure/logger"
)

type RabbitMQ struct {
	QueueChannel *amqp.Channel
}

var (
	QueueChannel *amqp.Channel
)


func (rmq *RabbitMQ) Connect() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		logger.Error(errors.New("set RabbitMQ url"))
		return
	}
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		logger.Error(errors.New("could not connect to rabbitmq"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return
	}
	logger.Info("connected to RabbitMQ")

	QueueChannel, err = conn.Channel()
	if err != nil {
		logger.Error(errors.New("could not open channel on rabbitmq"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return
	}
	return
}

func (rmq *RabbitMQ) CreateQueue(name string) error {
	_, err := QueueChannel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		logger.Error(errors.New("could not open channel on rabbitmq"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	return nil
}

func (rmq *RabbitMQ) PublishMessage(exchange string, name string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error(errors.New("could not marshal payload"),  logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	QueueChannel.Publish(exchange, name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	logger.Info("published to queue", logger.LoggerOptions{
		Key: "queue_name",
		Data: name,
	})
	return nil
}
