package dto

import (
	"kego.com/infrastructure/file_upload/types"
)

type CreateAccountDTO struct {
	Email      		  string                 `json:"email"`
	Password      		  string             `json:"password"`
	UserAgent 		  string     			 `json:"deviceType"`
	DeviceID  		  string                 `json:"deviceID"`
	// TransactionPin    string           		 `json:"transactionPin"`
	AppVersion        string       			 `json:"appVersion"`
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

type SetBVNDTOO struct {
}

type VerifyAccountData struct {
	ProfileImage     string `json:"profileImage" validate:"required,url"`
	BVN    			 string `json:"bvn" validate:"required,min=11,max=11"`
}

type SetTransactionPinDTO struct {
	TransactionPin   string `json:"transactionPin" validate:"required,min=4,max=4"`
	UserImage     string
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

type FileUploadOptions struct {
	Type 		string					  `json:"type" validate:"required,oneof=biometric profile_image"`
	Permissions types.SignedURLPermission `json:"permissions" validate:"required"`
}