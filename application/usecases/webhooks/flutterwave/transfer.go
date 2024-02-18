package flutterwave

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
	"kego.com/application/controllers/dto"
	"kego.com/application/repository"
	"kego.com/application/services"
	"kego.com/application/utils"
	"kego.com/infrastructure/logger"
	pushnotification "kego.com/infrastructure/messaging/push_notifications"
	currencyformatter "kego.com/infrastructure/currency_formatter"
	"kego.com/infrastructure/messaging/emails"
)

func FlwTransferWebhook(ctx any, body dto.FlutterwaveWebhookDTO) error {
		var err error
		if body.Transfer.Status == "SUCCESSFUL" {
			err = services.RemoveLockFunds(ctx, body.Transfer.Meta.WalletID, body.Transfer.Ref)
		}else if body.Transfer.Status == "FAILED" {
			err = services.ReverseLockFunds(ctx, body.Transfer.Meta.WalletID, body.Transfer.Ref)
		}
		if err != nil {
			logger.Error(errors.New("error removing locked funds after flutterwave webhook call"), logger.LoggerOptions{
				Key: "error",
				Data: err,
			}, logger.LoggerOptions{
				Key: "payload",
				Data: body,
			})
			return err
		}
		userRepo := repository.UserRepo()
		user, err := userRepo.FindByID(body.Transfer.Meta.UserID, options.FindOne().SetProjection(map[string]any{
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
				Data: body,
			})
			return err
		}
		if user == nil {
			err = errors.New("user account not found during flutterwave webhook run")
			logger.Error(err, logger.LoggerOptions{
				Key: "payload",
				Data: body,
			})
			return err
		}
		
		if user.NotificationOptions.PushNotification {
			pushnotification.PushNotificationService.PushOne(user.DeviceID, "Your payment was successful! 🚀",
				fmt.Sprintf("Your payment of %s%s to %s in %s has been processed successfully.", body.Transfer.Currency, currencyformatter.HumanReadableFloat32Currency(body.Transfer.Amount), body.Transfer.RecepientName, utils.CurrencyCodeToCountryCode(body.Transfer.Currency)))
		}
	
		if user.NotificationOptions.Emails {
			emails.EmailService.SendEmail(user.Email, "Your payment is on its way! 🚀", "payment_recieved", map[string]any{
				"FIRSTNAME": user.FirstName,
				"CURRENCY_CODE": utils.CurrencyCodeToCurrencySymbol("NGN"),
				"AMOUNT": currencyformatter.HumanReadableFloat32Currency(body.Transfer.Amount),
				"RECEPIENT_NAME": body.Transfer.RecepientName,
				"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName("Nigeria"),
			})
		}
	return nil
}
