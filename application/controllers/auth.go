package controllers

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	authusecases "kego.com/application/usecases/authUsecases"
	"kego.com/entities"
	"kego.com/infrastructure/auth"
	identityverification "kego.com/infrastructure/identity_verification"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
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
		BankDetails: 	ctx.Body.BankDetails,
		BVN: 			ctx.Body.BVN,
	})
	if err != nil {
		return
	}
	otp, err := auth.GenerateOTP(6, account.Email)
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	emails.EmailService.SendEmail(account.Email, "Welcome to Kego! Verify your account to continue", "otp", map[string]interface{}{
		"FIRSTNAME": account.FirstName,
		"OTP":      otp,
	},)
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "account created", map[string]interface{}{
		"account": account,
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

func VerifyAccount(ctx *interfaces.ApplicationContext[dto.VerifyData]) {
	msg, success := auth.VerifyOTP(ctx.Body.Email, ctx.Body.Otp)
	if !success {
		apperrors.ClientError(ctx.Ctx, msg, nil)
		return
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindOneByFilter(map[string]interface{}{
		"email": ctx.Body.Email,
	}, options.FindOne().SetProjection(map[string]any{
		"firstName": 1,
		"lastName": 1,
		"bankDetails": 1,
		"email": 1,
		"phone": 1,
		"bvn": 1,
		"emailVerified": 1,
	}))
	if err != nil {
		logger.Error(err)
		apperrors.FatalServerError(ctx.Ctx)
		return
	}
	if account.EmailVerified {
		apperrors.ClientError(ctx.Ctx, "account is already verified", nil)
		return
	}
	bankCode := ""
	for _, bank := range bankssupported.SupportedBanks {
		if bank.Name == account.BankDetails.BankName{
			bankCode = bank.Code
			break
		}
	}
	updateUserPayload := map[string]any{
		"emailVerified": true,
	}
	id, code, failureReason := identityverification.IdentityVerifier.CreateAndVerifyUser(identityverification.CustomerPayload{
		Email: account.Email,
		FirstName: account.FirstName,
		LastName: account.LastName,
		Phone: account.Phone.LocalNumber,
	}, identityverification.AccountPayload{
		AccountNumber: account.BankDetails.Number,
		BankCode: bankCode,
		BVN: account.BVN,
	})
	if failureReason == ""  {
		updateUserPayload["kycCompleted"] = true
		updateUserPayload["metadata"] = map[string]any{
			"customerID": id,
			"customerCode": *code,
		}
	}else {
		updateUserPayload["kycFailedReason"] = failureReason
	}
	userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email": ctx.Body.Email,
	}, updateUserPayload )
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "otp verified", nil, nil)
}

func RetryIdentityVerification(ctx *interfaces.ApplicationContext[any]){
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(ctx.GetStringContextData("UserID"))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx)
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

	bankCode := ""
	for _, bank := range bankssupported.SupportedBanks {
		if bank.Name == account.BankDetails.BankName{
			bankCode = bank.Code
			break
		}
	}
	updateUserPayload := map[string]any{}
	if account.KYCFailedReason != nil {
		failureReason := identityverification.IdentityVerifier.VerifyUser(&identityverification.CustomerVerificationPayload{
			FirstName: account.FirstName,
			LastName: account.LastName,
			Country: "NG",
			Type: "bank_account",
			BVN: account.BVN,
			BankCode: bankCode,
			AccountNumber: account.BankDetails.Number,
		}, *account.MetaData.CustomerCode)
		if failureReason == ""  {
			updateUserPayload["kycCompleted"] = true
		}else {
			updateUserPayload["kycFailedReason"] = failureReason
		}
	}else {
		id, code, failureReason := identityverification.IdentityVerifier.CreateAndVerifyUser(identityverification.CustomerPayload{
			Email: account.Email,
			FirstName: account.FirstName,
			LastName: account.LastName,
			Phone: account.Phone.LocalNumber,
		}, identityverification.AccountPayload{
			AccountNumber: account.BankDetails.Number,
			BankCode: bankCode,
			BVN: account.BVN,
		})
		if failureReason == ""  {
			updateUserPayload["kycCompleted"] = true
			updateUserPayload["metadata"] = map[string]any{
				"customerID": id,
				"customerCode": code,
			}
		}else {
			updateUserPayload["kycFailedReason"] = failureReason
		}
	}

	userRepo.UpdatePartialByFilter(map[string]interface{}{
		"email": ctx.GetStringContextData("Email"),
	}, updateUserPayload )
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
	})
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
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "success", response, nil)
}