package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/constants"
	"usepolymer.co/application/controllers/dto"
	countriessupported "usepolymer.co/application/countriesSupported"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/repository"
	"usepolymer.co/application/services"
	"usepolymer.co/application/services/types"
	authusecases "usepolymer.co/application/usecases/authUsecases"
	"usepolymer.co/application/usecases/wallet"
	"usepolymer.co/application/utils"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/auth"
	"usepolymer.co/infrastructure/background"
	"usepolymer.co/infrastructure/biometric"
	"usepolymer.co/infrastructure/cryptography"
	"usepolymer.co/infrastructure/database/repository/cache"
	fileupload "usepolymer.co/infrastructure/file_upload"
	file_upload_types "usepolymer.co/infrastructure/file_upload/types"
	identityverification "usepolymer.co/infrastructure/identity_verification"
	"usepolymer.co/infrastructure/logger"
	pushnotification "usepolymer.co/infrastructure/messaging/push_notifications"
	sms "usepolymer.co/infrastructure/messaging/whatsapp"

	server_response "usepolymer.co/infrastructure/serverResponse"
	"usepolymer.co/infrastructure/validator"
)

// func KeyExchange(ctx *interfaces.ApplicationContext[dto.KeyExchangeDTO]) {
// 	serverPublicKey, secretKey, err := authusecases.InitiateKeyExchange(ctx.Ctx, ctx.Body.DeviceID, ctx.Body.ClientPublicKey, ctx.GetHeader("Polymer-Device-Id"))
// 	if err != nil {
// 		return
// 	}
// 	payload := map[string]any{
// 		"pubKey": hex.EncodeToString(serverPublicKey),
// 	}
// 	if os.Getenv("ENV") != "prod" {
// 		payload["secret"] = secretKey
// 	}
// 	server_response.Responder.UnEncryptedRespond(ctx.Ctx, http.StatusCreated, "key exchanged", payload, nil, nil)
// }

// func EncryptForStaging(ctx *interfaces.ApplicationContext[dto.EncryptForStagingDTO]) {
// 	if os.Getenv("ENV") == "prod" {
// 		apperrors.ClientError(ctx.Ctx, "this endpoint cannot be used in a production environment", nil, utils.GetUIntPointer(401), nil)
// 		return
// 	}
// 	payloadByte, err := json.Marshal(ctx.Body.Payload)
// 	if err != nil {
// 		apperrors.FatalServerError(ctx.Ctx, err, nil)
// 		return
// 	}
// 	encrypted, err := cryptography.SymmetricEncryption(hex.EncodeToString(payloadByte), &ctx.Body.EncKey)
// 	if err != nil {
// 		apperrors.FatalServerError(ctx.Ctx, err, nil)
// 		return
// 	}
// 	server_response.Responder.UnEncryptedRespond(ctx.Ctx, http.StatusOK, "encrypted", encrypted, nil, nil)
// }

// func DecryptForStaging(ctx *interfaces.ApplicationContext[dto.DecryptForStagingDTO]) {
// 	if os.Getenv("ENV") == "prod" {
// 		apperrors.ClientError(ctx.Ctx, "this endpoint cannot be used in a production environment", nil, utils.GetUIntPointer(401), nil)
// 		return
// 	}
// 	decrypted, err := cryptography.DecryptData(ctx.Body.Payload, &ctx.Body.EncKey)
// 	if err != nil {
// 		apperrors.FatalServerError(ctx.Ctx, err, nil)
// 		return
// 	}
// 	byteBuffer := bytes.NewBuffer([]byte(decrypted))
// 	dec := gob.NewDecoder(byteBuffer) // Will read from byteBuffer
// 	var person map[string]string
// 	err = dec.Decode(&person)
// 	server_response.Responder.UnEncryptedRespond(ctx.Ctx, http.StatusOK, "decrypted", person, nil, nil)
// }

