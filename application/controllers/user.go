package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/constants"
	"usepolymer.co/application/controllers/dto"
	countriessupported "usepolymer.co/application/countriesSupported"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/repository"
	userusecases "usepolymer.co/application/usecases/userUseCases"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/auth"
	"usepolymer.co/infrastructure/biometric"
	"usepolymer.co/infrastructure/cryptography"
	"usepolymer.co/infrastructure/database/repository/cache"
	fileupload "usepolymer.co/infrastructure/file_upload"
	file_upload_types "usepolymer.co/infrastructure/file_upload/types"
	identityverification "usepolymer.co/infrastructure/identity_verification"
	"usepolymer.co/infrastructure/logger"
	sms "usepolymer.co/infrastructure/messaging/whatsapp"
	server_response "usepolymer.co/infrastructure/serverResponse"
	"usepolymer.co/infrastructure/validator"
)

func FetchUserProfile(ctx *interfaces.ApplicationContext[any]) {
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if user == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("This user profile was not found. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	walletRepo := repository.WalletRepo()
	wallet, err := walletRepo.FindByID(user.WalletID)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
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
		"wallet":  wallet,
		"country": country,
	}, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

// func UpdateUserProfile(ctx *interfaces.ApplicationContext[dto.UpdateUserDTO]){
// 	userRepo := repository.UserRepo()
// 	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"))
// 	if err != nil {
// 		logger.Error(errors.New("error fetching user profile for update"), logger.LoggerOptions{
// 			Key: "error",
// 			Data: err,
// 		})
// 		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
// 		return
// 	}
// 	if ctx.Body.FirstName != nil {
// 		user.FirstName = *ctx.Body.FirstName
// 	}
// 	if ctx.Body.LastName != nil {
// 		user.LastName = *ctx.Body.LastName
// 	}
// 	if ctx.Body.Phone != nil {
// 		user.Phone = ctx.Body.Phone
// 	}
// 	validationErr := validator.ValidatorInstance.ValidateStruct(user)
// 	if validationErr != nil {
// 		apperrors.ValidationFailedError(ctx, validationErr)
// 		return
// 	}
// 	userRepo.UpdateByID(ctx.GetStringContextData("UserID"), user)
// 	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "update completed", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
// }

func SetPaymentTag(ctx *interfaces.ApplicationContext[dto.SetPaymentTagDTO]) {
	err := userusecases.UpdateUserTag(ctx.Ctx, ctx.GetStringContextData("UserID"), *ctx.Body, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "Your payment tag has been set successfully", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func ToggleNotificationOptions(ctx *interfaces.ApplicationContext[dto.ToggleNotificationOptionsDTO]) {
	userRepo := repository.UserRepo()
	affected, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]any{
		"notificationOptions": ctx.Body,
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if affected == 0 {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("Notification setting could not be updated because profile was not found. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL), ctx.GetHeader("Polymer-Device-Id"))
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "Notification setting updated", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func EmailSubscription(ctx *interfaces.ApplicationContext[dto.EmailSubscriptionDTO]) {
	emailSubRepo := repository.EmailSubRepo()
	exists, err := emailSubRepo.CountDocs(map[string]interface{}{
		"email":   ctx.Body.Email,
		"channel": ctx.Body.Channel,
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if exists != 0 {
		server_response.Responder.Respond(ctx.Ctx, http.StatusOK,
			"Seems you have registered with this email previously.\nNot to worry, you still have access to exclusive insights, updates, and special offers delivered straight to your inbox. Thanks for staying connected with us! ", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	found := cache.Cache.FindOne(fmt.Sprintf("%s-email-blacklist", ctx.Body.Email))
	if found != nil {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("%s was flagged as a suspicious email and was not approved for signup on Polymer", ctx.Body.Email), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	valid, err := identityverification.IdentityVerifier.EmailVerification(ctx.Body.Email)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !valid {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("%s was not approved for signup on Polymer", ctx.Body.Email), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		cache.Cache.CreateEntry(fmt.Sprintf("%s-email-blacklist", ctx.Body.Email), ctx.Body.Email, time.Minute*0)
		return
	}
	payload := entities.Subscriptions{
		Email:   ctx.Body.Email,
		Channel: ctx.Body.Channel,
	}
	valiedationErr := validator.ValidatorInstance.ValidateStruct(payload)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	emailSubRepo.CreateOne(context.TODO(), payload)
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "You're in! You now have access to exclusive insights, updates, and special offers delivered straight to your inbox.", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func UpdateAddress(ctx *interfaces.ApplicationContext[dto.UpdateAddressDTO]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"phone":     1,
		"nextOfKin": 1,
		"bvn":       1,
		"nin":       1,
		"tier":      1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	var payload = map[string]any{
		"address": entities.Address{
			FullAddress: nil,
			State:       &ctx.Body.State,
			Street:      &ctx.Body.Street,
			LGA:         &ctx.Body.LGA,
			Verified:    true,
		},
	}

	if ctx.Body.AuthOne {
		if user.Phone.IsVerified && user.NextOfKin != nil && user.BVN != "" && user.NIN != "" {
			if user.Tier == 2 {
				payload["tier"] = 3
			}
		} else if (user.Phone.IsVerified && user.NextOfKin != nil) && (user.BVN == "" || user.NIN == "") {
			if user.Tier == 1 {
				payload["tier"] = 2
			}
		}
	}

	updated, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), payload)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if updated != 1 {
		apperrors.UnknownError(ctx.Ctx, errors.New("could not update users address"), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "address set", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func LinkNIN(ctx *interfaces.ApplicationContext[dto.LinkNINDTO]) {
	attemptsLeft := cache.Cache.FindOne(fmt.Sprintf("%s-nin-kyc-attempts-left", ctx.GetStringContextData("Email")))
	if attemptsLeft == nil {
		apperrors.ClientError(ctx.Ctx, `You cannot link your NIN to this account at this point, most likely because it has already been done before. If you think this is a mistake and want a manual review please click the button below`, nil, &constants.NIN_VERIFICATION_FAILED, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	parsedAttemptsLeft, err := strconv.Atoi(*attemptsLeft)
	if err != nil {
		apperrors.ClientError(ctx.Ctx, `You’ve reached the maximum number of tries allowed for this. If you think this is a mistake and want a manual review please click the button below`, nil, &constants.NIN_VERIFICATION_FAILED, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if parsedAttemptsLeft == 0 {
		apperrors.ClientError(ctx.Ctx, `You’ve reached the maximum number of tries allowed for this. If you think this is a mistake and want a manual review please click the button below`, nil, &constants.NIN_VERIFICATION_FAILED, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"profileImage": 1,
		"nextOfKin":    1,
		"tier":         1,
		"nin":          1,
		"bvn":          1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account.NIN != "" {
		apperrors.ClientError(ctx.Ctx, "You already have an NIN linked to your profile. If you do not remember doing so we probably got it from your BVN details 😉", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account.BVN == "" {
		apperrors.ClientError(ctx.Ctx, "This endpoint should be used only after a BVN has been used to verify the account", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	ninDetails, err := identityverification.IdentityVerifier.FetchNINDetails(ctx.Body.NIN)
	if err != nil {
		cache.Cache.CreateEntry(fmt.Sprintf("%s-nin-kyc-attempts-left", ctx.GetStringContextData("Email")), parsedAttemptsLeft-1, time.Hour*24*365)
		apperrors.CustomError(ctx.Ctx, err.Error(), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	url, err := fileupload.FileUploader.GeneratedSignedURL(account.ProfileImage, file_upload_types.SignedURLPermission{
		Read: true,
	})
	if err != nil {
		logger.Error(errors.New("signed url could not be generated for file download"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return
	}
	if url == nil {
		logger.Error(errors.New("signed url could not be generated for file download. nil url returned and no error"))
		return
	}
	result, err := biometric.BiometricService.FaceMatch(&ninDetails.Base64Image, url)
	if err != nil {
		apperrors.ClientError(ctx.Ctx, "something went wrong while performing face verification", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !result {
		cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft-1, time.Hour*24*365) // keep data cached for a yea
		apperrors.ClientError(ctx.Ctx, "We compared your face with that on your ID and it did not match. Please ensure you are in a well lit environment and have no coverings on your face.", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	encryptedNIN, err := cryptography.SymmetricEncryption(ctx.Body.NIN, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userUpdatedInfo := map[string]any{
		"nin":       encryptedNIN,
		"ninLinked": true,
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "this account no longer exists", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !ctx.Body.AuthOne {
		if account.NextOfKin != nil {
			if account.Tier == 2 {
				userUpdatedInfo["tier"] = 3
			}
		} else {
			apperrors.ClientError(ctx.Ctx, "could not upgrade your account", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
	} else {
		if account.NextOfKin != nil && account.Phone.IsVerified && account.Address.Verified {
			if account.NextOfKin != nil {
				if account.Tier == 2 {
					userUpdatedInfo["tier"] = 3
				}
			} else {
				apperrors.ClientError(ctx.Ctx, "could not upgrade your account", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
				return
			}
		} else {
			if account.Tier == 2 {
				userUpdatedInfo["tier"] = 3
			}
		}
	}
	success, err := userRepo.UpdatePartialByFilter(map[string]interface{}{
		"_id": ctx.GetStringContextData("UserID"),
	}, userUpdatedInfo)
	if err != nil {
		logger.Error(errors.New("an error occured while updating nin and tier"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !success {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	cache.Cache.DeleteOne(fmt.Sprintf("%s-nin-kyc-attempts-left", ctx.GetStringContextData("Email")))
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "NIN verified", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func UpdatePhone(ctx *interfaces.ApplicationContext[dto.UpdatePhoneDTO]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.CountDocs(map[string]interface{}{
		"phone.localNumber": ctx.Body.Phone,
	})
	if account != 0 {
		apperrors.EntityAlreadyExistsError(ctx.Ctx, "this number is already linked to another account", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	updated, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]any{
		"phone": entities.PhoneNumber{
			WhatsApp:    ctx.Body.WhatsApp,
			LocalNumber: ctx.Body.Phone,
			Prefix:      "234",
			ISOCode:     "NG",
			IsVerified:  false,
			Modified:    true,
		},
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if updated != 1 {
		apperrors.UnknownError(ctx.Ctx, errors.New("could not update users phone"), ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	var otp *string
	if ctx.Body.WhatsApp {
		otp, err = auth.GenerateOTP(6, ctx.Body.Phone)
		if err != nil {
			apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
	}
	ref := sms.SMSService.SendOTP(fmt.Sprintf("%s%s", "234", ctx.Body.Phone), ctx.Body.WhatsApp, otp)
	encryptedRef, err := cryptography.SymmetricEncryption(*ref, nil)
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	cache.Cache.CreateEntry(fmt.Sprintf("%s-sms-otp-ref", ctx.Body.Phone), encryptedRef, time.Minute*10)
	cache.Cache.CreateEntry(fmt.Sprintf("%s-otp-intent", ctx.Body.Phone), "verify_phone", time.Minute*10)
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "phone set", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func VerifyCurrentPhone(ctx *interfaces.ApplicationContext[dto.IsAuthOne]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"phone":     1,
		"nextOfKin": 1,
		"bvn":       1,
		"nin":       1,
		"address":   1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "this account does not exist", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account.Phone.Modified {
		apperrors.ClientError(ctx.Ctx, "current phone number cannot be verified because it has been modfied", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	account.Phone.IsVerified = true
	if account.Phone != nil {
		if account.Phone.LocalNumber != "" {
			account.Phone.IsVerified = true
		}
	} else {
		apperrors.ClientError(ctx.Ctx, "phone number not set", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	payload := map[string]any{
		"phone": account.Phone,
	}

	if ctx.Body.AuthOne {
		if account.Address.Verified && account.NextOfKin != nil && account.BVN != "" && account.NIN != "" {
			if account.Tier == 2 {
				payload["tier"] = 3
			}
		} else if (account.Address.Verified && account.NextOfKin != nil) && (account.BVN == "" || account.NIN == "") {
			if account.Tier == 1 {
				payload["tier"] = 2
			}
		}
	}

	success, err := userRepo.UpdatePartialByID(account.ID, payload)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if success != 1 {
		err = errors.New("could not modify account phone details")
		logger.Error(err, logger.LoggerOptions{
			Key:  "id",
			Data: account.ID,
		}, logger.LoggerOptions{
			Key:  "phone",
			Data: account.Phone,
		})
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "phone verified", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func VerifyCurrentAddress(ctx *interfaces.ApplicationContext[dto.IsAuthOne]) {
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"phone":     1,
		"nextOfKin": 1,
		"bvn":       1,
		"nin":       1,
		"tier":      1,
		"address":   1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	var payload = map[string]any{
		"address.verified": true,
	}
	if user.Address.FullAddress == nil || *user.Address.FullAddress == "" {
		apperrors.ClientError(ctx.Ctx, "current address is not set. manually set your address", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	if ctx.Body.AuthOne {
		if user.Phone.IsVerified && user.NextOfKin != nil && user.BVN != "" && user.NIN != "" {
			if user.Tier == 2 {
				payload["tier"] = 3
			}
		} else if (user.Phone.IsVerified && user.NextOfKin != nil) && (user.BVN == "" || user.NIN == "") {
			if user.Tier == 1 {
				payload["tier"] = 2
			}
		}
	}
	updated, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), payload)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if updated != 1 {
		apperrors.UnknownError(ctx.Ctx, errors.New("could not update users address"), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "address verified", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func SetNextOfKin(ctx *interfaces.ApplicationContext[dto.SetNextOfKin]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userRepo := repository.UserRepo()
	user, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"phone":     1,
		"address":   1,
		"nextOfKin": 1,
		"bvn":       1,
		"nin":       1,
		"tier":      1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	payload := map[string]any{
		"nextOfKin": entities.NextOfKin{
			FirstName:    ctx.Body.FirstName,
			LastName:     ctx.Body.LastName,
			Relationship: ctx.Body.Relationship,
		},
	}
	if !ctx.Body.AuthOne {
		if user.Phone.IsVerified && user.Address.Verified {
			if user.Tier == 1 {
				payload["tier"] = 2
			}
		} else {
			apperrors.ClientError(ctx.Ctx, "could not upgrade your account", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
	} else {
		if user.Phone.IsVerified && user.Address.Verified && user.BVN != "" && user.NIN != "" {
			if user.Tier == 2 {
				payload["tier"] = 3
			}
		} else if (user.Phone.IsVerified && user.Address.Verified) && (user.BVN == "" || user.NIN == "") {
			if user.Tier == 1 {
				payload["tier"] = 2
			}
		}
	}

	updated, err := userRepo.UpdatePartialByFilter(map[string]interface{}{
		"_id": ctx.GetStringContextData("UserID"),
	}, payload, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while updating next of kin and tier"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !updated {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "next of kin updated", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func GenerateFileURL(ctx *interfaces.ApplicationContext[dto.FileUploadOptions]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	fileName := fmt.Sprintf("%s/%s", ctx.GetStringContextData("UserID"), ctx.Body.Type)
	url, err := fileupload.FileUploader.GeneratedSignedURL(fileName, ctx.Body.Permissions)
	if err != nil {
		apperrors.CustomError(ctx.Ctx, err.Error(), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "url geenraed", map[string]string{
		"url":      *url,
		"fileName": fileName,
	}, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}
