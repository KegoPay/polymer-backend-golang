package dto

import "usepolymer.co/entities"

type AddBankDetailsDTO struct {
	BankDetails *entities.BankDetails `bson:"bankDetails" json:"bankDetails"`
}

type UpdateAddressDTO struct {
	State   string `bson:"state" json:"state" validate:"required,alpha_space,max=15,min=3"`
	LGA     string `bson:"lga" json:"lga" validate:"required,alpha_space,max=20"`
	Street  string `bson:"street" json:"street" validate:"required,max=300"`
	AuthOne bool   `json:"authone"`
}

type UpdatePhoneDTO struct {
	Phone    string `bson:"phone" json:"phone" validate:"required,numeric,len=11"`
	WhatsApp bool   `bson:"whatsapp" json:"whatsapp"`
}

type LinkNINDTO struct {
	NIN     string `bson:"nin" json:"nin" validate:"required,numeric,len=11"`
	AuthOne bool   `bson:"authOne" json:"authOne"`
}

type SetNextOfKin struct {
	AuthOne      bool   `json:"authOne"`
	FirstName    string `json:"firstName" validate:"required,max=30"`
	LastName     string `json:"lastName" validate:"required,max=30"`
	Relationship string `json:"relationship" validate:"required,max=30,oneof=mother father brother sister cousin aunt uncle grandparent wife husband child other"`
}
