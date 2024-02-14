package dto

import (
	"kego.com/infrastructure/file_upload/types"
)

type CreateAccountDTO struct {
	Email      		  string                 `json:"email"`
	Password      		  string             `json:"password"`
	UserAgent 		  string     			 `json:"deviceType"`
	DeviceID  		  string                 `json:"deviceID"`
	AppVersion        string       			 `json:"appVersion"`
}

type GenerateServerPublicKey struct {
	ClientPubKey	string		`json:"clientPubKey"`
	SessionID		string		`json:"sessionID"`
}

type VerifyOTPDTO struct {
	OTP		string		`json:"otp"`
	Email	*string		`json:"email"`
	Phone	*string		`json:"phone"`
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
	Path     		 string `json:"path" validate:"required,oneof=nin bvn"`
	BVN    			 *string `json:"bvn"`
	NIN    			 *string `json:"nin"`
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

type ResendOTP struct {
	Email 	*string		`json:"email"`
	Phone 	*string		`json:"phone"`
	Intent 	string		`json:"intent" validate:"required,oneof=verify_account update_password verify_phone"`
}