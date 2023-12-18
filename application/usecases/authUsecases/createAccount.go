package authusecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
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
	bankExists := ""
	for _, bank := range bankssupported.KYCSupportedBanks {
		if bank.Name == payload.BankDetails.BankName{
			bankExists = bank.Code
			break
		}
	}
	if bankExists  == "" {
		apperrors.NotFoundError(ctx, fmt.Sprintf("%s is not a supported bank on our platform yet.", payload.BankDetails.BankName))
		return nil, nil,  errors.New("")
	}
	passwordHash, err := cryptography.CryptoHahser.HashString(payload.Password)
	if err != nil {
		apperrors.ValidationFailedError(ctx, &[]error{err})
		return nil, nil, err
	}
	transactionPinHash, err := cryptography.CryptoHahser.HashString(payload.TransactionPin)
	if err != nil {
		apperrors.ValidationFailedError(ctx, &[]error{err})
		return nil, nil, err
	}
	payload.Password = string(passwordHash)
	payload.TransactionPin = string(transactionPinHash)
	userRepo := repository.UserRepo()
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