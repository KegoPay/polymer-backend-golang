package entities

import (
	"time"

	"kego.com/application/utils"
)

type User struct {
	FirstName         string       `bson:"firstName" json:"firstName" validate:"required"`
	LastName          string       `bson:"lastName" json:"lastName" validate:"required"`
	Email             *string      `bson:"email" json:"email,omitempty" validate:"exclusive_email_phone,omitempty,email"`
	Phone             *PhoneNumber `bson:"phone" json:"phone,omitempty" validate:"exclusive_email_phone,omitempty"`
	// BVN               string       `bson:"bvn" json:"bvn"`
	Password          string       `bson:"password" json:"-" validate:"password"`
	TransactionPin    string       `bson:"transactionPin" json:"transactionPin" validate:"password"`
	DeviceType        DeviceType   `bson:"deviceType" json:"deviceType" validate:"required,oneof=android ios"`
	DeviceID          string       `bson:"deviceID" json:"deviceID" validate:"required"`
	// AccountVerified   bool         `bson:"accountVerified" json:"accountVerified"`
	// KYCCompleted   bool         `bson:"kycCompleted" json:"accoukycCompletedntVerified"`
	EmailVerified     *bool        `bson:"emailVerified" json:"emailVerified"`
	PhoneVerified     *bool        `bson:"phoneVerified" json:"phoneVerified"`
	AccountRestricted bool         `bson:"accountRestricted" json:"accountRestricted"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (user User) ParseModel() any {
	if user.ID == "" {
		user.CreatedAt = time.Now()
		user.ID = utils.GenerateUUIDString()
	}
	user.UpdatedAt = time.Now()
	return &user
}
