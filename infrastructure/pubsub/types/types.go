package types

type MessagingQueueType interface {
	Connect()
	PublishMessage(exchange string, name string, payload interface{}) error
	CreateQueue(name string) error
}