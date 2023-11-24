package entities

import (
	"time"

	"kego.com/application/utils"
)

type UserMetaData struct {
	CustomerID         *int       `bson:"customerID" json:"customerID"`
	CustomerCode       *string       `bson:"customerCode" json:"customerCode"`
}

type User struct {
	FirstName         				string       `bson:"firstName" json:"firstName" validate:"required"`
	LastName          				string       `bson:"lastName" json:"lastName" validate:"required"`
	Email            				string       `bson:"email" json:"email,omitempty" validate:"required,email"`
	Phone             				PhoneNumber  `bson:"phone" json:"phone,omitempty" validate:"required"`
	Password          				string       `bson:"password" json:"-" validate:"password"`
	TransactionPin    				string       `bson:"transactionPin" json:"-" validate:"password"`
	UserAgent        				UserAgent    `bson:"userAgent" json:"userAgent" validate:"user_agent,required"`
	DeviceID          				string       `bson:"deviceID" json:"deviceID" validate:"required"`
	AppVersion          			string       `bson:"appVersion" json:"appVersion" validate:"required"`
	KYCFailedReason    				*string      `bson:"kycFailedReason" json:"kycFailedReason"`
	KYCCompleted   					bool         `bson:"kycCompleted" json:"kycCompleted"`
	EmailVerified     				bool         `bson:"emailVerified" json:"emailVerified"`
	PhoneVerified     				bool         `bson:"phoneVerified" json:"phoneVerified"`
	AccountRestricted 				bool         `bson:"accountRestricted" json:"accountRestricted"`
	Deactivated 					bool         `bson:"deactivated" json:"deactivated"`
	BankDetails		  				BankDetails  `bson:"bankDetails" json:"bankDetails" validate:"required"`
	BVN		  		  				string 	  	 `bson:"bvn" json:"bvn" validate:"required"`
	MetaData		  		  		*UserMetaData `bson:"metadata" json:"metadata"`

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
