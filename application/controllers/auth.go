package controllers

import (
	"crypto/ecdh"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/constants"
	"kego.com/application/controllers/dto"
	countriessupported "kego.com/application/countriesSupported"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/application/services"
	"kego.com/application/services/types"
	authusecases "kego.com/application/usecases/authUsecases"
	"kego.com/application/usecases/wallet"
	"kego.com/application/utils"
	"kego.com/entities"
	"kego.com/infrastructure/auth"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/database/repository/cache"
	fileupload "kego.com/infrastructure/file_upload"
	identityverification "kego.com/infrastructure/identity_verification"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
	pushnotification "kego.com/infrastructure/messaging/push_notifications"
	sms "kego.com/infrastructure/messaging/whatsapp"

	server_response "kego.com/infrastructure/serverResponse"
	"kego.com/infrastructure/validator"
)

func KeyExchange(ctx *interfaces.ApplicationContext[dto.GenerateServerPublicKey]) {
	keyBytes, err := hex.DecodeString(ctx.Body.ClientPubKey)
    if err != nil {
		logger.Error(errors.New("error decoding keys for key exchange"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.UnknownError(ctx.Ctx, errors.New("could not perform key exchange"))
        return
    }
    ClientPubKey, err := ecdh.P256().NewPublicKey(keyBytes)
    if err != nil {
		logger.Error(errors.New("error geting public key from key bytes"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.UnknownError(ctx.Ctx, errors.New("could not perform key exchange"))
        return
    }
	serverPubKey := cryptography.GeneratePublicKey(ctx.Body.SessionID, ClientPubKey)
	server_response.Responder.UnEncryptedRespond(ctx.Ctx, http.StatusCreated, "key generated", serverPubKey, nil)
}

func VerifyOTP(ctx *interfaces.ApplicationContext[dto.VerifyOTPDTO]) {
	if ctx.Body.Phone == nil && ctx.Body.Email == nil {
		apperrors.ClientError(ctx.Ctx, "pass in either a phone number or email", nil)
		return
	}
	var channel = ""
	var filter = map[string]any{}
	if ctx.Body.Email != nil {
		channel = *ctx.Body.Email
		filter["email"] = channel
	}else {
		channel = *ctx.Body.Phone
		filter["phone.localNumber"] = channel
	}
	msg, success := auth.VerifyOTP(*ctx.Body.Email, ctx.Body.OTP)
	if !success {
		apperrors.ClientError(ctx.Ctx, msg, nil)
		return
	}
	otpIntent := cache.Cache.FindOne(fmt.Sprintf("%s-otp-intent", channel))
	if otpIntent == nil {
		logger.Error(errors.New("otp intent missing"))
		apperrors.ClientError(ctx.Ctx, "otp expired", nil)
		return
	}
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email: ctx.Body.Email,
		PhoneNum: ctx.Body.Phone,
		OTPIntent: *otpIntent,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(15)).Unix(), //lasts for 10 mins
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}

	server_response.Responder.UnEncryptedRespond(ctx.Ctx, http.StatusCreated, "otp verified", token, nil)
}

func CreateAccount(ctx *interfaces.ApplicationContext[dto.CreateAccountDTO]) {
	account, _, err := authusecases.CreateAccount(ctx.Ctx, &entities.User{
		Email:          ctx.Body.Email,
		Password: 		ctx.Body.Password,
		UserAgent:      ctx.Body.UserAgent,
		DeviceID:       ctx.Body.DeviceID,
		AppVersion: 	ctx.Body.AppVersion,
		PushNotificationToken: ctx.Body.PushNotificationToken,
	})
	if err != nil {
		return
	}
	otp, err := auth.GenerateOTP(6, account.Email)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), 2, time.Hour * 24 * 365 ) // keep data cached for a year
		emails.EmailService.SendEmail(account.Email, "Welcome to Kego! Verify your account to continue", "otp", map[string]interface{}{
			"FIRSTNAME": account.FirstName,
			"OTP":      otp,
		},)
	cache.Cache.CreateEntry(fmt.Sprintf("%s-otp-intent", ctx.Body.Email), "verify_account", time.Minute * 10)
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "account created", nil, nil)
}

