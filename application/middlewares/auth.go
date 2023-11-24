package middlewares

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
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
		}))
		if account.Deactivated {
			apperrors.AuthenticationError(ctx.Ctx, "account has been deactivated")
			return nil, false
		}
		ctx.SetContextData("UserID", auth_token_claims["userID"])
		ctx.SetContextData("Email", auth_token_claims["email"])
		ctx.SetContextData("Phone", auth_token_claims["phone"])
		ctx.SetContextData("DeviceID", auth_token_claims["deviceID"])
		ctx.SetContextData("DeviceType", auth_token_claims["deviceType"])
		return ctx, true
}