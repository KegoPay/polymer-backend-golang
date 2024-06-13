package wallet

import (
	"errors"
	"fmt"
	"os"

	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/repository"
	"usepolymer.co/application/services"
	"usepolymer.co/infrastructure/logger"
	kora_local_payment_processor "usepolymer.co/infrastructure/payment_processor/kora"
)

func GenerateNGNDVA(ctx any, walletID string, firstName string, lastName string, email string, bvn string, nin string, device_id *string) (accountNumber *string, bankName *string, err error) {
	dva := services.GenerateDVA(ctx, &kora_local_payment_processor.CreateVirtualAccountPayload{
		Reference: walletID,
		Name:      fmt.Sprintf("%s %s", firstName, lastName),
		Permanent: true,
		BankCode: func() string {
			if os.Getenv("ENV") != "prod" {
				return "000"
			}
			return "035"
		}(),
		Customer: kora_local_payment_processor.DVACustomer{
			Name:  fmt.Sprintf("%s %s", firstName, lastName),
			Email: email,
		},
		KYC: kora_local_payment_processor.DVAKYC{
			BVN: bvn,
			NIN: nin,
		},
	}, device_id)
	walletRepo := repository.WalletRepo()
	affected, err := walletRepo.UpdatePartialByID(walletID, map[string]any{
		"accountNumber": dva.AccountNumber,
		"bankName":      dva.BankName,
	})
	if err != nil {
		logger.Error(errors.New("failed to update user wallet with dva details"), logger.LoggerOptions{
			Key:  "err",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "walletID",
			Data: walletID,
		}, logger.LoggerOptions{
			Key:  "dva",
			Data: dva,
		})
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, nil, err
	}
	if affected == 0 {
		logger.Error(errors.New("failed to update user wallet with dva details"), logger.LoggerOptions{
			Key:  "walletID",
			Data: walletID,
		}, logger.LoggerOptions{
			Key:  "dva",
			Data: dva,
		}, logger.LoggerOptions{
			Key:  "affected",
			Data: affected,
		})
		apperrors.UnknownError(ctx, errors.New("attempt to updated a users wallet with DVA details failed"), device_id)
		return nil, nil, err
	}
	return &dva.AccountNumber, &dva.BankName, nil
}
