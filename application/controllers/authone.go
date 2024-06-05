package controllers

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/repository"
	"usepolymer.co/infrastructure/background"
	identityverification "usepolymer.co/infrastructure/identity_verification"
	"usepolymer.co/infrastructure/logger"
	server_response "usepolymer.co/infrastructure/serverResponse"
	"usepolymer.co/infrastructure/validator"
)

func AuthOneFetchUserDetails(ctx *interfaces.ApplicationContext[any]) {
	userRepo := repository.UserRepo()
	user, err := userRepo.FindOneByByFilterStripped(map[string]interface{}{
		"email": ctx.Param["email"],
	}, options.FindOne().SetProjection(map[string]any{
		"_id": 1,
	}))
	if err != nil {
		logger.Error(errors.New("an error occured while fetching user details for authone"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "email",
			Data: ctx.Body,
		})
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if user == nil {
		apperrors.NotFoundError(ctx.Ctx, "User not found", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "account found", user, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func AuthOneSendEmail(ctx *interfaces.ApplicationContext[dto.AuthOneSendEmail]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	err := background.Scheduler.Emit("send_email", map[string]any{
		"email":        ctx.Body.Email,
		"subject":      ctx.Body.Subject,
		"templateName": ctx.Body.Template,
		"opts":         ctx.Body.Opts,
	})
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "email sent", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func AuthOneVerifyEmailStatus(ctx *interfaces.ApplicationContext[string]) {
	status, err := identityverification.IdentityVerifier.EmailVerification(*ctx.Body)
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, nil)
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "email sent", status, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func AuthOneCreateUser(ctx *interfaces.ApplicationContext[dto.AuthOneCreateUserDTO]) {
	CreateAccount(&interfaces.ApplicationContext[dto.CreateAccountDTO]{
		Body: &dto.CreateAccountDTO{
			Email:    ctx.Body.Email,
			Password: ctx.Body.Password,
			AuthOne:  true,
		},
		Ctx:    ctx.Ctx,
		Header: ctx.Header,
	})
}
