package dto

import "usepolymer.co/entities"

type CountryFilter struct {
	Service entities.CountryServiceType `json:"service"`
}

type CountryCode struct {
	Code string `json:"code"`
}

type FXRateDTO struct {
	Amount   *uint64 `json:"amount"`
	Currency *string `json:"currency"`
}
