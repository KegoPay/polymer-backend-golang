package middlewares

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/interfaces"
	"kego.com/infrastructure/auth"
	"kego.com/infrastructure/database/repository/cache"
	"kego.com/infrastructure/logger"
)


func OTPTokenMiddleware(ctx *interfaces.ApplicationContext[any], ipAddress string, intent string) (*interfaces.ApplicationContext[any], bool) {
	otpToken := ctx.GetHeader("Otp-Token")
	if otpToken == nil {
		apperrors.AuthenticationError(ctx.Ctx, "missing otp token")
		return nil, false
	}

	valid_access_token, err := auth.DecodeAuthToken(otpToken.(string))
	if err != nil {
		apperrors.AuthenticationError(ctx.Ctx, err.Error())
		return nil, false
	}
	if !valid_access_token.Valid {
		apperrors.AuthenticationError(ctx.Ctx, "invalid access token used")
		return nil, false
	}
	invalidToken := cache.Cache.FindOne(otpToken.(string))
	if invalidToken != nil {
		apperrors.AuthenticationError(ctx.Ctx, "expired access token used")
		return nil, false
	}
	auth_token_claims := valid_access_token.Claims.(jwt.MapClaims)
	if auth_token_claims["iss"] != os.Getenv("JWT_ISSUER") {
		logger.Warning("this should trigger a wallet lock")
		apperrors.AuthenticationError(ctx.Ctx, "this is not an authorized access token")
		return nil, false
	}
	var channel string
	if auth_token_claims["email"] != nil {
		channel = auth_token_claims["email"].(string)
	}else {
		channel = auth_token_claims["phoneNum"].(string)
	}
	otpIntent := cache.Cache.FindOne(fmt.Sprintf("%s-otp-intent", channel))
	if otpIntent == nil {
		logger.Error(errors.New("otp intent missing"))
		apperrors.ClientError(ctx.Ctx, "otp expired", nil)
		return nil, false
	}
	if *otpIntent != auth_token_claims["otpIntent"].(string) || auth_token_claims["otpIntent"].(string) != intent{
		logger.Warning("this should trigger a wallet lock")
		logger.Error(errors.New("wrong otp intent in token"))
		apperrors.ClientError(ctx.Ctx, "incorrect intent", nil)
		return nil, false
	}
	ctx.SetContextData("OTPToken", otpToken)
	ctx.SetContextData("OTPEmail", auth_token_claims["email"])
	ctx.SetContextData("OTPPhone", auth_token_claims["phoneNum"])
	return ctx, true
}