func LoginUser(ctx *interfaces.ApplicationContext[dto.LoginDTO]){
	appVersion := utils.ExtractAppVersionFromUserAgentHeader(ctx.GetHeader("User-Agent").(string))
	if appVersion == nil {
		apperrors.UnsupportedAppVersion(ctx.Ctx)
		return
	}
	account, wallet, token := authusecases.LoginAccount(ctx.Ctx, ctx.Body.Email, ctx.Body.Phone, &ctx.Body.Password, *appVersion, ctx.GetHeader("User-Agent").(string), ctx.Body.DeviceID, ctx.Body.PushNotificationToken)
	if account == nil || token == nil {
		return
	}
	signupCountries := countriessupported.FilterCountries(entities.SignUp)
	var country entities.Country
	for _, c := range signupCountries {
		if strings.Contains(c.Name, account.Nationality) {
			country = c
			country.ServicesAllowed = nil
			break
		}
	}
	responsePayload := map[string]interface{}{
		"account": account,
		"wallet": wallet,
		"token":   token,
		"country": country,
	}
	if account.TransactionPin == "" {
		responsePayload["unsetTrxPin"] = true
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "login successful", responsePayload, nil)
}

func ResetPassword(ctx *interfaces.ApplicationContext[dto.ResetPasswordDTO]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email":  ctx.GetStringContextData("OTPEmail"),
	})
	if err != nil {
		logger.Error(errors.New("error fetching a user account to reset password"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err)
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
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	success, err := userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email":  ctx.GetStringContextData("OTPEmail"),
	}, map[string]interface{}{
		"password": string(hashedPassword),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	
	if !success {
		logger.Error(errors.New("could not reset password"), logger.LoggerOptions{
			Key: "email",
			Data: ctx.GetStringContextData("OTPEmail"),
		})
		apperrors.UnknownError(ctx.Ctx, err)
		return
	}
	cache.Cache.CreateEntry(ctx.GetStringContextData("OTPToken"), true, time.Minute * 5)
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
		apperrors.FatalServerError(ctx.Ctx, err)
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
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	modified, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]interface{}{
		"password": string(hashed_password),
	})
	if err != nil {
		logger.Error(errors.New("error while updating user password"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err)
	}
	if modified == 0 {
		logger.Error(errors.New("error while updating user password"),  logger.LoggerOptions{
			Key: "modified",
			Data: modified,
		}, )
		apperrors.FatalServerError(ctx.Ctx, fmt.Errorf("failed to update users password userID %s", ctx.GetStringContextData("UserID")))
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "password updated", nil, nil)
}

func ResendOTP(ctx *interfaces.ApplicationContext[dto.ResendOTP]) {
	if ctx.Body.Phone == nil && ctx.Body.Email == nil {
		apperrors.ClientError(ctx.Ctx, "pass in either a phone number or email", nil)
		return
	}
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr)
		return
	}
	var channel = ""
	var filter = map[string]any{}
	if ctx.Body.Email != nil {
		channel = *ctx.Body.Email
		filter["email"] = channel
	}else {
		channel = *ctx.Body.Phone
		filter["phone.localNumber"] = channel
	}
	
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(filter, options.FindOne().SetProjection(map[string]any{
		"firstName": 1,
		"phone": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if account == nil {
		server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "otp sent", nil, nil)
		return
	}

	if ctx.Body.Email != nil {
		otp, err := auth.GenerateOTP(6, channel)
		if err != nil {
			apperrors.FatalServerError(ctx.Ctx, err)
		}
		emails.EmailService.SendEmail(channel, "An OTP was requested for your account", "otp", map[string]interface{}{
			"FIRSTNAME": account.FirstName,
			"OTP":      otp,
		},)
	}else if ctx.Body.Phone != nil {
		ref := sms.SMSService.SendSMS(fmt.Sprintf("%s%s", account.Phone.Prefix, account.Phone.LocalNumber), "< 123456 >is your Polymer OTP")
		encryptedRef, err := cryptography.SymmetricEncryption(*ref)
		if err != nil {
			apperrors.UnknownError(ctx.Ctx, err)
			return
		}
		cache.Cache.CreateEntry(fmt.Sprintf("%s-sms-otp-ref", channel), *encryptedRef, time.Minute * 10)
	}
	cache.Cache.CreateEntry(fmt.Sprintf("%s-otp-intent", channel), ctx.Body.Intent, time.Minute * 10)
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "otp sent", nil, nil)
}

