package dto

import "kego.com/entities"

type CreateAccountDTO struct {
	FirstName         string       	  		 `json:"firstName"`
	LastName          string       	 		 `json:"lastName"`
	Email      		  string                `json:"email,omitempty"`
	Phone      		  entities.PhoneNumber  `json:"phone,omitempty"`
	Password   		  string                 `json:"password"`
	DeviceType 		  entities.DeviceType    `json:"deviceType"`
	DeviceID  		  string                 `json:"deviceID"`
	TransactionPin    string           		 `json:"transactionPin"`
	BVN    			  string           		 `json:"bvn"`
	BankDetails 	  entities.BankDetails	 `json:"bankDetails"`
}

type LoginDTO struct {
	Email      *string                `json:"email,omitempty"`
	Phone      *string  			  `json:"phone,omitempty"`
	Password   string                 `json:"password"`
}

type VerifyAccountData struct {
	Otp     string `json:"otp"`
	Email	string `json:"email"`
}

type VerifyPassword struct {
	Password  string `json:"password"`
}

type ResetPasswordDTO struct {
	Otp         string `json:"otp"`
	NewPassword string `json:"newPassword"`
	Email       string `json:"email"`
}

type UpdatePassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
