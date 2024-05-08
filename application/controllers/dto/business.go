package dto

import cac_service "kego.com/infrastructure/cac"

type BusinessDTO struct {
	Name  string `json:"name"`
	Email string `json:"-"`
}

type UpdateBusinessDTO struct {
	Name string `json:"name" validate:"required"`
	ID 	 string `json:"id" validate:"required"`
}

type SearchCACByName struct {
	Name string `json:"name" validate:"required"`
}

type SetCACInfo struct {
	Info cac_service.CACBusiness `json:"info" validate:"required"`
}