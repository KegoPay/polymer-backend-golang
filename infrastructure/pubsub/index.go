package pubsub

import (
	"kego.com/infrastructure/pubsub/rabbitmq"
	"kego.com/infrastructure/pubsub/types"
)

var PubSub types.MessagingQueueType = (&rabbitmq.RabbitMQ{})