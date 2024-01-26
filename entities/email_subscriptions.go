package entities

import (
	"time"

	"kego.com/application/utils"
)

type SubscriptionChannels string
var NewsLetter SubscriptionChannels = "news_letter"


type Subscriptions struct {
	Email            string       			`bson:"email" json:"email" validate:"required,email"`
	Channel          SubscriptionChannels   `bson:"channel" json:"channel" validate:"oneof=news_letter"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (subs Subscriptions) ParseModel() any {
	if subs.ID == "" {
		subs.CreatedAt = time.Now()
		subs.ID = utils.GenerateUUIDString()
	}
	subs.UpdatedAt = time.Now()
	return &subs
}
