package middlewares

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/repository"
	"usepolymer.co/application/utils"
	"usepolymer.co/infrastructure/auth"
	"usepolymer.co/infrastructure/background"
	"usepolymer.co/infrastructure/cryptography"
	"usepolymer.co/infrastructure/database/repository/cache"
	"usepolymer.co/infrastructure/logger"
)

func AuthenticationMiddleware(ctx *interfaces.ApplicationContext[any], restricted bool, business_route bool) (*interfaces.ApplicationContext[any], bool) {
	authTokenHeaderPointer := ctx.GetHeader("Authorization")
	if authTokenHeaderPointer == nil {
		apperrors.AuthenticationError(ctx.Ctx, "provide an auth token", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	auth_token := strings.Split(*authTokenHeaderPointer, " ")[1]
	valid_access_token, err := auth.DecodeAuthToken(auth_token)
	if err != nil {
		apperrors.AuthenticationError(ctx.Ctx, "this session has expired", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	if !valid_access_token.Valid {
		apperrors.AuthenticationError(ctx.Ctx, "invalid access token used", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	auth_token_claims := valid_access_token.Claims.(jwt.MapClaims)
	if auth_token_claims["iss"] != os.Getenv("JWT_ISSUER") {
		logger.Warning("this triggers a wallet lock")
		background.Scheduler.Emit("lock_account", map[string]any{
			"id": auth_token_claims["userID"],
		})
		apperrors.AuthenticationError(ctx.Ctx, "this is not an authorized access token", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	valid_token := cache.Cache.FindOne(auth_token_claims["userID"].(string))
	if valid_token == nil {
		apperrors.AuthenticationError(ctx.Ctx, "this session has expired", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	match := cryptography.CryptoHahser.VerifyData(*valid_token, auth_token)
	if !match {
		apperrors.AuthenticationError(ctx.Ctx, "this session has expired", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	userRepo := repository.UserRepo()
	account, err := userRepo.FindByID(auth_token_claims["userID"].(string), options.FindOne().SetProjection(map[string]any{
		"deactivated":         1,
		"userAgent":           1,
		"deviceID":            1,
		"appVersion":          1,
		"notificationOptions": 1,
		"kycCompleted":        1,
		"tier":                1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	if account == nil {
		apperrors.NotFoundError(ctx.Ctx, "this account no longer exists", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	if account.Deactivated {
		apperrors.AuthenticationError(ctx.Ctx, "account has been deactivated", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}

	if restricted && !account.KYCCompleted {
		apperrors.AuthenticationError(ctx.Ctx, "verify your bvn before attempting this action", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}

	userAgent := ctx.GetHeader("User-Agent")
	if auth_token_claims["appVersion"] != account.AppVersion || account.AppVersion != *utils.ExtractAppVersionFromUserAgentHeader(*userAgent) || auth_token_claims["appVersion"] != *utils.ExtractAppVersionFromUserAgentHeader(*userAgent) {
		logger.Warning("client made request using app version different from that in access token", logger.LoggerOptions{
			Key:  "token appVersion",
			Data: auth_token_claims["appVersion"],
		}, logger.LoggerOptions{
			Key:  "client appVersion",
			Data: account.AppVersion,
		}, logger.LoggerOptions{
			Key:  "request appVersion",
			Data: *utils.ExtractAppVersionFromUserAgentHeader(*userAgent),
		})
		logger.Warning("this triggers a wallet lock")
		background.Scheduler.Emit("lock_account", map[string]any{
			"id": auth_token_claims["userID"],
		})
		auth.SignOutUser(ctx.Ctx, account.ID, "client made request using app version different from that in access token")
		apperrors.AuthenticationError(ctx.Ctx, "unauthorized access", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	deviceID := ctx.GetHeader("Polymer-Device-Id")
	if deviceID == nil {
		auth.SignOutUser(ctx.Ctx, account.ID, "client made request without a device id")
		logger.Info("device id missing from client")
		apperrors.AuthenticationError(ctx.Ctx, "unauthorized access", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	if auth_token_claims["deviceID"] != account.DeviceID || account.DeviceID != *deviceID || auth_token_claims["deviceID"] != *deviceID {
		logger.Warning("client made request using device id different from that in access token", logger.LoggerOptions{
			Key:  "token device id",
			Data: auth_token_claims["deviceID"],
		}, logger.LoggerOptions{
			Key:  "client  device id",
			Data: account.DeviceID,
		}, logger.LoggerOptions{
			Key:  "request  device id",
			Data: deviceID,
		})
		logger.Warning("this triggers a wallet lock")
		background.Scheduler.Emit("lock_account", map[string]any{
			"id": auth_token_claims["userID"],
		})
		auth.SignOutUser(ctx.Ctx, account.ID, "client made request using device id different from that in access token")
		apperrors.AuthenticationError(ctx.Ctx, "unauthorized access", ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}

	if business_route {
		businessRepo := repository.BusinessRepo()
		exists, err := businessRepo.CountDocs(map[string]interface{}{
			"_id":    ctx.GetStringParameter("businessID"),
			"userID": account.ID,
		})
		if err != nil {
			apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
			return nil, false
		}
		if exists == 0 {
			apperrors.NotFoundError(ctx.Ctx, "business not found", ctx.GetHeader("Polymer-Device-Id"))
			return nil, false
		}
	}

	ctx.SetContextData("UserID", auth_token_claims["userID"])
	ctx.SetContextData("LastName", auth_token_claims["lastName"])
	ctx.SetContextData("FirstName", auth_token_claims["firstName"])
	ctx.SetContextData("Email", auth_token_claims["email"])
	ctx.SetContextData("Phone", auth_token_claims["phone"])
	ctx.SetContextData("DeviceID", auth_token_claims["deviceID"])
	ctx.SetContextData("UserAgent", auth_token_claims["userAgent"])
	ctx.SetContextData("AppVersion", auth_token_claims["appVersion"])
	ctx.SetContextData("PushNotificationToken", auth_token_claims["pushNotificationToken"])
	ctx.SetContextData("EmailOptions", account.NotificationOptions.Emails)
	ctx.SetContextData("PushNotifOptions", account.NotificationOptions.PushNotification)
	ctx.SetContextData("Tier", account.Tier)
	return ctx, true
}
