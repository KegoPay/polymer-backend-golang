package entities

import "fmt"

type PhoneNumber struct	{
	ISOCode 	 string `bson:"isoCode" json:"isoCode" validate:"iso3166_1_alpha2"` // Two-letter country code (ISO 3166-1 alpha-3)
	LocalNumber  string `bson:"localNumber" json:"localNumber"`
	Prefix		 string `bson:"prefix" json:"prefix"`
}

func (pn *PhoneNumber) ParsePhoneNumber() string {
	return fmt.Sprintf("+%s%s", pn.Prefix, pn.LocalNumber)
}