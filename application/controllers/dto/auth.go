package dto

import "kego.com/entities"


type CreateAccountDTO struct {
	Email       *string  				  `json:"email,omitempty"`
	Phone       *entities.PhoneNumber     `json:"phone,omitempty"`
	Password    string     			 	  `json:"password"`
}
