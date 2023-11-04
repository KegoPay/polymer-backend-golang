package controllers

import (
	"net/http"

	"kego.com/application/controllers/dto"
	countriessupported "kego.com/application/countriesSupported"
	"kego.com/application/interfaces"
	server_response "kego.com/infrastructure/serverResponse"
)

func FilterCountries(ctx *interfaces.ApplicationContext[dto.CountryFilter]){
	countries := countriessupported.FilterCountries(ctx.Body.Service)
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "countries fetched", countries, nil)
}