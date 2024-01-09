package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/constants"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/application/services"
	"kego.com/application/services/types"
	authusecases "kego.com/application/usecases/authUsecases"
	"kego.com/application/utils"
	"kego.com/entities"
	"kego.com/infrastructure/auth"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/database/repository/cache"
	fileupload "kego.com/infrastructure/file_upload"
	identityverification "kego.com/infrastructure/identity_verification"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
	server_response "kego.com/infrastructure/serverResponse"
	"kego.com/infrastructure/validator"
)

func CreateAccount(ctx *interfaces.ApplicationContext[dto.CreateAccountDTO]) {
	account, _, err := authusecases.CreateAccount(ctx.Ctx, &entities.User{
		Email:          ctx.Body.Email,
		Password: 		ctx.Body.Password,
		TransactionPin: ctx.Body.TransactionPin,
		UserAgent:      ctx.Body.UserAgent,
		DeviceID:       ctx.Body.DeviceID,
		BVN: 			ctx.Body.BVN,
		AppVersion: ctx.Body.AppVersion,
	})
	if err != nil {
		return
	}
	otp, err := auth.GenerateOTP(6, account.Email)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), 2, time.Hour * 24 * 365 ) // keep data cached for a year
	emails.EmailService.SendEmail(account.Email, "Welcome to Kego! Verify your account to continue", "otp", map[string]interface{}{
		"FIRSTNAME": account.FirstName,
		"OTP":      otp,
	},)
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "account created", nil, nil)
}

func LoginUser(ctx *interfaces.ApplicationContext[dto.LoginDTO]){
	appVersion := utils.ExtractAppVersionFromUserAgentHeader(ctx.GetHeader("User-Agent").(string))
	if appVersion == nil {
		apperrors.UnsupportedAppVersion(ctx.Ctx)
		return
	}
	account, token := authusecases.LoginAccount(ctx.Ctx, ctx.Body.Email, ctx.Body.Phone, &ctx.Body.Password, *appVersion, ctx.GetHeader("User-Agent").(string), ctx.Body.DeviceID)
	if account == nil || token == nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "login successful", map[string]interface{}{
		"account": account,
		"token":   token,
	}, nil)
}


func ResetPassword(ctx *interfaces.ApplicationContext[dto.ResetPasswordDTO]) {
	msg, success := auth.VerifyOTP(ctx.Body.Email, ctx.Body.Otp)
	if !success {
		apperrors.ClientError(ctx.Ctx, msg, nil)
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.Body.Email,
	})
	if err != nil {
		logger.Error(errors.New("error fetching a user account to reset password"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "account with email not found")
		return
	}
	account.Password = ctx.Body.NewPassword
	valiedationErr := validator.ValidatorInstance.ValidateStruct(*account)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr)
		return
	}
	hashedPassword, err := cryptography.CryptoHahser.HashString(ctx.Body.NewPassword)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	success, err = userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email": ctx.Body.Email,
	}, map[string]interface{}{
		"password": string(hashedPassword),
	})
	if !success || err != nil {
		apperrors.FatalServerError(ctx.Ctx)
	}
	
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "password reset", nil, nil)
}

func UpdatePassword(ctx *interfaces.ApplicationContext[dto.UpdatePassword]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("Email"),
	})
	if err != nil {
		logger.Error(errors.New("error fetching a user account to reset password"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "account with email not found")
		return
	}
	success := cryptography.CryptoHahser.VerifyData(account.Password, ctx.Body.CurrentPassword)
	if !success {
		apperrors.ClientError(ctx.Ctx, "incorrect password", nil)
		return
	}
	account.Password = ctx.Body.NewPassword
	valiedationErr := validator.ValidatorInstance.ValidateStruct(*account)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr)
		return
	}
	hashed_password, err := cryptography.CryptoHahser.HashString(ctx.Body.NewPassword)
	if err != nil {
		logger.Error(errors.New("error hashing users new password"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	modified, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]interface{}{
		"password": string(hashed_password),
	})
	if !success || err != nil {
		logger.Error(errors.New("error while updating user password"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		},  logger.LoggerOptions{
			Key: "modified",
			Data: modified,
		}, )
		apperrors.FatalServerError(ctx.Ctx)
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "password updated", nil, nil)
}

func ResendOTP(ctx *interfaces.ApplicationContext[any]) {
	email := ctx.Query["email"].(string)
	if email == "" {
		server_response.Responder.Respond(ctx.Ctx, http.StatusBadRequest, "pass in a valid email to recieve the otp", nil, nil)
		return
	}
	otp, err := auth.GenerateOTP(6, email)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": email,
	}, options.FindOne().SetProjection(map[string]any{
		"firstName": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if account == nil {
		server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "otp sent", nil, nil)
		return
	}
	emails.EmailService.SendEmail(email, "An OTP was requested for your account", "otp", map[string]interface{}{
		"FIRSTNAME": account.FirstName,
		"OTP":      otp,
	},)
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "otp sent", nil, nil)
}

