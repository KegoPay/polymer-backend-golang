package dto

import "kego.com/entities"

type AddBankDetailsDTO struct {
	BankDetails		  *entities.BankDetails  	`bson:"bankDetails" json:"bankDetails"`
}

type UpdateAddressDTO struct {
	State 	  string       `bson:"state" json:"state" validate:"required,alpha_space,max=15,min=3"`
	LGA    	  string       `bson:"lga" json:"lga" validate:"required,alpha_space,max=20"`
	Street	  string       `bson:"street" json:"street" validate:"required,max=300"`
}


type UpdatePhoneDTO struct {
	Phone 	  string     `bson:"phone" json:"phone" validate:"required,numeric,len=11"`
	WhatsApp  bool       `bson:"whatsapp" json:"whatsapp"`
}
