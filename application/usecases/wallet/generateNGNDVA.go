package wallet

import (
	"errors"
	"fmt"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/repository"
	"kego.com/application/services"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/payment_processor/types"
)

func GenerateNGNDVA(ctx any, walletID string, firstName string, lastName string, email string, bvn string) (accountNumber *string, bankName *string, err error) {
	dva := services.GenerateDVA(ctx, &types.CreateVirtualAccountPayload{
		Permanent: true,
		Currency: "NGN",
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		TransactionReference: walletID,
		Narration: fmt.Sprintf("%s %s", firstName, lastName),
		BVN: bvn,
	})
	walletRepo := repository.WalletRepo()
	affected, err := walletRepo.UpdatePartialByID(walletID, map[string]any{
		"accountNumber": dva.AccountNumber,
		"bankName": dva.BankName,
	})
	if err != nil {
		logger.Error(errors.New("failed to update user wallet with dva details"), logger.LoggerOptions{
			Key: "err",
			Data: err,
		}, logger.LoggerOptions{
			Key: "walletID",
			Data: walletID,
		}, logger.LoggerOptions{
			Key: "dva",
			Data: dva,
		})
		apperrors.FatalServerError(ctx, err)
		return nil, nil, err
	}
	if affected == 0 {
		logger.Error(errors.New("failed to update user wallet with dva details"), logger.LoggerOptions{
			Key: "walletID",
			Data: walletID,
		}, logger.LoggerOptions{
			Key: "dva",
			Data: dva,
		}, logger.LoggerOptions{
			Key: "affected",
			Data: affected,
		})
		apperrors.UnknownError(ctx, errors.New("attempt to updated a users wallet with DVA details failed"))
		return nil, nil, err
	}
	return &dva.AccountNumber, &dva.BankName, nil
}