func VerifyEmail(ctx *interfaces.ApplicationContext[dto.VerifyEmailData]) {
	msg, success := auth.VerifyOTP(ctx.Body.Email, ctx.Body.Otp)
	if !success {
		apperrors.ClientError(ctx.Ctx, msg, nil)
		return
	}
	userRepo := repository.UserRepo()
	userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email": ctx.Body.Email,
	}, map[string]bool{
		"emailVerified": true,
	})
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "account verified", nil, nil)
}

func VerifyAccount(ctx *interfaces.ApplicationContext[dto.VerifyAccountData]){
	attemptsLeft := cache.Cache.FindOne(fmt.Sprintf("%s-kyc-attempts-left", ctx.Body.Email))
	if attemptsLeft == nil {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf( "You cannot perform kyc at this moment. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL), nil)
		return
	}
	parsedAttemptsLeft, err := strconv.Atoi(*attemptsLeft)
	if err != nil {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf( "You cannot perform kyc at this moment. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL), nil)
		return
	}
	if parsedAttemptsLeft == 0 {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf( "You cannot perform kyc at this moment. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL), nil)
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.Body.Email,
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("Account with email %s does not exist. Please contact support on %s to help resolve this issue.", ctx.Body.Email, constants.SUPPORT_EMAIL))
		return
	}
	if !account.EmailVerified {
		apperrors.ClientError(ctx.Ctx, "verify your email before attempting identity verification", nil)
		return
	}

	if account.KYCCompleted {
		apperrors.ClientError(ctx.Ctx, "you have completed your identity verification", nil)
		return
	}
	bvnDetails, err := identityverification.IdentityVerifier.FetchBVNDetails(account.BVN)
	if err != nil {
		cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft - 1 , time.Hour * 24 * 365 ) // keep data cached for a year
		apperrors.CustomError(ctx.Ctx, err.Error())
		return
	}
	url, err := fileupload.FileUploader.UploadSingleFile(ctx.Body.ProfileImage, &account.ID)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	result, err := identityverification.IdentityVerifier.FaceMatch(*url, bvnDetails.Base64Image)
	if err != nil {
		cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft - 1 , time.Hour * 24 * 365 ) // keep data cached for a year
		cldErr := fileupload.FileUploader.DeleteSingleFile(account.ID)
		if cldErr != nil {
			apperrors.FatalServerError(ctx.Ctx)
			return
		}
		apperrors.ClientError(ctx.Ctx, err.Error(), nil)
		return
	}
	if *result < 80 {
		cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft - 1 , time.Hour * 24 * 365 ) // keep data cached for a year
		err = fileupload.FileUploader.DeleteSingleFile(account.ID)
		if err != nil {
			apperrors.FatalServerError(ctx.Ctx)
			return
		}
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("Your picture does not match with your Image on the BVN provided. If you think this is a mistake please contact support on %s", constants.SUPPORT_EMAIL), nil)
		return
	}
	userUpdatedInfo := map[string]any{
		"gender": bvnDetails.Gender,
		"dob": bvnDetails.DateOfBirth,
		"lastName": bvnDetails.LastName,
		"firstName": bvnDetails.FirstName,
		"middleName": bvnDetails.MiddleName,
		"watchListed": bvnDetails.WatchListed,
		"nationality": bvnDetails.Nationality,
		"phone": entities.PhoneNumber{
			Prefix: "234",
			ISOCode: "NG",
			LocalNumber: bvnDetails.PhoneNumber,
		},
		"profileImage": *url,
		"kycCompleted": true,
	}
	userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email": ctx.Body.Email,
	}, userUpdatedInfo)
	cache.Cache.DeleteOne(fmt.Sprintf("%s-kyc-attempts-left", account.Email))
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "kyc completed", nil, nil)
}

func AccountWithEmailExists(ctx *interfaces.ApplicationContext[any]){
	email := ctx.Query["email"]
	if email == "" {
		server_response.Responder.Respond(ctx.Ctx, http.StatusBadRequest, "pass in a valid email", nil, nil)
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": email,
	}, options.FindOne().SetProjection(map[string]interface{}{
		"emailVerified": 1,
		"kycCompleted": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	response := map[string]any{}
	if account == nil {
		response["exists"] = false
	}else {
		response["exists"] = true
		response["emailVerified"] = account.EmailVerified
		response["KYCCompleted"] = account.KYCCompleted
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "success", response, nil)
}

func DeactivateAccount(ctx *interfaces.ApplicationContext[dto.ConfirmPin]){
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"deactivated": 1,
		"password": 1,
		"transactionPin": 1,
	}))
	if err != nil {
		logger.Error(errors.New("error fetching a user account to deactivate account"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("This user profile was not found. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL))
		return
	} 
	if account.Deactivated {
		apperrors.ClientError(ctx.Ctx, "account has already been deactivated", nil)
		return
	}
	match := services.VerifyPin(ctx.Ctx, account, ctx.Body.Pin, &types.PinSelectionType{
		Password: true,
	})
	if !match {
		return
	}
	success, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]interface{}{
		"deactivated": true,
	})
	if success == 0 || err != nil {
		logger.Error(errors.New("error while deactivating user account"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		},  logger.LoggerOptions{
			Key: "success",
			Data: success,
		}, )
		apperrors.FatalServerError(ctx.Ctx)
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "deactivated", nil, nil)
}