func VerifyEmail(ctx *interfaces.ApplicationContext[any]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("OTPEmail"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "this account no longer exists")
		return
	}
	if account.EmailVerified {
		apperrors.ClientError(ctx.Ctx, "this email has already been verified", nil)
		return
	}
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email:     &account.Email,
		Phone:     account.Phone,
		UserID:    account.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(15)).Unix(), //lasts for 15 mins
		UserAgent: account.UserAgent,
		FirstName: account.FirstName,
		LastName: account.LastName,
		DeviceID:   account.DeviceID,
		AppVersion: account.AppVersion,
		PushNotificationToken: account.PushNotificationToken,
	})
	account.EmailVerified = true
	success, err := userRepo.UpdateByID(account.ID, account)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if !success {
		apperrors.UnknownError(ctx.Ctx, err)
		return
	}
	cache.Cache.CreateEntry(ctx.GetStringContextData("OTPToken"), true, time.Minute * 5)
	hashedToken, err := cryptography.CryptoHahser.HashString(*token)
	if err != nil {
		apperrors.FatalServerError(ctx, err)
		return
	}
	cache.Cache.CreateEntry(account.ID, hashedToken, time.Minute * time.Duration(10)) // cache authentication token for 10 mins
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "account verified", map[string]string{
		"token": *token,
	}, nil)
}

