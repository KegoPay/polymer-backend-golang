package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/constants"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/repository"
	"usepolymer.co/application/usecases/business"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/auth"
	"usepolymer.co/infrastructure/background"
	cac_service "usepolymer.co/infrastructure/cac"
	"usepolymer.co/infrastructure/database/repository/cache"
	"usepolymer.co/infrastructure/logger"
	pushnotification "usepolymer.co/infrastructure/messaging/push_notifications"
	server_response "usepolymer.co/infrastructure/serverResponse"
	"usepolymer.co/infrastructure/validator"
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
	businessRepo := repository.BusinessRepo()
	business, err := businessRepo.FindByID(ctx.GetStringParameter("businessID"), options.FindOne().SetProjection(map[string]any{
		"cacInfo": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if business.CACInfo != nil {
		if business.CACInfo.Verified {
			apperrors.ClientError(ctx.Ctx, "Your business has already been verified! If you wish to update it's details please contact support", nil, &constants.BUSINESS_ALREADY_VERIFIED, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
	}
	rcExists, err := businessRepo.CountDocs(map[string]interface{}{
		"cacInfo.rcNumber": ctx.Body.RCNumber,
		"verified":         true,
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if rcExists != 0 {
		apperrors.ClientError(ctx.Ctx, "This business profile has been attached to and verified on another business on Polymer. If you think this is a mistake please click the button below and we will review the situation", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	names, err := cac_service.CACServiceInstance.FetchBusinessDetailsByName(ctx.Body.RCNumber)
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
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

func VerifyBusinessManual(ctx *interfaces.ApplicationContext[any]) {
	token := ctx.Query["token"]
	if token == nil {
		apperrors.ClientError(ctx.Ctx, "Business verification faileds. Request a new token from the app", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	tokenExists := cache.Cache.FindOne(token.(string))
	if tokenExists != nil {
		apperrors.ClientError(ctx.Ctx, "This token has already been used in verifying your business. If you think this is a mistake and want a manual review please click the button below", nil, &constants.ESCALATE_TO_SUPPORT, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	decoded, err := auth.DecodeAuthToken(token.(string))
	if err != nil {
		apperrors.ClientError(ctx.Ctx, "Business verification faileds. Request a new token from the app", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !decoded.Valid {
		apperrors.AuthenticationError(ctx.Ctx, "Business verification faileds. Request a new token from the app", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	auth_token_claims := decoded.Claims.(jwt.MapClaims)
	if auth_token_claims["iss"] != os.Getenv("JWT_ISSUER") {
		background.Scheduler.Emit("lock_account", map[string]any{
			"id": auth_token_claims["userID"],
		})
		apperrors.AuthenticationError(ctx.Ctx, "Business verification faileds. Request a new token from the app", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(auth_token_claims["userID"].(string), options.FindOne().SetProjection(map[string]any{
		"deactivated":           1,
		"pushNotificationToken": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "this account no longer exists", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account.Deactivated {
		apperrors.AuthenticationError(ctx.Ctx, "account has been deactivated", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	directorCache := cache.Cache.FindOneByteArray(fmt.Sprintf("%s-kyc-info-directors", auth_token_claims["businessID"].(string)))
	shareHoldersCache := cache.Cache.FindOneByteArray(fmt.Sprintf("%s-kyc-info-shareholders", auth_token_claims["businessID"].(string)))
	address := cache.Cache.FindOne(fmt.Sprintf("%s-kyc-info-address", auth_token_claims["businessID"].(string)))
	if directorCache == nil || shareHoldersCache == nil || address == nil {
		apperrors.ClientError(ctx.Ctx, "Business kyc details were not found in our system. This is probaly because this link is too old to use and the data has expired. If you think this is a mistake and want a manual review please click the button below", nil, &constants.ESCALATE_TO_SUPPORT, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	var directors []entities.Director
	var shareholders []entities.ShareHolder
	err = json.Unmarshal(*directorCache, &directors)
	if err != nil {
		logger.Error(err, logger.LoggerOptions{
			Key:  "directors",
			Data: directorCache,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	err = json.Unmarshal(*shareHoldersCache, &shareholders)
	if err != nil {
		logger.Error(err, logger.LoggerOptions{
			Key:  "shareholders",
			Data: shareHoldersCache,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	businessRepo := repository.BusinessRepo()
	updated, err := businessRepo.UpdatePartialByID(auth_token_claims["businessID"].(string), map[string]any{
		"cacInfo.verified":    true,
		"shareholders":        shareholders,
		"directors":           directors,
		"cacInfo.fulladdress": address,
	})
	if err != nil {
		logger.Error(errors.New("an error occured while verifying cac information manually"), logger.LoggerOptions{
			Key:  "err",
			Data: err,
		})
		apperrors.AuthenticationError(ctx.Ctx, "Business verification faileds. Request a new token from the app", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if updated != 1 {
		logger.Error(errors.New("could not update business with shareholder and director info"), logger.LoggerOptions{
			Key:  "updated",
			Data: updated,
		})
		apperrors.AuthenticationError(ctx.Ctx, "Business verification faileds. Request a new token from the app", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	cache.Cache.CreateEntry(token.(string), true, time.Hour*24*10)
	cache.Cache.DeleteOne(fmt.Sprintf("%s-kyc-info-directors", auth_token_claims["businessID"].(string)))
	cache.Cache.DeleteOne(fmt.Sprintf("%s-kyc-info-shareholders", auth_token_claims["businessID"].(string)))
	cache.Cache.DeleteOne(fmt.Sprintf("%s-kyc-info-address", auth_token_claims["businessID"].(string)))
	pushnotification.PushNotificationService.PushOne(account.PushNotificationToken, "Your business has been verified!🥳", "That was easy, wasn't it? Now you're 1 step closer to unlimited transfer amounts.")
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "business verified", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}
