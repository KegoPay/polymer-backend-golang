package controllers

import (
	"context"
	"net/http"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/entities"
	server_response "kego.com/infrastructure/serverResponse"
	"kego.com/infrastructure/validator"
)

func ErrSupportRequest(ctx *interfaces.ApplicationContext[dto.ErrorSupportRequestDTO]){
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	errSupportRequestRepo := repository.ErrorSupportRequestRepo()
	_, err := errSupportRequestRepo.CreateOne(context.TODO(), entities.ErrorSupportRequest{
		UserID: ctx.GetStringContextData("UserID"),
		Message: ctx.Body.Message,
		Email: ctx.GetStringContextData("Email"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "support request sent", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}
