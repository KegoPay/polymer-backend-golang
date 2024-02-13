package dto

import "kego.com/entities"

type UpdateUserDTO struct {
	FirstName         *string       			`bson:"firstName" json:"firstName"`
	LastName          *string       			`bson:"lastName" json:"lastName"`
	Phone             *entities.PhoneNumber		`bson:"phone" json:"phone,omitempty"`
	BankDetails		  *entities.BankDetails  	`bson:"bankDetails" json:"bankDetails"`
}

type UpdateAddressDTO struct {
	State 	  string       `bson:"state" json:"state" validate:"required,alpha_space,max=15,min=3"`
	LGA    	  string       `bson:"lga" json:"lga" validate:"required,alpha_space,max=20"`
	Street	  string       `bson:"street" json:"street" validate:"required,max=300"`
}
