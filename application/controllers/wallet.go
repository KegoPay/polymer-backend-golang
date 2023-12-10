package controllers

import (
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/services"
	"kego.com/application/utils"
	"kego.com/entities"
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