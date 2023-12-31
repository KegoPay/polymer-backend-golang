package dto

import "mime/multipart"

type CreateAccountDTO struct {
	Email      		  string                 `json:"email"`
	Password      		  string                 `json:"password"`
	UserAgent 		  string     			 `json:"deviceType"`
	DeviceID  		  string                 `json:"deviceID"`
	TransactionPin    string           		 `json:"transactionPin"`
	AppVersion        string       			 `json:"appVersion"`
	BVN    			  string           		 `json:"bvn"`
}

type LoginDTO struct {
	Email      *string                `json:"email,omitempty"`
	Phone      *string  			  `json:"phone,omitempty"`
	Password   string                 `json:"password"`
	DeviceID   string                 `json:"deviceID"`
}

type VerifyEmailData struct {
	Otp     string `json:"otp"`
	Email	string `json:"email"`
}

type VerifyAccountData struct {
	ProfileImage     *multipart.FileHeader
	Email 			 string
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

type ConfirmPin struct {
	Pin    string           		 `json:"pin"`
}
