package middlewares

import (
	"usepolymer.co/application/interfaces"
)

func InhouseAuthMiddleware(ctx *interfaces.ApplicationContext[any], signingKey string, issuer string) (*interfaces.ApplicationContext[any], bool) {
	// authTokenPointer := ctx.GetHeader("X-Is-Token")
	// if authTokenPointer == nil {
	// 	apperrors.AuthenticationError(ctx.Ctx, "provide an auth token", nil)
	// 	return nil, false
	// }
	// authToken := *authTokenPointer
	// valid_access_token, err := auth.DecodeInterserviceAuthToken(authToken, signingKey)
	// if err != nil {
	// 	apperrors.AuthenticationError(ctx.Ctx, "this session has expired", nil)
	// 	return nil, false
	// }
	// if !valid_access_token.Valid {
	// 	apperrors.AuthenticationError(ctx.Ctx, "invalid access token used", nil)
	// 	return nil, false
	// }
	// auth_token_claims := valid_access_token.Claims.(jwt.MapClaims)
	// if auth_token_claims["iss"] != issuer {
	// 	apperrors.AuthenticationError(ctx.Ctx, "this is not an authorized access token", nil)
	// 	return nil, false
	// }
	return ctx, true
}
