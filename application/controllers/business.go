package controllers

import (
	"net/http"

	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/usecases/business"
	"kego.com/entities"
	server_response "kego.com/infrastructure/serverResponse"
)


func CreateBusiness(ctx *interfaces.ApplicationContext[dto.BusinessDTO]){
	business, wallet, err := business.CreateBusiness(ctx.Ctx, &entities.Business{
		Name: ctx.Body.Name,
		UserID: ctx.GetStringContextData("UserID"),
	})
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "business created", map[string]any{
		"business": business,
		"wallet": wallet,
	}, nil)
}