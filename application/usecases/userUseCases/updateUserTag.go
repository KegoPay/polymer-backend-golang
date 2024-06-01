package userusecases

import (
	"errors"
	"fmt"
	"strings"

	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/repository"
	"usepolymer.co/infrastructure/validator"
)

func UpdateUserTag(ctx any, id string, tag dto.SetPaymentTagDTO, device_id *string) error {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(tag)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx, valiedationErr, device_id)
		return errors.New("")
	}
	userRepo := repository.UserRepo()
	tag.Tag = strings.ToLower(tag.Tag)
	tagExists, err := userRepo.CountDocs(map[string]interface{}{
		"tag": tag.Tag,
	})
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return err
	}
	if tagExists != 0 {
		err = fmt.Errorf("payment tag %s is associated with another account", tag.Tag)
		apperrors.EntityAlreadyExistsError(ctx, err.Error(), device_id)
		return err
	}
	count, err := userRepo.UpdatePartialByID(id, map[string]any{
		"tag": tag.Tag,
	})
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return err
	}
	if count != 1 {
		apperrors.UnknownError(ctx, fmt.Errorf(`could not update user tag "%s" id - %s`, tag.Tag, id), device_id)
		return errors.New("")
	}
	return nil
}
