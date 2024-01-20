package controllers

import (
	"net/http"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/application/usecases/business"
	"kego.com/entities"
	server_response "kego.com/infrastructure/serverResponse"
)


func CreateBusiness(ctx *interfaces.ApplicationContext[dto.BusinessDTO]){
	business, wallet, err := business.CreateBusiness(ctx.Ctx, &entities.Business{
		Name: ctx.Body.Name,
		UserID: ctx.GetStringContextData("UserID"),
		Email: ctx.Body.Email,
	})
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "business created", map[string]any{
		"business": business,
		"wallet": wallet,
	}, nil)
}

func UpdateBusiness(ctx *interfaces.ApplicationContext[dto.UpdateBusinessDTO]){
	ctx.Body.ID = ctx.GetStringParameter("businessID")
	err := business.UpdateBusiness(ctx.Ctx, ctx.Body)
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business updated", nil, nil)
}

func DeleteBusiness(ctx *interfaces.ApplicationContext[any]){
	err := business.DeleteBusiness(ctx.Ctx, ctx.GetStringParameter("businessID"))
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business deleted", nil, nil)
}

func FetchBusinesses(ctx *interfaces.ApplicationContext[any]){
	businessRepo := repository.BusinessRepo()
	businesses, err := businessRepo.FindMany(map[string]interface{}{
		"userID": ctx.GetStringContextData("UserID"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "businesses fetched", businesses, nil)
}