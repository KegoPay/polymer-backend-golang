package userusecases

import (
	"errors"
	"fmt"
	"strings"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers/dto"
	"kego.com/application/repository"
	"kego.com/infrastructure/validator"
)

func UpdateUserTag(ctx any, id string, tag dto.SetPaymentTagDTO) error {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(tag)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx, valiedationErr)
		return errors.New("")
	}
	userRepo := repository.UserRepo()
	tag.Tag = strings.ToLower(tag.Tag)
	tagExists, err := userRepo.CountDocs(map[string]interface{}{
		"tag": tag.Tag,
	})
	if err != nil {
		apperrors.FatalServerError(ctx, err)
		return err
	}
	if tagExists != 0 {
		err = fmt.Errorf("payment tag %s is associated with another account", tag.Tag)
		apperrors.EntityAlreadyExistsError(ctx, err.Error())
		return err
	}
	count, err := userRepo.UpdatePartialByID(id, map[string]any{
		"tag": tag.Tag,
	})
	if err != nil {
		apperrors.FatalServerError(ctx, err)
		return err
	}
	if count != 1 {
		apperrors.UnknownError(ctx, fmt.Errorf(`could not update user tag "%s" id - %s`, tag.Tag, id))
		return errors.New("")
	}
	return nil
}