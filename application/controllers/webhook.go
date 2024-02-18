package controllers

import (
	"net/http"

	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/usecases/webhooks/flutterwave"
	"kego.com/infrastructure/logger"
	server_response "kego.com/infrastructure/serverResponse"
)

func FlutterwaveWebhook(ctx *interfaces.ApplicationContext[dto.FlutterwaveWebhookDTO]) {
	event := *ctx.Body.EventType
	if event == "Transfer" {
		flutterwave.FlwTransferWebhook(ctx.Ctx, *ctx.Body)
	} else if event == "BANK_TRANSFER_TRANSACTION" {
		flutterwave.CreditWebHook(*ctx.Body)
	} else {
		logger.Warning("flutterwave webook hit without reaction", logger.LoggerOptions{
			Key: "payload",
			Data: ctx.Body,
		})
	}

	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "processed successfully", nil, nil)
}