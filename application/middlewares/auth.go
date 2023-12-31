package middlewares

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/application/utils"
	"kego.com/infrastructure/auth"
	"kego.com/infrastructure/database/repository/cache"
	"kego.com/infrastructure/logger"
)


func AuthenticationMiddleware(ctx *interfaces.ApplicationContext[any]) (*interfaces.ApplicationContext[any], bool) {
	auth_token_header := ctx.GetHeader("Authorization")
		if auth_token_header == "" || auth_token_header == nil {
			apperrors.AuthenticationError(ctx.Ctx, "provide an auth token")
			return nil, false
		}
		auth_token := strings.Split(auth_token_header.(string), " ")[1]
		valid_access_token, err := auth.DecodeAuthToken(auth_token)
		if err != nil {
			apperrors.AuthenticationError(ctx.Ctx, err.Error())
			return nil, false
		}
		if !valid_access_token.Valid {
			apperrors.AuthenticationError(ctx.Ctx, "invalid access token used")
			return nil, false
		}
		auth_token_claims := valid_access_token.Claims.(jwt.MapClaims)
		if auth_token_claims["iss"] != os.Getenv("JWT_ISSUER") {
			logger.Warning("this should trigger a wallet lock")
			apperrors.AuthenticationError(ctx.Ctx, "this is not an authorized access token")
			return nil, false
		}
		valid_token := cache.Cache.FindOne(auth_token_claims["userID"].(string))
		if valid_token == nil {
			apperrors.AuthenticationError(ctx.Ctx, "this session has expired")
			return nil, false
		}
		userRepo := repository.UserRepo()
		account, err := userRepo.FindByID(auth_token_claims["userID"].(string), options.FindOne().SetProjection(map[string]any{
			"deactivated": 1,
			"userAgent": 1,
			"deviceID": 1,
			"appVersion": 1,
		}))
		if account == nil {
			apperrors.NotFoundError(ctx.Ctx, "this account no longer exists")
			return nil, false
		}
		if account.Deactivated {
			apperrors.AuthenticationError(ctx.Ctx, "account has been deactivated")
			return nil, false
		}

		userAgent := ctx.GetHeader("User-Agent").(string)
		if auth_token_claims["appVersion"] != account.AppVersion || account.AppVersion != *utils.ExtractAppVersionFromUserAgentHeader(userAgent) ||  auth_token_claims["appVersion"] != *utils.ExtractAppVersionFromUserAgentHeader(userAgent) {
			logger.Warning("client made request using app version different from that in access token", logger.LoggerOptions{
				Key: "token appVersion",
				Data: auth_token_claims["appVersion"],
			}, logger.LoggerOptions{
				Key: "client appVersion",
				Data: account.AppVersion,
			}, logger.LoggerOptions{
				Key: "request appVersion",
				Data: *utils.ExtractAppVersionFromUserAgentHeader(userAgent),
			})
			auth.SignOutUser(ctx.Ctx, account.ID, "client made request using app version different from that in access token")
			apperrors.AuthenticationError(ctx.Ctx, "unauthorized access")
			return nil, false
		}
		deviceID := ctx.GetHeader("Polymer-Device-Id")
		if deviceID == nil {
			auth.SignOutUser(ctx.Ctx, account.ID, "client made request without a device id")
			apperrors.AuthenticationError(ctx.Ctx, "unauthorized access")
			return nil, false
		}
		if auth_token_claims["deviceID"] != account.DeviceID || account.DeviceID != deviceID.(string) ||  auth_token_claims["deviceID"] != deviceID.(string) {
			logger.Warning("client made request using device id different from that in access token",logger.LoggerOptions{
				Key: "token appVersion",
				Data: auth_token_claims["appVersion"],
			}, logger.LoggerOptions{
				Key: "client appVersion",
				Data: account.AppVersion,
			}, logger.LoggerOptions{
				Key: "request appVersion",
				Data: *utils.ExtractAppVersionFromUserAgentHeader(userAgent),
			})
			auth.SignOutUser(ctx.Ctx, account.ID, "client made request using device id different from that in access token")
			apperrors.AuthenticationError(ctx.Ctx, "unauthorized access")
			return nil, false
		}

		ctx.SetContextData("UserID", auth_token_claims["userID"])
		ctx.SetContextData("LastName", auth_token_claims["lastName"])
		ctx.SetContextData("FirstName", auth_token_claims["firstName"])
		ctx.SetContextData("Email", auth_token_claims["email"])
		ctx.SetContextData("Phone", auth_token_claims["phone"])
		ctx.SetContextData("DeviceID", auth_token_claims["deviceID"])
		ctx.SetContextData("UserAgent", auth_token_claims["userAgent"])
		ctx.SetContextData("AppVersion", auth_token_claims["appVersion"])
		return ctx, true
}