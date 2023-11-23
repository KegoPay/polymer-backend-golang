package dto

import "kego.com/entities"

type UpdateUserDTO struct {
	FirstName         *string       			`bson:"firstName" json:"firstName"`
	LastName          *string       			`bson:"lastName" json:"lastName"`
	Phone             *entities.PhoneNumber		`bson:"phone" json:"phone,omitempty"`
	BankDetails		  *entities.BankDetails  	`bson:"bankDetails" json:"bankDetails"`
}