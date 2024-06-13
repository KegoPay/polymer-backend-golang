package business

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/repository"
	walletUsecases "usepolymer.co/application/usecases/wallet"
	"usepolymer.co/application/utils"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/validator"
)

func CreateBusiness(ctx any, payload *entities.Business, device_id *string) (*entities.Business, *entities.Wallet, error) {
	validationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if validationErr != nil {
		apperrors.ValidationFailedError(ctx, validationErr, device_id)
		return nil, nil, errors.New("")
	}
	businessRepo := repository.BusinessRepo()
	var business *entities.Business
	var wallet *entities.Wallet
	var err error
	businessRepo.StartTransaction(func(sc mongo.Session, c context.Context) error {
		payload = payload.ParseModel().(*entities.Business)
		walletPayload := &entities.Wallet{
			BusinessID:     &payload.ID,
			UserID:         payload.UserID,
			BusinessName:   &payload.Name,
			Frozen:         false,
			Balance:        0,
			LedgerBalance:  0,
			Currency:       "NGN",
			LockedFundsLog: []entities.LockedFunds{},
		}
		walletPayload = walletPayload.ParseModel().(*entities.Wallet)
		payload.WalletID = walletPayload.ID
		b, e := businessRepo.CreateOne(c, *payload)
		if e != nil {
			logger.Error(errors.New("error creating users business"), logger.LoggerOptions{
				Key:  "error",
				Data: e,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: payload,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		business = b
		w, e := walletUsecases.CreateWallet(ctx, c, walletPayload)
		if e != nil {
			logger.Error(errors.New("error creating users business wallet"), logger.LoggerOptions{
				Key:  "error",
				Data: e,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: walletPayload,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		userRepo := repository.UserRepo()
		affected, e := userRepo.UpdatePartialByID(b.UserID, map[string]any{
			"hasBusiness": true,
		})
		if e != nil {
			logger.Error(errors.New("error updated users hasBusiness field after business creation"), logger.LoggerOptions{
				Key:  "error",
				Data: e,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: business,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		if affected == 0 {
			e = errors.New("error updated users hasBusiness field after business creation")
			logger.Error(e, logger.LoggerOptions{
				Key:  "error",
				Data: e,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: business,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		wallet = w
		(sc).CommitTransaction(c)
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			apperrors.EntityAlreadyExistsError(ctx, err.Error(), device_id)
			return nil, nil, err
		} else {
			apperrors.ClientError(ctx, err.Error(), nil, nil, device_id)
			return nil, nil, err
		}
	}
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(payload.UserID, options.FindOne().SetProjection(map[string]any{
		"bvn": 1,
		"nin": 1,
	}))
	if err != nil {
		logger.Error(errors.New("error fetching user bvn to create business account dva"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: payload,
		})
		return nil, nil, err
	}
	if user.NIN != "" {
		wallet.AccountNumber, wallet.BankName, err = walletUsecases.GenerateNGNDVA(ctx, business.WalletID, business.Name, "Polymer Software", business.Email, *utils.GenerateDummyKYCID(), user.NIN, device_id)
	} else {
		wallet.AccountNumber, wallet.BankName, err = walletUsecases.GenerateNGNDVA(ctx, business.WalletID, business.Name, "Polymer Software", business.Email, user.BVN, *utils.GenerateDummyKYCID(), device_id)
	}
	if err != nil {
		return business, wallet, nil
	}
	return business, wallet, nil
}
