package controllers

import (
	"fmt"
	"net/http"

	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
	"kego.com/application/controllers/dto"
	countriessupported "kego.com/application/countriesSupported"
	"kego.com/application/interfaces"
	international_payment_processor "kego.com/infrastructure/payment_processor/chimoney"
	server_response "kego.com/infrastructure/serverResponse"
)

func FilterCountries(ctx *interfaces.ApplicationContext[dto.CountryFilter]){
	countries := countriessupported.FilterCountries(ctx.Body.Service)
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "countries fetched", countries, nil)
}

func FetchBanks(ctx *interfaces.ApplicationContext[any]){
	banks := bankssupported.KYCSupportedBanks
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "countries fetched", banks, nil)
}

func FetchExchangeRates(ctx *interfaces.ApplicationContext[any]){
	rates, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetExchangeRates()
	if err != nil {
		apperrors.ExternalDependencyError(ctx.Ctx, "chimoney", fmt.Sprintf("%d", statusCode), err)
		return
	}
	if statusCode != 200 {
		apperrors.UnknownError(ctx.Ctx)
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "rates fetched", rates.Data.FormatAllRates(), nil)
}
