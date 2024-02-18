package flutterwave

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
	"kego.com/application/controllers/dto"
	"kego.com/application/repository"
	"kego.com/application/services"
	"kego.com/application/utils"
	"kego.com/entities"
	currencyformatter "kego.com/infrastructure/currency_formatter"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
	pushnotification "kego.com/infrastructure/messaging/push_notifications"
)

func CreditWebHook(body dto.FlutterwaveWebhookDTO) error {
		var err error
		userRepo := repository.UserRepo()
		user, err := userRepo.FindByID(body.Customer.Email, options.FindOne().SetProjection(map[string]any{
			"id": 1,
			"pushNotificationToken": 1,
			"notificationOptions": 1,
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
		if *body.Status == "successful" {
			services.CreditWallet(*body.TrxRef, utils.Float32ToUint64Currency(body.Transfer.Amount), entities.FlutterwaveDVACredit, &entities.Transaction{
				TransactionReference: *body.TrxRef,
				Amount: utils.Float32ToUint64Currency(body.Transfer.Amount),
				AmountInNGN: utils.Float32ToUint64Currency(body.Transfer.Amount),
				Fee: 0,
				ProcessorFee: 0,
				Currency: *body.Currency,
				ProcessorFeeCurrency: *body.Currency,
				WalletID: *body.TrxRef,
				UserID: user.ID,
				Description: fmt.Sprintf("Transfer from %s %s", body.Entity.FirstName, body.Entity.LastName),
				Location: entities.Location{
					IPAddress: *body.IPAddress,
				},
				Intent: entities.FlutterwaveDVACredit,
				Sender: entities.TransactionSender{
					FullName: fmt.Sprintf("%s %s", body.Entity.FirstName, body.Entity.LastName),
					AccountNumber: body.Entity.AccNumber,
				},
				Recepient: entities.TransactionRecepient{
					FullName: body.Customer.FullName,
				},
			})
		}else if body.Transfer.Status == "FAILED" {
			err = errors.New("credit attempt failed")
		}
		if err == nil {
			if user.NotificationOptions.PushNotification {
				pushnotification.PushNotificationService.PushOne(user.DeviceID, "Money In!🤪",
					fmt.Sprintf("You just got sent ₦%s by %s", currencyformatter.HumanReadableFloat32Currency(body.Transfer.Amount), body.Entity.FirstName))
			}
			if user.NotificationOptions.Emails {
				emails.EmailService.SendEmail(user.Email, "Money In!", "credit", map[string]any{
					"FIRSTNAME": user.FirstName,
					"CURRENCY_CODE": "₦",
					"AMOUNT": currencyformatter.HumanReadableFloat32Currency(body.Transfer.Amount),
					"RECEPIENT_NAME": body.Transfer.RecepientName,
				})
			}
			return nil
		}
		logger.Error(err, logger.LoggerOptions{
			Key: "payload",
			Data: body,
		})
	return err
}