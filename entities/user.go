package entities

import (
	"time"

	"kego.com/application/utils"
)

type User struct {
	Email       *string      `bson:"email" json:"email,omitempty" validate:"exclusive_email_phone,omitempty,email"`
	Phone       *PhoneNumber `bson:"phone" json:"phone,omitempty" validate:"exclusive_email_phone,omitempty"`
	Password    string       `bson:"password" json:"-" validate:"password"`

	ID          string       `bson:"_id" json:"id"`
	CreatedAt   time.Time    `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time    `bson:"updatedAt" json:"updatedAt"`
}

func (user User) ParseModel() any {
	if user.ID == "" {
		user.CreatedAt = time.Now()
		user.ID = utils.GenerateUUIDString()
	}
	user.UpdatedAt = time.Now()
	return &user
}
