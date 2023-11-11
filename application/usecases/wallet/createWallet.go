package wallet

import (
	"context"
	"errors"
	"fmt"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/constants"
	"kego.com/application/repository"
	"kego.com/entities"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/validator"
)

func CreateWallet(ctx any, trxCtx *context.Context, payload *entities.Wallet) (*entities.Wallet, error) {
	validationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if validationErr != nil {
		apperrors.ValidationFailedError(ctx, validationErr)
		return nil, errors.New("")
	}
	walletRepo := repository.WalletRepo()
	walletCount, err := walletRepo.CountDocs(map[string]interface{}{
		"userID": payload.UserID,
	})
	if err != nil {
		logger.Error(errors.New("error fetching number of wallets a user has in wallet creation"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx)
		return nil, err
	}
	if walletCount == int64(constants.BUSINESS_WALLET_LIMIT) {
		err = fmt.Errorf("You have reached your wallet limit. If you think this is an error contact %s.", constants.SUPPORT_EMAIL)
		apperrors.ClientError(ctx, err.Error(), nil)
		return nil, err
	}
	wallet, err := walletRepo.CreateOne(trxCtx, *payload)
	if err != nil {
		logger.Error(errors.New("error creating users wallet"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "payload",
			Data: payload,
		})
		apperrors.FatalServerError(ctx)
		return nil, err
	}
	return wallet, nil
}