package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/options"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/application/services"
	"kego.com/application/utils"
	currencyformatter "kego.com/infrastructure/currency_formatter"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
	pushnotification "kego.com/infrastructure/messaging/push_notifications"
	server_response "kego.com/infrastructure/serverResponse"
)

func FlutterwaveWebhook(ctx *interfaces.ApplicationContext[dto.FlutterwaveWebhookDTO]) {
	if ctx.Body.EventType == "Transfer" {
		var err error
		if ctx.Body.Transfer.Status == "SUCCESSFUL" {
			err = services.RemoveLockFunds(ctx.Ctx, ctx.Body.Transfer.Meta.WalletID, ctx.Body.Transfer.Ref)
		}else if ctx.Body.Transfer.Status == "FAILED" {
			err = services.ReverseLockFunds(ctx.Ctx, ctx.Body.Transfer.Meta.WalletID, ctx.Body.Transfer.Ref)
		}
		if err != nil {
			logger.Error(errors.New("error removing locked funds after flutterwave webhook call"), logger.LoggerOptions{
				Key: "error",
				Data: err,
			}, logger.LoggerOptions{
				Key: "payload",
				Data: ctx.Body,
			})
			return
		}
		userRepo := repository.UserRepo()
		user, err := userRepo.FindByID(ctx.Body.Transfer.Meta.UserID, options.FindOne().SetProjection(map[string]any{
			"notificationOptions": 1,
			"email": 1,
			"deviceID": 1,
			"firstName": 1,
		}))
		if err != nil {
			logger.Error(errors.New("error fetching user information after flutterwave webhook call"), logger.LoggerOptions{
				Key: "error",
				Data: err,
			}, logger.LoggerOptions{
				Key: "payload",
				Data: ctx.Body,
			})
			return
		}
		if user == nil {
			logger.Error(errors.New("user account not found during flutterwave webhook run"), logger.LoggerOptions{
				Key: "payload",
				Data: ctx.Body,
			})
			return
		}
		
		if user.NotificationOptions.PushNotification {
			pushnotification.PushNotificationService.PushOne(user.DeviceID, "Your payment was successful! 🚀",
				fmt.Sprintf("Your payment of %s%s to %s in %s has been processed successfully.", ctx.Body.Transfer.Currency, currencyformatter.HumanReadableFloat32Currency(ctx.Body.Transfer.Amount), ctx.Body.Transfer.RecepientName, utils.CurrencyCodeToCountryCode(ctx.Body.Transfer.Currency)))
		}
	
		if user.NotificationOptions.Emails {
			emails.EmailService.SendEmail(user.Email, "Your payment is on its way! 🚀", "payment_recieved", map[string]any{
				"FIRSTNAME": user.FirstName,
				"CURRENCY_CODE": utils.CurrencyCodeToCurrencySymbol("NGN"),
				"AMOUNT": currencyformatter.HumanReadableFloat32Currency(ctx.Body.Transfer.Amount),
				"RECEPIENT_NAME": ctx.Body.Transfer.RecepientName,
				"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName("Nigeria"),
			})
		}
	}else {
		logger.Warning("flutterwave webook hit without reaction", logger.LoggerOptions{
			Key: "payload",
			Data: ctx.Body,
		})
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "processed successfully", nil, nil)
}