package controllers

import (
	"fmt"
	"net/http"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/constants"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	server_response "kego.com/infrastructure/serverResponse"
)


func FetchUserProfile(ctx *interfaces.ApplicationContext[any]){
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("USER_ID"))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if user == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("This user profile was not found. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "countries fetched", user, nil)
}