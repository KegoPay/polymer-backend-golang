package authusecases

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/repository"
	walletUsecases "usepolymer.co/application/usecases/wallet"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/cryptography"
	"usepolymer.co/infrastructure/database/repository/cache"
	identityverification "usepolymer.co/infrastructure/identity_verification"
	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/validator"
)

func CreateAccount(ctx any, payload *entities.User, device_id *string, verifiedEmail bool) (*entities.User, *entities.Wallet, error) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx, valiedationErr, device_id)
		return nil, nil, errors.New("")
	}
	userRepo := repository.UserRepo()
	passwordHash, err := cryptography.CryptoHahser.HashString(payload.Password, nil)
	if err != nil {
		apperrors.ValidationFailedError(ctx, &[]error{err}, device_id)
		return nil, nil, err
	}
	exists, err := userRepo.CountDocs(map[string]any{
		"email": payload.Email,
	})
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, nil, err
	}
	if exists != 0 {
		err = errors.New("user with email already exists")
		apperrors.EntityAlreadyExistsError(ctx, err.Error(), device_id)
		return nil, nil, err
	}
	if os.Getenv("ENV") == "prod" {
		found := cache.Cache.FindOne(fmt.Sprintf("%s-email-blacklist", payload.Email))
		if found != nil {
			apperrors.ClientError(ctx, fmt.Sprintf("%s was not approved for signup on Polymer", payload.Email), nil, nil, device_id)
			return nil, nil, err
		}
		valid, err := identityverification.IdentityVerifier.EmailVerification(payload.Email)
		if err != nil {
			apperrors.FatalServerError(ctx, err, device_id)
			return nil, nil, err
		}
		if !valid {
			apperrors.ClientError(ctx, fmt.Sprintf("%s was not approved for signup on Polymer", payload.Email), nil, nil, device_id)
			cache.Cache.CreateEntry(fmt.Sprintf("%s-email-blacklist", payload.Email), payload.Email, time.Minute*0)
			return nil, nil, err
		}
	}

	payload.Password = string(passwordHash)
	var user *entities.User
	var wallet *entities.Wallet
	userRepo.StartTransaction(func(sc mongo.Session, c context.Context) error {
		userPayload := payload.ParseModel().(*entities.User)
		walletPayload := &entities.Wallet{
			UserID:         userPayload.ID,
			Frozen:         false,
			Balance:        0,
			LedgerBalance:  0,
			Currency:       "NGN",
			LockedFundsLog: []entities.LockedFunds{},
		}

		walletPayload = walletPayload.ParseModel().(*entities.Wallet)
		userPayload.WalletID = walletPayload.ID
		userPayload.EmailVerified = verifiedEmail
		userPayload.NotificationOptions = entities.NotificationOptions{
			Emails:           true,
			PushNotification: true,
		}
		userData, e := userRepo.CreateOne(c, *userPayload)
		if e != nil {
			logger.Error(errors.New("error creating users account"), logger.LoggerOptions{
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
		user = userData
		wallet, e = walletUsecases.CreateWallet(ctx, c, walletPayload)
		if e != nil {
			logger.Error(errors.New("error creating users wallet"), logger.LoggerOptions{
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
	return user, wallet, err
}
