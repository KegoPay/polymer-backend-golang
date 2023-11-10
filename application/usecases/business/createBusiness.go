package business

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/repository"
	walletUsecases "kego.com/application/usecases/wallet"
	"kego.com/entities"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/validator"
)

func CreateBusiness(ctx any, payload *entities.Business) (*entities.Business, *entities.Wallet, error) {
	validationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if validationErr != nil {
		apperrors.ValidationFailedError(ctx, validationErr)
		return nil, nil, errors.New("")
	}
	businessRepo := repository.BusinessRepo()
	var business *entities.Business
	var wallet *entities.Wallet
	var err error
	businessRepo.StartTransaction(func(sc *mongo.SessionContext, c *context.Context) error {
		payload = payload.ParseModel().(*entities.Business)
		walletPayload := &entities.Wallet{
			BusinessID: payload.ID,
			UserID: payload.UserID,
			Frozen: false,
			Balance: 0,
			AvailableBalance: 0,
		}
		walletPayload = walletPayload.ParseModel().(*entities.Wallet)
		payload.WalletID = walletPayload.ID
		b, e := businessRepo.CreateOne(c, *payload)
		if e != nil {
			logger.Error(errors.New("error creating users business"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "payload",
				Data: payload,
			})
			err = e
			return e
		}
		business = b
		w, e := walletUsecases.CreateWallet(ctx, c, walletPayload)
		if e != nil {
			logger.Error(errors.New("error creating users business"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "payload",
				Data: walletPayload,
			})
			err = e
			return e
		}
		walletPayload = w
		return nil
	})
	if err != nil {
		apperrors.FatalServerError(ctx)
		return nil, nil, err
	}
	return business, wallet, nil
}