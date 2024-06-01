package pubsub

import (
	"usepolymer.co/infrastructure/pubsub/rabbitmq"
	"usepolymer.co/infrastructure/pubsub/types"
)

var PubSub types.MessagingQueueType = (&rabbitmq.RabbitMQ{})
