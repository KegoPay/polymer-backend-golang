package auth

import "usepolymer.co/entities"

type ClaimsData struct {
	Issuer                string
	UserID                string
	BusinessID            *string
	FirstName             string
	LastName              string
	Email                 *string
	PhoneNum              *string
	Phone                 *entities.PhoneNumber
	ExpiresAt             int64
	IssuedAt              int64
	UserAgent             string
	DeviceID              string
	PushNotificationToken string
	AppVersion            string
	OTPIntent             string
}
