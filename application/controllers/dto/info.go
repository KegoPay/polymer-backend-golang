package dto

import "kego.com/entities"

type CountryFilter struct {
	Service entities.CountryServiceType	`json:"service"`
}

type CountryCode struct {
	Code string	`json:"code"`
}
