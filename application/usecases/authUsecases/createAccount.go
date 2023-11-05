package authusecases

import (
	"errors"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/cryptography"
	"kego.com/application/repository"
	"kego.com/entities"
	"kego.com/infrastructure/validator"
)

func CreateAccount(ctx any, payload *entities.User)(*entities.User, error){
	valiedationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx, valiedationErr)
		return nil, errors.New("")
	}
	passwordHash, err := cryptography.CryptoHahser.HashString(payload.Password)
	if err != nil {
		apperrors.ValidationFailedError(ctx, &[]error{err})
		return nil, err
	}
	transactionPinHash, err := cryptography.CryptoHahser.HashString(payload.TransactionPin)
	if err != nil {
		apperrors.ValidationFailedError(ctx, &[]error{err})
		return nil, err
	}
	payload.Password = string(passwordHash)
	payload.TransactionPin = string(transactionPinHash)
	accountVerificationStatus := false
	if payload.Email != nil {
		payload.EmailVerified = &accountVerificationStatus
	}else {
		payload.PhoneVerified = &accountVerificationStatus
	}
	result, err := repository.UserRepo().CreateOne(*payload)
	if err != nil {
		apperrors.ValidationFailedError(ctx, &[]error{err})
		return nil, err
	}
	return result, err
}