func VerifyAccount(ctx *interfaces.ApplicationContext[dto.VerifyAccountData]){
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr)
		return
	}
	attemptsLeft := cache.Cache.FindOne(fmt.Sprintf("%s-kyc-attempts-left", ctx.GetStringContextData("Email")))
	if attemptsLeft == nil {
		apperrors.ClientError(ctx.Ctx, `You’ve reach the maximum number of tries allowed for this.`, nil)
		return
	}
	parsedAttemptsLeft, err := strconv.Atoi(*attemptsLeft)
	if err != nil {
		apperrors.ClientError(ctx.Ctx, `You’ve reach the maximum number of tries allowed for this.`, nil)
		return
	}
	if parsedAttemptsLeft == 0 {
		apperrors.ClientError(ctx.Ctx, `You’ve reach the maximum number of tries allowed for this.`, nil)
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("Email"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("Account with email %s does not exist. Please contact support on %s to help resolve this issue.", ctx.GetStringContextData("Email"), constants.SUPPORT_EMAIL))
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
	kycDetails := struct{
		Gender            string
		WatchListed       *string
		FirstName         string
		MiddleName        *string 
		LastName          string
		DateOfBirth       string 
		PhoneNumber       *string
		Nationality       string
		Base64Image       string 
		Address		      string 
	}{}
	if ctx.Body.Path == "bvn" {
		bvnDetails, err := identityverification.IdentityVerifier.FetchBVNDetails(*ctx.Body.BVN)
		if err != nil {
			cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft - 1 , time.Hour * 24 * 365 ) // keep data cached for a year
			apperrors.CustomError(ctx.Ctx, err.Error())
			return
		}
		kycDetails.Base64Image = bvnDetails.Base64Image
		kycDetails.WatchListed = &bvnDetails.WatchListed
		kycDetails.FirstName = bvnDetails.FirstName
		kycDetails.MiddleName = bvnDetails.MiddleName
		kycDetails.LastName = bvnDetails.LastName
		kycDetails.Gender = bvnDetails.Gender
		kycDetails.PhoneNumber = &bvnDetails.PhoneNumber
		kycDetails.Nationality = bvnDetails.Nationality
		kycDetails.DateOfBirth = bvnDetails.DateOfBirth
		kycDetails.Address = bvnDetails.Address
	}else if ctx.Body.Path == "nin" {
		apperrors.ClientError(ctx.Ctx, "Verification by NIN is currently not supported", nil)
		return
		ninDetails, err := identityverification.IdentityVerifier.FetchNINDetails(*ctx.Body.NIN)
		if err != nil {
			cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft - 1 , time.Hour * 24 * 365 ) // keep data cached for a year
			apperrors.CustomError(ctx.Ctx, err.Error())
			return
		}
		kycDetails.Base64Image = ninDetails.Base64Image
		kycDetails.WatchListed = nil
		kycDetails.FirstName = ninDetails.FirstName
		kycDetails.MiddleName = ninDetails.MiddleName
		kycDetails.LastName = ninDetails.LastName
		kycDetails.Gender = ninDetails.Gender
		kycDetails.PhoneNumber = ninDetails.PhoneNumber
		kycDetails.Nationality = ninDetails.Nationality
		kycDetails.DateOfBirth = ninDetails.DateOfBirth
	}
	// result, err := identityverification.IdentityVerifier.FaceMatch(*&ctx.Body.ProfileImage, bvnDetails.Base64Image)
	if err != nil {
		cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft - 1 , time.Hour * 24 * 365 ) // keep data cached for a year
		// _, _ := fileupload.FileUploader.DeleteFileByURL(ctx.Body.ProfileImage)
		// if cldErr != nil {
		// 	apperrors.FatalServerError(ctx.Ctx, cldErr)
		// 	return
		// }
		apperrors.ClientError(ctx.Ctx, err.Error(), nil)
		return
	}
	// if *result < 80 {
		// cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft - 1 , time.Hour * 24 * 365 ) // keep data cached for a year
		// err = fileupload.FileUploader.DeleteSingleFile(account.ID)
		// if err != nil {
		// 	apperrors.FatalServerError(ctx.Ctx, err)
		// 	return
		// }
	// 	apperrors.ClientError(ctx.Ctx, fmt.Sprintf("Your picture does not match with your Image on the BVN provided. If you think this is a mistake please contact support on %s", constants.SUPPORT_EMAIL), nil)
	// 	return
	// }
	watchListed := false
	if  kycDetails.WatchListed  != nil {
		if *kycDetails.WatchListed == "True" {
			watchListed = true
		}
	}
	userUpdatedInfo := map[string]any{
		"gender": kycDetails.Gender,
		"dob": kycDetails.DateOfBirth,
		"lastName": cases.Title(language.Und).String(kycDetails.LastName),
		"firstName": cases.Title(language.Und).String(kycDetails.FirstName),
		"middleName": func () *string {
			if kycDetails.MiddleName != nil {
				return utils.GetStringPointer(cases.Title(language.Und).String(*kycDetails.MiddleName))
			}
			return nil
		}(),
		"watchListed": watchListed,
		"nationality": kycDetails.Nationality,
		"phone": func () *entities.PhoneNumber {
			if kycDetails.PhoneNumber != nil {
				return &entities.PhoneNumber{
					Prefix: "234",
					ISOCode: "NG",
					LocalNumber: *kycDetails.PhoneNumber,
				}
			}	
			return nil
		}(),
		"profileImage": ctx.Body.ProfileImage,
		"kycCompleted": true,
		"bvn": ctx.Body.BVN,
		"nin": ctx.Body.NIN,
		"accountRestricted": watchListed,
		"address": entities.Address{
			FullAddress: &kycDetails.Address,
		},
	}
	userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("Email"),
	}, userUpdatedInfo)
	if ctx.Body.Path == "bvn" {
	wallet.GenerateNGNDVA(ctx.Ctx, account.WalletID,  account.FirstName, account.LastName, account.Email, *ctx.Body.BVN)
	}
	cache.Cache.DeleteOne(fmt.Sprintf("%s-kyc-attempts-left", account.Email))
	now := time.Now()
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email:     &account.Email,
		Phone:     &entities.PhoneNumber{
			Prefix: "234",
			ISOCode: "NG",
			LocalNumber: *kycDetails.PhoneNumber,
		},
		UserID:    account.ID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Local().Add(time.Minute * time.Duration(15)).Unix(), //lasts for 10 mins
		UserAgent: account.UserAgent,
		FirstName: userUpdatedInfo["firstName"].(string),
		LastName: userUpdatedInfo["lastName"].(string),
		DeviceID:   account.DeviceID,
		PushNotificationToken: account.PushNotificationToken,
		AppVersion: account.AppVersion,
	})
	pushnotification.PushNotificationService.PushOne(account.PushNotificationToken, "Welcome to Polymer!😃",
	"You now have global payments at your finger tips! Make payments with crypto, Mobile Money and to bank accounts in over 40+ countries!🤯")
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "kyc completed", token, nil)
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
		apperrors.FatalServerError(ctx.Ctx, err)
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

