package controllers

import (
	"net/http"

	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/usecases/webhooks/flutterwave"
	"usepolymer.co/infrastructure/logger"
	server_response "usepolymer.co/infrastructure/serverResponse"
)

func FlutterwaveWebhook(ctx *interfaces.ApplicationContext[dto.FlutterwaveWebhookDTO]) {
	event := ctx.Body.EventType
	if event == "Transfer" {
		flutterwave.FlwTransferWebhook(*ctx.Body)
	} else if event == "BANK_TRANSFER_TRANSACTION" {
		flutterwave.CreditWebHook(*ctx.Body)
	} else {
		logger.Warning("flutterwave webhook hit without reaction", logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
	}

	server_response.Responder.UnEncryptedRespond(ctx.Ctx, http.StatusOK, "processed successfully", nil, nil, nil)
}

func ChimoneyWebhook(ctx *interfaces.ApplicationContext[dto.FlutterwaveWebhookDTO]) {
	event := ctx.Body.EventType
	if event == "Transfer" {
		flutterwave.FlwTransferWebhook(*ctx.Body)
	} else if event == "BANK_TRANSFER_TRANSACTION" {
		flutterwave.CreditWebHook(*ctx.Body)
	} else {
		logger.Warning("flutterwave webhook hit without reaction", logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
	}

	server_response.Responder.UnEncryptedRespond(ctx.Ctx, http.StatusOK, "processed successfully", nil, nil, nil)
}
