package authusecases

import (
	"errors"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/cryptography"
	"kego.com/application/repository"
	"kego.com/infrastructure/validator"
	"kego.com/entities"
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
	payload.Password = string(passwordHash)
	result, err := repository.UserRepo().CreateOne(*payload)
	return result, err
}