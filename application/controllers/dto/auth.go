package dto

import "kego.com/entities"

type CreateAccountDTO struct {
	FirstName         string       	  `json:"firstName"`
	LastName          string       	  `json:"lastName"`
	Email      *string                `json:"email,omitempty"`
	Phone      *entities.PhoneNumber  `json:"phone,omitempty"`
	Password   string                 `json:"password"`
	DeviceType entities.DeviceType    `json:"deviceType"`
	DeviceID   string                 `json:"deviceID"`
	TransactionPin   string           `json:"transactionPin"`
}
