package wallet

import (
	"context"
	"errors"
	"fmt"

	"kego.com/application/constants"
	"kego.com/application/repository"
	"kego.com/entities"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/validator"
)

func CreateWallet(ctx any, trxCtx context.Context, payload *entities.Wallet) (*entities.Wallet, error) {
	validationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if validationErr != nil {
		return nil, (*validationErr)[0]
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
		return nil, err
	}
	if walletCount == int64(constants.BUSINESS_WALLET_LIMIT) {
		err = fmt.Errorf("You have reached your wallet limit. If you think this is an error contact %s.", constants.SUPPORT_EMAIL)
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
		return nil, err
	}
	return wallet, nil
}