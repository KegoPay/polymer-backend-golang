package controllers

import (
	"context"
	"fmt"
	"net/http"

	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/application/services"
	"kego.com/application/utils"
	"kego.com/entities"
	international_payment_processor "kego.com/infrastructure/payment_processor/chimoney"
	server_response "kego.com/infrastructure/serverResponse"
)

func InitiateBusinessInternationalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]){
	businessID := ctx.GetStringParameter("businessID") 
	rates, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetExchangeRates(ctx.Body.DestinationCountryCode, ctx.Body.Amount)
	if err != nil {
		apperrors.ExternalDependencyError(ctx.Ctx, "chimoney", fmt.Sprintf("%d", statusCode), err)
		return
	}
	if statusCode != 200 {
		apperrors.UnknownError(ctx.Ctx)
		return
	}
	ctx.Body.Amount = utils.Float32ToUint64Currency((*rates)["convertedValue"])
	wallet , err := services.InitiatePreAuth(ctx.Ctx, businessID, ctx.GetStringContextData("UserID"), ctx.Body.Amount, ctx.Body.Pin)
	if err != nil {
		return
	}
	err = services.LockFunds(ctx.Ctx, wallet, ctx.Body.Amount, entities.ChimoneyDebitInternational)
	if err != nil {
		return
	}
	response := services.InitiateInternationalPayment(ctx.Ctx, &international_payment_processor.InternationalPaymentRequestPayload{
		DestinationCountry: "Nigeria",
		AccountNumber: ctx.Body.AccountNumber,
		BankCode: ctx.Body.BankCode,
		ValueInUSD: utils.Float32ToUint64Currency((*rates)["convertToUSD"]),
	})
	if response == nil {
		return
	}
	transaction := entities.Transaction{
		TransactionReference: response.Chimoneys[0].ChiRef,
		MetaData: response.Chimoneys[0],
		AmountInUSD: utils.Float32ToUint64Currency(response.Chimoneys[0].ValueInUSD),
		Amount: ctx.Body.Amount,
		Currency: wallet.Currency,
		WalletID: wallet.ID,
		UserID: wallet.UserID,
		BusinessID: wallet.BusinessID,
		Description: func () string {
			if	ctx.Body.Description == nil {
				des := fmt.Sprintf("International transfer from %s %s to %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName"), ctx.Body.FullName)
				return des
			}
			return *ctx.Body.Description
		}(),
		Location: entities.Location{
			IPAddress: ctx.Body.IPAddress,
		},
		Intent: entities.ChimoneyDebitInternational,
		DeviceInfo: entities.DeviceInfo{
			IPAddress: ctx.Body.IPAddress,
			DeviceID: ctx.GetStringContextData("DeviceID"),
			UserAgent: ctx.GetStringContextData("UserAgent"),
		},
		Sender: entities.TransactionSender{
			BusinessName: *wallet.BusinessName,
			FirstName: ctx.GetStringContextData("FirstName"),
			LastName: ctx.GetStringContextData("LastName"),
			Email: ctx.GetStringContextData("Email"),
		},
		Recepient: entities.TransactionRecepient{
			Name: ctx.Body.FullName,
			BankCode: ctx.Body.BankCode,
			AccountNumber: ctx.Body.AccountNumber,
			BranchCode: ctx.Body.BranchCode,
			Country: ctx.Body.DestinationCountryCode,
		},
	}
	trxRepository := repository.TransactionRepo()
	trx, err := trxRepository.CreateOne(context.TODO(), transaction)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "Your payment is on its way! 🚀", trx, nil)
}

func VerifyLocalAccountName(ctx *interfaces.ApplicationContext[dto.NameVerificationDTO]){
	bankCode := ""
	for _, bank := range bankssupported.KYCSupportedBanks {
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