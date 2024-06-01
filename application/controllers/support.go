package controllers

import (
	"context"
	"net/http"

	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/repository"
	"usepolymer.co/entities"
	server_response "usepolymer.co/infrastructure/serverResponse"
	"usepolymer.co/infrastructure/validator"
)

func ErrSupportRequest(ctx *interfaces.ApplicationContext[dto.ErrorSupportRequestDTO]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	errSupportRequestRepo := repository.ErrorSupportRequestRepo()
	_, err := errSupportRequestRepo.CreateOne(context.TODO(), entities.ErrorSupportRequest{
		UserID:  ctx.GetStringContextData("UserID"),
		Message: ctx.Body.Message,
		Email:   ctx.GetStringContextData("Email"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "support request sent", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}
