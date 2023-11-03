package controllers

import (
	"net/http"

	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	authusecases "kego.com/application/usecases/authUsecases"
	"kego.com/entities"
	server_response "kego.com/infrastructure/serverResponse"
)

func CreateAccount(ctx *interfaces.ApplicationContext[dto.CreateAccountDTO]){
	result, err := authusecases.CreateAccount(ctx.Ctx, &entities.User{
		Email: ctx.Body.Email,
		Phone: ctx.Body.Phone,
		Password: ctx.Body.Password,
	})
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "account created", *result, nil)
}