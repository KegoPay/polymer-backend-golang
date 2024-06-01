package flutterwave

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/repository"
	"usepolymer.co/application/services"
	"usepolymer.co/application/utils"
	"usepolymer.co/infrastructure/background"
	currencyformatter "usepolymer.co/infrastructure/currency_formatter"
	"usepolymer.co/infrastructure/logger"
	pushnotification "usepolymer.co/infrastructure/messaging/push_notifications"
)

func FlwTransferWebhook(body dto.FlutterwaveWebhookDTO) error {
	var err error
	if body.Transfer.Status == "SUCCESSFUL" {
		err = services.RemoveLockFunds(body.Transfer.Meta.WalletID, body.Transfer.Ref)
	} else if body.Transfer.Status == "FAILED" {
		err = services.ReverseLockFunds(body.Transfer.Meta.WalletID, body.Transfer.Ref)
	}
	if err != nil {
		logger.Error(errors.New("error removing locked funds after flutterwave webhook call"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: body,
		})
		return err
	}
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(body.Transfer.Meta.UserID, options.FindOne().SetProjection(map[string]any{
		"notificationOptions": 1,
		"email":               1,
		"deviceID":            1,
		"firstName":           1,
	}))
	if err != nil {
		logger.Error(errors.New("error fetching user information after flutterwave webhook call"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: body,
		})
		return err
	}
	if user == nil {
		err = errors.New("user account not found during flutterwave webhook run")
		logger.Error(err, logger.LoggerOptions{
			Key:  "payload",
			Data: body,
		})
		return err
	}

	if user.NotificationOptions.PushNotification {
		pushnotification.PushNotificationService.PushOne(user.PushNotificationToken, "Your payment was successful! 🚀",
			fmt.Sprintf("Your payment of %s%s to %s in %s has been processed successfully.", body.Transfer.Currency, currencyformatter.HumanReadableFloat32Currency(body.Transfer.Amount), body.Transfer.RecepientName, utils.CurrencyCodeToCountryCode(body.Transfer.Currency)))
	}

	if user.NotificationOptions.Emails {
		background.Scheduler.Emit("send_email", map[string]any{
			"email":        user.Email,
			"subject":      "Your payment is on its way! 🚀",
			"templateName": "payment_recieved",
			"opts": map[string]any{
				"FIRSTNAME":         user.FirstName,
				"CURRENCY_CODE":     utils.CurrencyCodeToCurrencySymbol("NGN"),
				"AMOUNT":            currencyformatter.HumanReadableFloat32Currency(body.Transfer.Amount),
				"RECEPIENT_NAME":    body.Transfer.RecepientName,
				"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName("Nigeria"),
			},
		})

	}
	return nil
}
