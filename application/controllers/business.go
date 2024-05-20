package controllers

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/constants"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/application/usecases/business"
	"kego.com/entities"
	"kego.com/infrastructure/background"
	cac_service "kego.com/infrastructure/cac"
	"kego.com/infrastructure/logger"
	server_response "kego.com/infrastructure/serverResponse"
	"kego.com/infrastructure/validator"
)

func CreateBusiness(ctx *interfaces.ApplicationContext[dto.BusinessDTO]) {
	ctx.Body.Email = ctx.GetStringContextData("Email")
	business, wallet, err := business.CreateBusiness(ctx.Ctx, &entities.Business{
		Name:   ctx.Body.Name,
		UserID: ctx.GetStringContextData("UserID"),
		Email:  ctx.Body.Email,
	}, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "business created", map[string]any{
		"business": business,
		"wallet":   wallet,
	}, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func UpdateBusiness(ctx *interfaces.ApplicationContext[dto.UpdateBusinessDTO]) {
	ctx.Body.ID = ctx.GetStringParameter("businessID")
	err := business.UpdateBusiness(ctx.Ctx, ctx.Body, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business updated", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func DeleteBusiness(ctx *interfaces.ApplicationContext[any]) {
	err := business.DeleteBusiness(ctx.Ctx, ctx.GetStringParameter("businessID"), ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business deleted", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func FetchBusinesses(ctx *interfaces.ApplicationContext[any]) {
	businessRepo := repository.BusinessRepo()
	business, err := businessRepo.FindOneByFilter(map[string]interface{}{
		"userID": ctx.GetStringContextData("UserID"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if business == nil {
		server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business fetched", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	walletRepo := repository.WalletRepo()
	wallet, err := walletRepo.FindOneByFilter(map[string]interface{}{
		"businessID": business.ID,
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business fetched", map[string]any{
		"business": business,
		"wallet":   wallet,
	}, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func SearchCACByName(ctx *interfaces.ApplicationContext[dto.SearchCACByName]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	names, err := cac_service.CACServiceInstance.FetchBusinessDetailsByName(ctx.Body.Name)
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business fetched", names, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func SetCACInfo(ctx *interfaces.ApplicationContext[dto.SetCACInfo]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	names, err := cac_service.CACServiceInstance.FetchBusinessDetailsByName(ctx.Body.RCNumber)
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	businessRepo := repository.BusinessRepo()
	business, err := businessRepo.FindByID(ctx.GetStringParameter("businessID"), options.FindOne().SetProjection(map[string]any {
		"cacInfo": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if business.CACInfo.Verified {
		apperrors.ClientError(ctx.Ctx, "Your business has already been verified! If you wish to update it's details please contact support", nil, &constants.BUSINESS_ALREADY_VERIFIED,  ctx.GetHeader("Polymer-Device-Id"))
	}
	updated, err := businessRepo.UpdatePartialByID(ctx.GetStringParameter("businessID"), map[string]any{
		"cacInfo": (*names)[0],
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if updated != 1 {
		logger.Error(errors.New("something went wrong with SetCACInfo"), logger.LoggerOptions{
			Key:  "updated",
			Data: updated,
		}, logger.LoggerOptions{
			Key:  "businessID",
			Data: ctx.GetStringParameter("businessID"),
		})
		return
	}
	background.Scheduler.Emit("verify_business", map[string]any{
		"id": ctx.GetStringParameter("businessID"),
	})
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business details saved", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}