func GenerateFileURL(ctx *interfaces.ApplicationContext[dto.FileUploadOptions]){
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr)
		return
	}
	url, err := fileupload.FileUploader.GeneratedSignedURL(fmt.Sprintf("%s/%s", ctx.Body.Type, ctx.GetStringContextData("UserID")), ctx.Body.Permissions)
	if err != nil {
		apperrors.CustomError(ctx.Ctx, err.Error())
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "url geenraed", url, nil)
}

func SetTransactionPin(ctx *interfaces.ApplicationContext[dto.SetTransactionPinDTO]){
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr)
		return
	}
	hashedPin, err := cryptography.CryptoHahser.HashString(ctx.Body.TransactionPin)
	if err != nil {
		logger.Error(errors.New("an error occured while hashing users transaction pin"), logger.LoggerOptions{
			Key: "userID",
			Data: ctx.GetStringContextData("UserID"),
		})
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"transactionPin": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if account.TransactionPin != "" {
		server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "transaction pin has already been set", nil, nil)
		return
	}
	affected, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"),map[string]any{
		"transactionPin": string(hashedPin),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	if affected == 0 {
		apperrors.UnknownError(ctx.Ctx, errors.New("failed to update users transaction pin"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "pin set", nil, nil)
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
		apperrors.FatalServerError(ctx.Ctx, err)
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
	if err != nil {
		logger.Error(errors.New("error while deactivating user account"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		},  logger.LoggerOptions{
			Key: "success",
			Data: success,
		}, )
		apperrors.FatalServerError(ctx.Ctx, err)
	}
	if success == 0 {
		logger.Error(errors.New("error while deactivating user account"), logger.LoggerOptions{
			Key: "userID",
			Data: ctx.GetStringContextData("UserID"),
		},  logger.LoggerOptions{
			Key: "success",
			Data: success,
		},)
		apperrors.FatalServerError(ctx.Ctx, fmt.Errorf("error while deactivating user account userID - %s", ctx.GetStringContextData("UserID")))
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "deactivated", nil, nil)
}

func LogOut(ctx *interfaces.ApplicationContext[any]){
	success := cache.Cache.DeleteOne(ctx.GetStringContextData("UserID"))
	if !success {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("log out user failed - %s", ctx.GetStringContextData("UserID")))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "deactivated", nil, nil)
}
