package events

import "kego.com/infrastructure/pubsub"

type Event string


func SendEmail(toEmail string, subject string, template string, options interface{}){
	pubsub.PubSub.PublishMessage("send_email", map[string]any{
		"subject": subject,
		"template": template,
		"options": options,
	})
}