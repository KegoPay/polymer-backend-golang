package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/constants"
	"kego.com/application/controllers/dto"
	countriessupported "kego.com/application/countriesSupported"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	userusecases "kego.com/application/usecases/userUseCases"
	"kego.com/entities"
	"kego.com/infrastructure/database/repository/cache"
	identityverification "kego.com/infrastructure/identity_verification"
	"kego.com/infrastructure/logger"
	server_response "kego.com/infrastructure/serverResponse"
	"kego.com/infrastructure/validator"
)


func FetchUserProfile(ctx *interfaces.ApplicationContext[any]){
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if user == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("This user profile was not found. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL))
		return
	}
	walletRepo := repository.WalletRepo()
	wallet, err := walletRepo.FindByID(user.WalletID)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	signupCountries := countriessupported.FilterCountries(entities.SignUp)
	var country entities.Country
	for _, c := range signupCountries {
		if strings.Contains(c.Name, user.Nationality) {
			country = c
			country.ServicesAllowed = nil
			break
		}
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "profile fetched", map[string]any{
		"account": user,
		"wallet": wallet,
		"country": country,
	}, nil)
}

func UpdateUserProfile(ctx *interfaces.ApplicationContext[dto.UpdateUserDTO]){userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"))
	if err != nil {
		logger.Error(errors.New("error fetching user profile for update"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if ctx.Body.FirstName != nil {
		user.FirstName = *ctx.Body.FirstName
	}
	if ctx.Body.LastName != nil {
		user.LastName = *ctx.Body.LastName
	}
	if ctx.Body.Phone != nil {
		user.Phone = ctx.Body.Phone
	}
	validationErr := validator.ValidatorInstance.ValidateStruct(user)
	if validationErr != nil {
		apperrors.ValidationFailedError(ctx, validationErr)
		return
	}
	userRepo.UpdateByID(ctx.GetStringContextData("UserID"), user)
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "update completed", nil, nil)
}

func SetPaymentTag(ctx *interfaces.ApplicationContext[dto.SetPaymentTagDTO]){
	err := userusecases.UpdateUserTag(ctx.Ctx, ctx.GetStringContextData("UserID"), *ctx.Body)
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "Your payment tag has been set successfully", nil, nil)
}


func ToggleNotificationOptions(ctx *interfaces.ApplicationContext[dto.ToggleNotificationOptionsDTO]){
	userRepo := repository.UserRepo()
	affected, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]any{
		"notificationOptions": ctx.Body,
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if affected == 0 {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("Notification setting could not be updated because profile was not found. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL))
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "Notification setting updated", nil, nil)
}

func EmailSubscription(ctx *interfaces.ApplicationContext[dto.EmailSubscriptionDTO]) {
	emailSubRepo := repository.EmailSubRepo()
	exists, err := emailSubRepo.CountDocs(map[string]interface{}{
		"email": ctx.Body.Email,
		"channel": ctx.Body.Channel,
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if exists != 0 {
		server_response.Responder.Respond(ctx.Ctx, http.StatusOK,
			"Seems you have registered with this email previously.\nNot to worry, you still have access to exclusive insights, updates, and special offers delivered straight to your inbox. Thanks for staying connected with us! ", nil, nil)
		return
	}
	found := cache.Cache.FindOne(fmt.Sprintf("%s-email-blacklist", ctx.Body.Email))
	if found != nil {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("%s was not approved for signup on Polymer", ctx.Body.Email), nil)
		return 
	}
	valid, err := identityverification.IdentityVerifier.EmailVerification(ctx.Body.Email)
	if err != nil {
		apperrors.FatalServerError(ctx, err)
		return  
	}
	if !valid {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("%s was not approved for signup on Polymer", ctx.Body.Email), nil)
		cache.Cache.CreateEntry(fmt.Sprintf("%s-email-blacklist", ctx.Body.Email), ctx.Body.Email, time.Minute * 0 )
		return 
	}
	payload := entities.Subscriptions{
		Email: ctx.Body.Email,
		Channel: ctx.Body.Channel,
	}
	valiedationErr := validator.ValidatorInstance.ValidateStruct(payload)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr)
		return
	}
	emailSubRepo.CreateOne(context.TODO(), payload)
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "You're in! You now have access to exclusive insights, updates, and special offers delivered straight to your inbox.", nil, nil)
}


func UpdateAddress(ctx *interfaces.ApplicationContext[dto.UpdateAddressDTO]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr)
		return
	}
	userRepo := repository.UserRepo()
	updated, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]any{
		"address": entities.Address{
			FullAddress: nil,
			State: &ctx.Body.State,
			Street: &ctx.Body.Street,
			LGA: &ctx.Body.LGA,
			Verified: true,
		},
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if updated != 1 {
		apperrors.UnknownError(ctx.Ctx, errors.New("could not update users address"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "address set", nil, nil)
}


func VerifyCurrentAddress(ctx *interfaces.ApplicationContext[any]) {
	userRepo := repository.UserRepo()
	updated, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]any{
		"address": entities.Address{
			Verified: true,
		},
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if updated != 1 {
		apperrors.UnknownError(ctx.Ctx, errors.New("could not update users address"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "address verified", nil, nil)
}