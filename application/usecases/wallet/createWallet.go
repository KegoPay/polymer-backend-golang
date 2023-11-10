package wallet

import (
	"context"
	"errors"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/repository"
	"kego.com/entities"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/validator"
)

func CreateWallet(ctx any,trxCtx *context.Context, payload *entities.Wallet) (*entities.Wallet, error) {
	validationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if validationErr != nil {
		apperrors.ValidationFailedError(ctx, validationErr)
		return nil, errors.New("")
	}
	walletRepo := repository.WalletRepo()
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