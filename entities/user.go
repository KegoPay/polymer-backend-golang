package entities

import (
	"time"

	"kego.com/application/utils"
)


type User struct {
	FirstName         				string       `bson:"firstName" json:"firstName"`
	LastName          				string       `bson:"lastName" json:"lastName"`
	MiddleName          			*string      `bson:"middleName" json:"middleName"`
	Email            				string       `bson:"email" json:"email" validate:"required,email"`
	Phone             				PhoneNumber  `bson:"phone" json:"phone"`
	Password          				string       `bson:"password" json:"-" validate:"password"`
	TransactionPin    				string       `bson:"transactionPin" json:"-" validate:"password"`
	UserAgent        				string    	 `bson:"userAgent" json:"userAgent" validate:"user_agent,required"`
	DeviceID          				string       `bson:"deviceID" json:"deviceID"`
	AppVersion          			string       `bson:"appVersion" json:"appVersion"`
	WalletID  						string    	 `bson:"walletID" json:"walletID"`
	KYCCompleted   					bool         `bson:"kycCompleted" json:"kycCompleted"`
	EmailVerified     				bool         `bson:"emailVerified" json:"emailVerified"`
	AccountRestricted 				bool         `bson:"accountRestricted" json:"accountRestricted"`
	Deactivated 					bool         `bson:"deactivated" json:"deactivated"`
	BVN		  		  				string 	  	 `bson:"bvn" json:"bvn" validate:"required"`
	Gender		  		  			string 	  	 `bson:"gender" json:"gender"`
	DOB		  		  				string 	  	 `bson:"dob" json:"dob"`
	WatchListed		  		  		string 	  	 `bson:"watchListed" json:"watchListed"`
	Nationality		  		  		string 	  	 `bson:"nationality" json:"nationality"`
	ProfileImage		  		  	string 	  	 `bson:"profileImage" json:"profileImage"`

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
