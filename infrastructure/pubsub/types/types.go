package types

type MessagingQueueType interface {
	Connect()
	PublishMessage(name string, payload any)
	CreateQueue(name string) error
}