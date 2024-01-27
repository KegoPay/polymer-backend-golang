package middlewares

import (
	"os"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/interfaces"
	"kego.com/infrastructure/cryptography"
)


func WebAgentMiddleware(ctx *interfaces.ApplicationContext[any], ipAddress string) (*interfaces.ApplicationContext[any], bool) {
	clientKey := ctx.GetHeader("Web-Api-Key")
	if clientKey == nil {
		apperrors.AuthenticationError(ctx.Ctx, "missing web token")
		return nil, false
	}
	valid := cryptography.CryptoHahser.VerifyData(os.Getenv("WEB_CLIENT_API_KEY"), clientKey.(string))
	if !valid {
		apperrors.AuthenticationError(ctx.Ctx, "invalid web token")
		return nil, false
	}
	return ctx, true
}