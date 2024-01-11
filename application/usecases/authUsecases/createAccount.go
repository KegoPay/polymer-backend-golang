package authusecases

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/repository"
	walletUsecases "kego.com/application/usecases/wallet"
	"kego.com/entities"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/validator"
)

func CreateAccount(ctx any, payload *entities.User)(*entities.User, *entities.Wallet, error){
	valiedationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx, valiedationErr)
		return nil, nil, errors.New("")
	}
	userRepo := repository.UserRepo()
	bvnExists, err := userRepo.CountDocs(map[string]any{
		"bvn": payload.BVN,
	})
	if err != nil {
		apperrors.FatalServerError(ctx)
		return nil, nil, err
	}
	if bvnExists != 0 {
		err = errors.New("bvn is already registered to another account")
		apperrors.EntityAlreadyExistsError(ctx, err.Error())
		return nil, nil, err
	}
	passwordHash, err := cryptography.CryptoHahser.HashString(payload.Password)
	if err != nil {
		apperrors.ValidationFailedError(ctx, &[]error{err})
		return nil, nil, err
	}
	payload.Password = string(passwordHash)
	var user *entities.User
	var wallet *entities.Wallet
	userRepo.StartTransaction(func(sc mongo.Session, c context.Context) error {
		userPayload := payload.ParseModel().(*entities.User)
		walletPayload := &entities.Wallet{
			UserID: userPayload.ID,
			Frozen: false,
			Balance: 0,
			LedgerBalance: 0,
			Currency: "NGN",
			LockedFundsLog: []entities.LockedFunds{},
		}

		walletPayload = walletPayload.ParseModel().(*entities.Wallet)
		userPayload.WalletID = walletPayload.ID
		userPayload.NotificationOptions = entities.NotificationOptions{
			Emails: true,
			PushNotification: true,
		}
		userData, e := userRepo.CreateOne(c, *userPayload)
		if e != nil {
			logger.Error(errors.New("error creating users account"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "payload",
				Data: payload,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		user = userData
		wallet, e = walletUsecases.CreateWallet(ctx, c, walletPayload)
		if e != nil {
			logger.Error(errors.New("error creating users wallet"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "payload",
				Data: walletPayload,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		(sc).CommitTransaction(c)
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists"){
			apperrors.EntityAlreadyExistsError(ctx, err.Error())
			return nil, nil, err
		}else {
			apperrors.ClientError(ctx, err.Error(), nil)
			return nil, nil, err
		}
	}
	return user, wallet,  err
}