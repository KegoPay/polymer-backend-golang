package wallet

import (
	"errors"
	"fmt"
	"os"

	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/repository"
	"usepolymer.co/application/services"
	"usepolymer.co/application/utils"
	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/payment_processor/types"
)

func GenerateNGNDVA(ctx any, walletID string, firstName string, lastName string, email string, bvn string, device_id *string) (accountNumber *string, bankName *string, err error) {
	dva := services.GenerateDVA(ctx, &types.CreateVirtualAccountPayload{
		Permanent:            os.Getenv("ENV") == "production",
		Currency:             "NGN",
		FirstName:            firstName,
		LastName:             lastName,
		Email:                email,
		TransactionReference: walletID,
		Narration:            fmt.Sprintf("%s %s", firstName, lastName),
		BVN:                  bvn,
		Amount: func() *uint64 {
			if os.Getenv("ENV") != "production" {
				return utils.GetUInt64Pointer(10000000)
			}
			return utils.GetUInt64Pointer(0)
		}(),
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
