package controllers

import (
	"fmt"
	"net/http"

	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/services"
	"kego.com/application/utils"
	"kego.com/entities"
	server_response "kego.com/infrastructure/serverResponse"
)

func SendInternationalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]){
	businessID := ctx.GetStringParameter("businessID") 
	wallet , err := services.InitiatePreAuth(ctx.Ctx, businessID, ctx.GetStringContextData("UserID"), utils.ParseAmountToSmallerDenomination(ctx.Body.Amount), ctx.Body.Pin)
	if err != nil {
		return
	}
	err = services.LockFunds(ctx.Ctx, wallet, utils.ParseAmountToSmallerDenomination(ctx.Body.Amount), entities.ChimoneyDebitInternational)
	if err != nil {
		return
	}
}

func VerifyLocalAccountName(ctx *interfaces.ApplicationContext[dto.NameVerificationDTO]){
	bankCode := ""
	for _, bank := range bankssupported.SupportedBanks {
		if bank.Name == ctx.Body.BankName {
			bankCode = bank.Code
			break
		}
	}
	if bankCode  == "" {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("%s is not a supported bank on our platform yet.", ctx.Body.BankName))
		return
	}
	name := services.NameVerification(ctx.Ctx, ctx.Body.AccountNumber, bankCode)
	if name == nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "name verification complete", map[string]string{
		"name": *name,
	}, nil)
}