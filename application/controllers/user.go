package controllers

import (
	"errors"
	"fmt"
	"net/http"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/constants"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/infrastructure/logger"
	server_response "kego.com/infrastructure/serverResponse"
	"kego.com/infrastructure/validator"
)


func FetchUserProfile(ctx *interfaces.ApplicationContext[any]){
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if user == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("This user profile was not found. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "profile fetched", user, nil)
}

func UpdateUserProfile(ctx *interfaces.ApplicationContext[dto.UpdateUserDTO]){userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"))
	if err != nil {
		logger.Error(errors.New("error fetching user profile for update"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if ctx.Body.FirstName != nil {
		user.FirstName = *ctx.Body.FirstName
	}
	if ctx.Body.LastName != nil {
		user.LastName = *ctx.Body.LastName
	}
	if ctx.Body.Phone != nil {
		user.Phone = *ctx.Body.Phone
	}
	validationErr := validator.ValidatorInstance.ValidateStruct(user)
	if validationErr != nil {
		apperrors.ValidationFailedError(ctx, validationErr)
		return
	}
	userRepo.UpdateByID(ctx.GetStringContextData("UserID"), user)
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "update completed", nil, nil)
}