package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
	"kego.com/application/constants"
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

func FetchStateData(ctx *interfaces.ApplicationContext[any]){
	states := constants.States
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "states fetched", states, nil)
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
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("The amount %d is not a valid amount. Put in a valid amount", amountAsUInt), nil)
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
	var countries []entities.Country
	for _, country := range countriessupported.CountriesSupported {
		for c, currency := range *rates {
			if strings.Contains(c, country.ISOCode) {
				country.Rate = currencyformatter.HumanReadableFloat32Currency(currency)
				countries = append(countries, country)
				continue
			}
		}
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "rates fetched", countries, nil)
}
