package dto

import (
	"crypto/ecdh"

	"usepolymer.co/infrastructure/file_upload/types"
)

type CreateAccountDTO struct {
	Email                 string `json:"email"`
	Password              string `json:"password"`
	UserAgent             string `json:"deviceType"`
	DeviceID              string `json:"deviceID"`
	PushNotificationToken string `json:"pushNotificationToken"`
	AppVersion            string `json:"appVersion"`
	IPAddress             string
}

type KeyExchangeDTO struct {
	ClientPublicKey *ecdh.PublicKey `json:"clientPubKey"`
	DeviceID        string
}

type EncryptForStagingDTO struct {
	EncKey   string `json:"enc_key"`
	Payload  any    `json:"payload"`
	DeviceID string
}

type DecryptForStagingDTO struct {
	Payload  string `json:"payload"`
	EncKey   string `json:"enc_key"`
	DeviceID string
}

type ClientKeyExchangeMockDTO struct {
	ServerPublicKey *ecdh.PublicKey `json:"serverPubKey"`
	ClientPublicKey *ecdh.PublicKey `json:"clientPubKey"`
	DeviceID        string
}

type VerifyOTPDTO struct {
	OTP   string  `json:"otp"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

type LoginDTO struct {
	Email                 *string `json:"email,omitempty"`
	Phone                 *string `json:"phone,omitempty"`
	Password              string  `json:"password"`
	DeviceID              string  `json:"deviceID"`
	PushNotificationToken string  `json:"pushNotificationToken"`
}

type VerifyEmailData struct {
	Otp   string `json:"otp"`
	Email string `json:"email"`
}

type SetBVNDTOO struct {
}

type VerifyAccountData struct {
	ProfileImage string `json:"profileImage" validate:"required,url"`
	Path         string `json:"path" validate:"required,oneof=nin bvn"`
}

type SetIDForBiometricVerificationDTO struct {
	ID   string `json:"id"`
	Path string `json:"path" validate:"required,oneof=nin bvn"`
}

type SetTransactionPinDTO struct {
	TransactionPin string `json:"transactionPin" validate:"required,min=4,max=4"`
	UserImage      string
	Email          string
}

type VerifyPassword struct {
	Password string `json:"password"`
}

type ResetPasswordDTO struct {
	NewPassword string `json:"newPassword"`
}

type UpdatePassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type ConfirmPin struct {
	Pin string `json:"pin"`
}

type FileUploadOptions struct {
	Permissions types.SignedURLPermission `json:"permissions" validate:"required"`
}

type ResendOTP struct {
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Whatsapp *bool   `json:"whatsapp"`
	Intent   string  `json:"intent" validate:"required,oneof=verify_account update_password verify_phone"`
}
