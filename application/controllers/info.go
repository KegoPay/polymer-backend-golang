package controllers

import (
	"fmt"
	"net/http"
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
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "countries fetched", countries, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func FetchLocalBanks(ctx *interfaces.ApplicationContext[any]){
	banks := bankssupported.SupportedLocalBanks
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "banks fetched", banks, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func FetchStateData(ctx *interfaces.ApplicationContext[any]){
	states := constants.States
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "states fetched", states, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func FetchInternationalBanks(ctx *interfaces.ApplicationContext[dto.CountryCode]){
	var banks *[]entities.Bank

	if ctx.Body.Code == "NG" {
		banks = &bankssupported.SupportedLocalBanks
	} else {
		banks = services.FetchInternationalBanks(ctx.Ctx, ctx.Body.Code, ctx.GetHeader("Polymer-Device-Id"))
	}
	
	if banks == nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "banks fetched", banks, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func FetchExchangeRates(ctx *interfaces.ApplicationContext[dto.FXRateDTO]){
	rates, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetExchangeRates(ctx.Body.Amount)
	if err != nil {
			apperrors.ExternalDependencyError(ctx.Ctx, "chimoney", fmt.Sprintf("%d", statusCode), err, ctx.GetHeader("Polymer-Device-Id"))
			return
	}
	if statusCode != 200 {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("chimoney failed to return 200 response code"), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	
	if ctx.Body.Currency != nil {
		var country entities.Country
		for _, c := range countriessupported.CountriesSupported {
			if c.ISOCode == *ctx.Body.Currency {
				for key, rate := range *rates {
					if strings.Contains(key, c.ISOCode) {
						c.Rate = currencyformatter.HumanReadableFloat32Currency(rate.NGNRate)
						country = c
						break
					}
				}
				break
			}
		}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "rate fetched", country, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
	return
	}
	var countries []entities.Country
	for _, country := range countriessupported.CountriesSupported {
		for c, currency := range *rates {
			if strings.Contains(c, country.ISOCode) {
				country.Rate = currencyformatter.HumanReadableFloat32Currency(currency.NGNRate)
				fmt.Println(country.Name)
				fmt.Println(currency.NGNRate)
				countries = append(countries, country)
				continue
			}
		}
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "rates fetched", countries, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}
