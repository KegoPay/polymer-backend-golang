package entities

import "fmt"

type PhoneNumber struct	{
	ISOCode 	 string `bson:"isoCode" json:"isoCode"` // Two-letter country code (ISO 3166-1 alpha-2)
	LocalNumber  string `bson:"localNumber" json:"localNumber"`
	Prefix		 string `bson:"prefix" json:"prefix"`
	IsVerified	 bool   `bson:"isVerified" json:"isVerified"`
	WhatsApp	 bool   `bson:"whatsapp" json:"whatsapp"`
	Modified	 bool   `bson:"modified" json:"modified"`
}

func (pn *PhoneNumber) ParsePhoneNumber() string {
	return fmt.Sprintf("+%s%s", pn.Prefix, pn.LocalNumber)
}