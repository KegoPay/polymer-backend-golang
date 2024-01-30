package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
	"kego.com/application/controllers/dto"
	countriessupported "kego.com/application/countriesSupported"
	"kego.com/application/interfaces"
	"kego.com/application/services"
	"kego.com/entities"
	currencyformatter "kego.com/infrastructure/currency_formatter"
	international_payment_processor "kego.com/infrastructure/payment_processor/chimoney"
	server_response "kego.com/infrastructure/serverResponse"
)

func FilterCountries(ctx *interfaces.ApplicationContext[dto.CountryFilter]){
	countries := countriessupported.FilterCountries(ctx.Body.Service)
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "countries fetched", countries, nil)
}

func FetchLocalBanks(ctx *interfaces.ApplicationContext[any]){
	banks := bankssupported.SupportedLocalBanks
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "banks fetched", banks, nil)
}

func FetchInternationalBanks(ctx *interfaces.ApplicationContext[dto.CountryCode]){
	var banks *[]entities.Bank

	if ctx.Body.Code == "NG" {
		banks = &bankssupported.SupportedLocalBanks
	} else {
		banks = services.FetchInternationalBanks(ctx.Ctx, ctx.Body.Code)
	}
	
	if banks == nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "banks fetched", banks, nil)
}

func FetchExchangeRates(ctx *interfaces.ApplicationContext[any]){
	amountAsUInt, err := strconv.ParseUint(ctx.Query["amount"].(string), 10, 64)
	if err != nil {
		apperrors.ClientError(ctx, fmt.Sprintf("The amount %d is not a valid amount. Put in a valid amount", amountAsUInt), nil)
		return
	}
	rates, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetExchangeRates(ctx.Query["currency"], &amountAsUInt)
	if err != nil {
		apperrors.ExternalDependencyError(ctx.Ctx, "chimoney", fmt.Sprintf("%d", statusCode), err)
		return
	}
	if statusCode != 200 {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("chimoney failed to return 200 response code"))
		return
	}
	var formattedRates = map[string]string{}
	for i, r := range *rates {
		formattedRates[i] = currencyformatter.HumanReadableFloat32Currency(r)
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "rates fetched", formattedRates, nil)
}
