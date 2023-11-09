package controllers

import (
	"net/http"
	"time"

	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	authusecases "kego.com/application/usecases/authUsecases"
	"kego.com/entities"
	"kego.com/infrastructure/auth"
	"kego.com/infrastructure/database/repository/cache"
	server_response "kego.com/infrastructure/serverResponse"
)

func CreateAccount(ctx *interfaces.ApplicationContext[dto.CreateAccountDTO]) {
	account, err := authusecases.CreateAccount(ctx.Ctx, &entities.User{
		Email:          ctx.Body.Email,
		Phone:          ctx.Body.Phone,
		Password:       ctx.Body.Password,
		TransactionPin: ctx.Body.TransactionPin,
		DeviceType:     ctx.Body.DeviceType,
		DeviceID:       ctx.Body.DeviceID,
		FirstName:      ctx.Body.FirstName,
		LastName:       ctx.Body.LastName,
	})
	if err != nil {
		return
	}
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email:     account.Email,
		Phone:     account.Phone,
		UserID:    account.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(10)).Unix(), //lasts for 10 mins
		DeviceType: account.DeviceType,
		DeviceID:   account.DeviceID,
	})
	cache.Cache.CreateEntry(account.ID, *token, time.Minute * time.Duration(10)) // cache authentication token for 10 mins
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "account created", map[string]interface{}{
		"account": account,
		"token":   token,
	}, nil)
}

func LoginUser(ctx *interfaces.ApplicationContext[dto.LoginDTO]){
	account, token := authusecases.LoginAccount(ctx.Ctx, ctx.Body.Email, ctx.Body.Phone, &ctx.Body.Password)
	if account == nil || token == nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "login successful", map[string]interface{}{
		"account": account,
		"token":   token,
	}, nil)
}