func VerifyOTP(ctx *interfaces.ApplicationContext[dto.VerifyOTPDTO]) {
	if ctx.Body.Phone == nil && ctx.Body.Email == nil {
		apperrors.ClientError(ctx.Ctx, "pass in either a phone number or email", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	var channel = ""
	var filter = map[string]any{}
	if ctx.Body.Email != nil {
		channel = *ctx.Body.Email
		filter["email"] = channel
		msg, success := auth.VerifyOTP(channel, ctx.Body.OTP)
		if !success {
			apperrors.ClientError(ctx.Ctx, msg, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
	} else {
		channel = *ctx.Body.Phone
		filter["phone.localNumber"] = channel
		msg, success := auth.VerifyOTP(channel, ctx.Body.OTP)
		if !success {
			logger.Info("possible sms otp attempted to be verified as whatsapp otp", logger.LoggerOptions{
				Key:  "message",
				Data: msg,
			})
			otpRef := cache.Cache.FindOne(fmt.Sprintf("%s-sms-otp-ref", channel))
			if otpRef == nil {
				apperrors.NotFoundError(ctx.Ctx, "otp has expired", ctx.GetHeader("Polymer-Device-Id"))
				return
			}
			d, err := cryptography.DecryptData(*otpRef, nil)
			logger.Error(errors.New("error dcrypting sms otp ref"), logger.LoggerOptions{
				Key:  "ref",
				Data: *otpRef,
			}, logger.LoggerOptions{
				Key:  "channel",
				Data: channel,
			}, logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			success := sms.SMSService.VerifyOTP(d, ctx.Body.OTP)
			if !success {
				apperrors.ClientError(ctx.Ctx, "wrong otp", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
				return
			}
			cache.Cache.DeleteOne(fmt.Sprintf("%s-sms-otp-ref", channel))
		}
	}
	otpIntent := cache.Cache.FindOne(fmt.Sprintf("%s-otp-intent", channel))
	if otpIntent == nil {
		logger.Error(errors.New("otp intent missing"))
		apperrors.ClientError(ctx.Ctx, "otp expired", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email:     ctx.Body.Email,
		PhoneNum:  ctx.Body.Phone,
		OTPIntent: *otpIntent,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(15)).Unix(), //lasts for 10 mins
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "otp verified", token, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func CreateAccount(ctx *interfaces.ApplicationContext[dto.CreateAccountDTO]) {
	account, _, err := authusecases.CreateAccount(ctx.Ctx, &entities.User{
		Email:                 ctx.Body.Email,
		Password:              ctx.Body.Password,
		UserAgent:             ctx.Body.UserAgent,
		DeviceID:              ctx.Body.DeviceID,
		AppVersion:            ctx.Body.AppVersion,
		PushNotificationToken: ctx.Body.PushNotificationToken,
		Tier:                  0,
		Longitude:             ctx.GetFloat64ContextData("Longitude"),
		Latitude:              ctx.GetFloat64ContextData("Latitude"),
	}, &ctx.Body.DeviceID, ctx.Body.AuthOne)
	if err != nil {
		return
	}
	cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), 3, time.Hour*24*100)
	cache.Cache.CreateEntry(fmt.Sprintf("%s-nin-kyc-attempts-left", account.Email), 2, time.Hour*24*100)
	cache.Cache.CreateEntry(fmt.Sprintf("%s-bvn-kyc-attempts-left", account.Email), 2, time.Hour*24*100)
	if !ctx.Body.AuthOne {
		otp, err := auth.GenerateOTP(6, account.Email)
		if err != nil {
			apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
		background.Scheduler.Emit("send_email", map[string]any{
			"email":        account.Email,
			"subject":      "Welcome to Polymer! Verify your account to continue",
			"templateName": "otp",
			"opts": map[string]interface{}{
				"FIRSTNAME": account.FirstName,
				"OTP":       otp,
			},
		})
		cache.Cache.CreateEntry(fmt.Sprintf("%s-otp-intent", ctx.Body.Email), "verify_account", time.Minute*10)
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "account created", nil, nil, nil, &ctx.Body.DeviceID)
}

func LoginUser(ctx *interfaces.ApplicationContext[dto.LoginDTO]) {
	appVersion := utils.ExtractAppVersionFromUserAgentHeader(*ctx.GetHeader("User-Agent"))
	if appVersion == nil {
		apperrors.UnsupportedAppVersion(ctx.Ctx, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	account, wallet, token := authusecases.LoginAccount(ctx.Ctx, ctx.Body.Email, ctx.Body.Phone, &ctx.Body.Password, *appVersion, *ctx.GetHeader("User-Agent"), ctx.Body.DeviceID, ctx.Body.PushNotificationToken, ctx.GetFloat64ContextData("Longitude"), ctx.GetFloat64ContextData("Latitude"), ctx.GetHeader("Polymer-Device-Id"))
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
	businessRepo := repository.BusinessRepo()
	business, err := businessRepo.FindOneByFilter(map[string]interface{}{
		"userID": account.ID,
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	walletRepo := repository.WalletRepo()
	var businessWallet *entities.Wallet
	if business != nil {
		businessWallet, err = walletRepo.FindOneByFilter(map[string]interface{}{
			"businessID": business.ID,
		})
		if err != nil {
			apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
	}
	responsePayload := map[string]interface{}{
		"account":        account,
		"wallet":         wallet,
		"token":          token,
		"country":        country,
		"business":       business,
		"businessWallet": businessWallet,
	}
	if account.TransactionPin == "" {
		responsePayload["unsetTrxPin"] = true
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "login successful", responsePayload, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func ResetPassword(ctx *interfaces.ApplicationContext[dto.ResetPasswordDTO]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("OTPEmail"),
	})
	if err != nil {
		logger.Error(errors.New("error fetching a user account to reset password"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "account with email not found", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	account.Password = ctx.Body.NewPassword
	valiedationErr := validator.ValidatorInstance.ValidateStruct(*account)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	hashedPassword, err := cryptography.CryptoHahser.HashString(ctx.Body.NewPassword, nil)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	success, err := userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("OTPEmail"),
	}, map[string]interface{}{
		"password": string(hashedPassword),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	if !success {
		logger.Error(errors.New("could not reset password"), logger.LoggerOptions{
			Key:  "email",
			Data: ctx.GetStringContextData("OTPEmail"),
		})
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	cache.Cache.CreateEntry(ctx.GetStringContextData("OTPToken"), true, time.Minute*5)
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "password reset", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func UpdatePassword(ctx *interfaces.ApplicationContext[dto.UpdatePassword]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("Email"),
	})
	if err != nil {
		logger.Error(errors.New("error fetching a user account to reset password"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "account with email not found", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	success := cryptography.CryptoHahser.VerifyData(account.Password, ctx.Body.CurrentPassword)
	if !success {
		apperrors.ClientError(ctx.Ctx, "incorrect password", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	account.Password = ctx.Body.NewPassword
	valiedationErr := validator.ValidatorInstance.ValidateStruct(*account)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	hashed_password, err := cryptography.CryptoHahser.HashString(ctx.Body.NewPassword, nil)
	if err != nil {
		logger.Error(errors.New("error hashing users new password"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	modified, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]interface{}{
		"password": string(hashed_password),
	})
	if err != nil {
		logger.Error(errors.New("error while updating user password"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
	}
	if modified == 0 {
		logger.Error(errors.New("error while updating user password"), logger.LoggerOptions{
			Key:  "modified",
			Data: modified,
		})
		apperrors.FatalServerError(ctx.Ctx, fmt.Errorf("failed to update users password userID %s", ctx.GetStringContextData("UserID")), ctx.GetHeader("Polymer-Device-Id"))
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "password updated", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func ResendOTP(ctx *interfaces.ApplicationContext[dto.ResendOTP]) {
	if ctx.Body.Phone == nil && ctx.Body.Email == nil {
		apperrors.ClientError(ctx.Ctx, "pass in either a phone number or email", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	var channel = ""
	var filter = map[string]any{}
	if ctx.Body.Email != nil {
		channel = *ctx.Body.Email
		filter["email"] = channel
	} else {
		channel = *ctx.Body.Phone
		filter["phone.localNumber"] = channel
	}

	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(filter, options.FindOne().SetProjection(map[string]any{
		"firstName": 1,
		"phone":     1,
		"email":     1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "User not found", ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	if ctx.Body.Email != nil {
		otp, err := auth.GenerateOTP(6, channel)
		if err != nil {
			apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
		background.Scheduler.Emit("send_email", map[string]any{
			"email":        account.Email,
			"subject":      "An OTP was requested for your account",
			"templateName": "otp",
			"opts": map[string]interface{}{
				"FIRSTNAME": account.FirstName,
				"OTP":       otp,
			},
		})
	} else if ctx.Body.Phone != nil {
		var otp *string
		var err error
		if ctx.Body.Whatsapp != nil || account.Phone.WhatsApp {
			if account.Phone.WhatsApp || (ctx.Body.Whatsapp != nil && *ctx.Body.Whatsapp) {
				otp, err = auth.GenerateOTP(6, channel)
				if err != nil {
					apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
					return
				}
			}
		}
		ref := sms.SMSService.SendOTP(fmt.Sprintf("%s%s", account.Phone.Prefix, account.Phone.LocalNumber), account.Phone.WhatsApp || (ctx.Body.Whatsapp != nil && *ctx.Body.Whatsapp), otp)
		encryptedRef, err := cryptography.SymmetricEncryption(*ref, nil)
		if err != nil {
			apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
			return
		}
		cache.Cache.CreateEntry(fmt.Sprintf("%s-sms-otp-ref", channel), encryptedRef, time.Minute*10)
	}
	cache.Cache.CreateEntry(fmt.Sprintf("%s-otp-intent", channel), ctx.Body.Intent, time.Minute*10)
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "otp sent", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func VerifyEmail(ctx *interfaces.ApplicationContext[any]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("OTPEmail"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "this account no longer exists", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account.EmailVerified {
		apperrors.ClientError(ctx.Ctx, "this email has already been verified", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email:                 &account.Email,
		Phone:                 account.Phone,
		UserID:                account.ID,
		IssuedAt:              time.Now().Unix(),
		ExpiresAt:             time.Now().Local().Add(time.Minute * time.Duration(15)).Unix(), //lasts for 15 mins
		UserAgent:             account.UserAgent,
		FirstName:             account.FirstName,
		LastName:              account.LastName,
		DeviceID:              account.DeviceID,
		AppVersion:            account.AppVersion,
		PushNotificationToken: account.PushNotificationToken,
	})
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	account.EmailVerified = true
	success, err := userRepo.UpdateByID(account.ID, account)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !success {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	cache.Cache.CreateEntry(ctx.GetStringContextData("OTPToken"), true, time.Minute*5)
	hashedToken, err := cryptography.CryptoHahser.HashString(*token, nil)
	if err != nil {
		apperrors.FatalServerError(ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	cache.Cache.CreateEntry(account.ID, hashedToken, time.Minute*time.Duration(10)) // cache authentication token for 10 mins
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "account verified", map[string]string{
		"token": *token,
	}, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func VerifyPhone(ctx *interfaces.ApplicationContext[dto.IsAuthOne]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"phone.localNumber": ctx.GetStringContextData("OTPPhone"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "this number is not assigned to any account", ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	if ctx.Body.AuthOne {
		if account.Phone.IsVerified && account.NextOfKin != nil && account.BVN != "" && account.NIN != "" {
			if account.Tier == 2 {
				account.Tier = 3
			}
		} else if (account.Phone.IsVerified && account.NextOfKin != nil) && (account.BVN == "" || account.NIN == "") {
			if account.Tier == 1 {
				account.Tier = 2
			}
		}
	}
	account.Phone.IsVerified = true
	success, err := userRepo.UpdateByID(account.ID, account)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !success {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "phone verified", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func VerifyAccount(ctx *interfaces.ApplicationContext[dto.VerifyAccountData]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	attemptsLeft := cache.Cache.FindOne(fmt.Sprintf("%s-kyc-attempts-left", ctx.GetStringContextData("Email")))
	if attemptsLeft == nil {
		apperrors.ClientError(ctx.Ctx, `You’ve reach the maximum number of tries allowed for this.`, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	parsedAttemptsLeft, err := strconv.Atoi(*attemptsLeft)
	if err != nil {
		apperrors.ClientError(ctx.Ctx, `You’ve reach the maximum number of tries allowed for this.`, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if parsedAttemptsLeft == 0 {
		apperrors.ClientError(ctx.Ctx, `You’ve reach the maximum number of tries allowed for this.`, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("Email"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("Account with email %s does not exist. Please contact support on %s to help resolve this issue.", ctx.GetStringContextData("Email"), constants.SUPPORT_EMAIL), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !account.EmailVerified {
		apperrors.ClientError(ctx.Ctx, "verify your email before attempting identity verification", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	if account.KYCCompleted {
		apperrors.ClientError(ctx.Ctx, "you have completed your identity verification", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	url, err := fileupload.FileUploader.GeneratedSignedURL(ctx.Body.ProfileImage, file_upload_types.SignedURLPermission{
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

	success, err := biometric.BiometricService.LivenessCheck(url)
	if err != nil {
		apperrors.ClientError(ctx.Ctx, "something went wrong while performing face verification", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !success {
		cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft-1, time.Hour*24*365) // keep data cached for a yea
		apperrors.ClientError(ctx.Ctx, "Face verification failed. Please ensure you are in a well lit environment and have no coverings on your face.", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	kycDetails := struct {
		Gender        string
		WatchListed   *string
		FirstName     string
		MiddleName    *string
		LastName      string
		DateOfBirth   string
		PhoneNumber   *string
		Nationality   string
		Base64Image   string
		Address       string
		StateOfOrigin string
		LGAOfOrigin   string
		Title         string
	}{}
	if ctx.Body.Path == "bvn" {
		bvnDetails, err := identityverification.IdentityVerifier.FetchBVNDetails(ctx.Body.ID)
		if err != nil {
			cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft-1, time.Hour*24*365) // keep data cached for a yea
			apperrors.CustomError(ctx.Ctx, err.Error(), ctx.GetHeader("Polymer-Device-Id"))
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
		kycDetails.Title = bvnDetails.Title
		kycDetails.LGAOfOrigin = bvnDetails.LGAOfOrigin
		kycDetails.StateOfOrigin = bvnDetails.StateOfOrigin
	} else if ctx.Body.Path == "nin" {
		ninDetails, err := identityverification.IdentityVerifier.FetchNINDetails(ctx.Body.ID)
		if err != nil {
			cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft-1, time.Hour*24*365) // keep data cached for a yea
			apperrors.CustomError(ctx.Ctx, err.Error(), ctx.GetHeader("Polymer-Device-Id"))
			return
		}
		kycDetails.Base64Image = ninDetails.Base64Image
		kycDetails.FirstName = ninDetails.FirstName
		kycDetails.MiddleName = ninDetails.MiddleName
		kycDetails.LastName = ninDetails.LastName
		kycDetails.Gender = ninDetails.Gender
		kycDetails.PhoneNumber = ninDetails.PhoneNumber
		kycDetails.DateOfBirth = ninDetails.DateOfBirth
		kycDetails.Address = ninDetails.Address
	} else {
		apperrors.ClientError(ctx.Ctx, "unknown id type selected", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	result, err := biometric.BiometricService.FaceMatch(url, &kycDetails.Base64Image)
	if err != nil {
		apperrors.ClientError(ctx.Ctx, "something went wrong while performing face verification", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !result {
		cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-attempts-left", account.Email), parsedAttemptsLeft-1, time.Hour*24*365) // keep data cached for a yea
		apperrors.ClientError(ctx.Ctx, "We compared your face with that on your ID and it did not match. Please ensure you are in a well lit environment and have no coverings on your face.", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	watchListed := false
	if kycDetails.WatchListed != nil {
		if *kycDetails.WatchListed == "True" {
			watchListed = true
		}
	}
	profileImageFilePath := fmt.Sprintf("%s/%s", ctx.GetStringContextData("UserID"), "profile-image")
	fileupload.FileUploader.UploadBase64File(profileImageFilePath, &kycDetails.Base64Image)
	encryptedID, err := cryptography.CryptoHahser.HashString(ctx.Body.ID, []byte{})
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userUpdatedInfo := map[string]any{
		"gender":    kycDetails.Gender,
		"dob":       kycDetails.DateOfBirth,
		"lastName":  cases.Title(language.Und).String(kycDetails.LastName),
		"firstName": cases.Title(language.Und).String(kycDetails.FirstName),
		"middleName": func() *string {
			if kycDetails.MiddleName != nil {
				return utils.GetStringPointer(cases.Title(language.Und).String(*kycDetails.MiddleName))
			}
			return nil
		}(),
		"watchListed": watchListed,
		"nationality": kycDetails.Nationality,
		"phone": func() *entities.PhoneNumber {
			if kycDetails.PhoneNumber != nil {
				return &entities.PhoneNumber{
					Prefix:      "234",
					ISOCode:     "NG",
					LocalNumber: *kycDetails.PhoneNumber,
				}
			}
			return nil
		}(),
		"profileImage": ctx.Body.ProfileImage,
		"kycCompleted": true,
		"bvn": func() *string {
			if ctx.Body.Path == "bvn" {
				return utils.GetStringPointer(string(encryptedID))
			}
			return nil
		}(),
		"nin": func() *string {
			if ctx.Body.Path == "nin" {
				return utils.GetStringPointer(string(encryptedID))
			}
			return nil
		}(),
		"accountRestricted": watchListed,
		"address": entities.Address{
			FullAddress: &kycDetails.Address,
		},
		"tier": 1,
	}
	success, err = userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("Email"),
	}, userUpdatedInfo)
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if !success {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if ctx.Body.Path == "bvn" {
		wallet.GenerateNGNDVA(ctx.Ctx, account.WalletID, kycDetails.FirstName, kycDetails.LastName, ctx.GetStringContextData("Email"), *utils.GenerateDummyKYCID(), *utils.GenerateDummyKYCID(), ctx.GetHeader("Polymer-Device-Id"))
	} else {
		wallet.GenerateNGNDVA(ctx.Ctx, account.WalletID, kycDetails.FirstName, kycDetails.LastName, ctx.GetStringContextData("Email"), *utils.GenerateDummyKYCID(), *utils.GenerateDummyKYCID(), ctx.GetHeader("Polymer-Device-Id"))
	}
	cache.Cache.DeleteOne(fmt.Sprintf("%s-kyc-attempts-left", ctx.GetStringContextData("Email")))
	cache.Cache.DeleteOne(fmt.Sprintf("%s-%s-kyc-attempts-left", ctx.GetStringContextData("Email"), ctx.Body.Path))
	now := time.Now()
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email: utils.GetStringPointer(ctx.GetStringContextData("Email")),
		Phone: func() *entities.PhoneNumber {
			if kycDetails.PhoneNumber != nil {
				return &entities.PhoneNumber{
					Prefix:      "234",
					ISOCode:     "NG",
					LocalNumber: *kycDetails.PhoneNumber,
				}
			}
			return nil
		}(),
		UserID:                ctx.GetStringContextData("UserID"),
		IssuedAt:              now.Unix(),
		ExpiresAt:             now.Local().Add(time.Minute * time.Duration(15)).Unix(), //lasts for 10 mins
		UserAgent:             ctx.GetStringContextData("UserAgent"),
		FirstName:             userUpdatedInfo["firstName"].(string),
		LastName:              userUpdatedInfo["lastName"].(string),
		DeviceID:              ctx.GetStringContextData("DeviceID"),
		PushNotificationToken: ctx.GetStringContextData("PushNotificationToken"),
		AppVersion:            ctx.GetStringContextData("AppVersion"),
	})
	if err != nil {
		apperrors.UnknownError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	pushnotification.PushNotificationService.PushOne(ctx.GetStringContextData("PushNotificationToken"), "Welcome to Polymer!😃",
		"You now have global payments at your finger tips! Make payments with crypto, Mobile Money and to bank accounts in over 40+ countries!🤯")
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "kyc completed", token, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func AccountWithEmailExists(ctx *interfaces.ApplicationContext[any]) {
	email := ctx.Query["email"]
	if email == "" {
		server_response.Responder.Respond(ctx.Ctx, http.StatusBadRequest, "pass in a valid email", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": email,
	}, options.FindOne().SetProjection(map[string]interface{}{
		"emailVerified": 1,
		"kycCompleted":  1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	response := map[string]any{}
	if account == nil {
		response["exists"] = false
	} else {
		response["exists"] = true
		response["emailVerified"] = account.EmailVerified
		response["KYCCompleted"] = account.KYCCompleted
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "success", response, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func SetTransactionPin(ctx *interfaces.ApplicationContext[dto.SetTransactionPinDTO]) {
	valiedationErr := validator.ValidatorInstance.ValidateStruct(ctx.Body)
	if valiedationErr != nil {
		apperrors.ValidationFailedError(ctx.Ctx, valiedationErr, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	hashedPin, err := cryptography.CryptoHahser.HashString(ctx.Body.TransactionPin, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while hashing users transaction pin"), logger.LoggerOptions{
			Key:  "userID",
			Data: ctx.GetStringContextData("UserID"),
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"transactionPin": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account.TransactionPin != "" {
		server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "transaction pin has already been set", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	affected, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]any{
		"transactionPin": string(hashedPin),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if affected == 0 {
		apperrors.UnknownError(ctx.Ctx, errors.New("failed to update users transaction pin"), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "pin set", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func DeactivateAccount(ctx *interfaces.ApplicationContext[dto.ConfirmPin]) {
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(ctx.GetStringContextData("UserID"), options.FindOne().SetProjection(map[string]any{
		"deactivated":    1,
		"password":       1,
		"transactionPin": 1,
	}))
	if err != nil {
		logger.Error(errors.New("error fetching a user account to deactivate account"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("This user profile was not found. Please contact support on %s to help resolve this issue.", constants.SUPPORT_EMAIL), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if account.Deactivated {
		apperrors.ClientError(ctx.Ctx, "account has already been deactivated", nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	match := services.VerifyPin(ctx.Ctx, account, ctx.Body.Pin, &types.PinSelectionType{
		Password: true,
	}, ctx.GetHeader("Polymer-Device-Id"))
	if !match {
		return
	}
	success, err := userRepo.UpdatePartialByID(ctx.GetStringContextData("UserID"), map[string]interface{}{
		"deactivated": true,
	})
	if err != nil {
		logger.Error(errors.New("error while deactivating user account"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "success",
			Data: success,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
	}
	if success == 0 {
		logger.Error(errors.New("error while deactivating user account"), logger.LoggerOptions{
			Key:  "userID",
			Data: ctx.GetStringContextData("UserID"),
		}, logger.LoggerOptions{
			Key:  "success",
			Data: success,
		})
		apperrors.FatalServerError(ctx.Ctx, fmt.Errorf("error while deactivating user account userID - %s", ctx.GetStringContextData("UserID")), ctx.GetHeader("Polymer-Device-Id"))
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "deactivated", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func LogOut(ctx *interfaces.ApplicationContext[any]) {
	success := cache.Cache.DeleteOne(ctx.GetStringContextData("UserID"))
	if !success {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("log out user failed - %s", ctx.GetStringContextData("UserID")), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "logged out", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}
