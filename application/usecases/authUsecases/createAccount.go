package authusecases

import (
	"errors"
	"fmt"

	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
	"kego.com/application/repository"
	"kego.com/entities"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/validator"
)

func CreateAccount(ctx any, payload *entities.User)(*entities.User, error){
	valiedationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx, valiedationErr)
		return nil, errors.New("")
	}
	bankExists := ""
	for _, bank := range bankssupported.SupportedBanks {
		if bank.Name == payload.BankDetails.BankName{
			bankExists = bank.Code
			break
		}
	}
	if bankExists  == "" {
		apperrors.NotFoundError(ctx, fmt.Sprintf("%s is not a supported bank on our platform yet.", payload.BankDetails.BankName))
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
	result, err := repository.UserRepo().CreateOne(nil, *payload)
	if err != nil {
		apperrors.EntityAlreadyExistsError(ctx, err.Error())
		return nil, err
	}
	return result,  err
}