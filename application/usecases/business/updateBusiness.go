package business

import (
	"errors"

	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/repository"
	"usepolymer.co/infrastructure/validator"
)

func UpdateBusiness(ctx any, payload *dto.UpdateBusinessDTO, device_id *string) error {
	validationErr := validator.ValidatorInstance.ValidateStruct(*payload)
	if validationErr != nil {
		apperrors.ValidationFailedError(ctx, validationErr, device_id)
		return (*validationErr)[0]
	}

	businessRepo := repository.BusinessRepo()
	success, err := businessRepo.UpdatePartialByID(payload.ID, map[string]any{
		"name": payload.Name,
	})
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return err
	}
	if success == 0 {
		apperrors.NotFoundError(ctx, "business does not exist", device_id)
		return errors.New("")
	}
	return nil
}
