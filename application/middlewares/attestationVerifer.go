package middlewares

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/interfaces"
	"kego.com/infrastructure/logger"
)

var jwksURL = "https://firebaseappcheck.googleapis.com/v1beta/jwks"

func AttestationVerifier(ctx *interfaces.ApplicationContext[any]) (*interfaces.ApplicationContext[any], bool) {
	attestationToken := ctx.GetHeader("X-Firebase-Token")
	if attestationToken == nil {
		apperrors.AuthenticationError(ctx.Ctx, "attestation token missing",  ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	if os.Getenv("ENV") == "development" {
		return ctx, true
	}
	options := keyfunc.Options{
		Ctx: context.TODO(),
		RefreshErrorHandler: func(err error) {
			logger.Error(errors.New("there was an error with the jwt.Keyfunc"), logger.LoggerOptions{
				Key: "error",
				Data: err,
			})
		},
		RefreshInterval: time.Hour * 6,
	}
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		logger.Error(errors.New("failed to create JWKS from resource at the given URL"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.AuthenticationError(ctx.Ctx, "client verification failed",  ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}

	payload, err := jwt.Parse(*attestationToken, jwks.Keyfunc)
	if err != nil {
		logger.Error(errors.New("failed to parse token"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.AuthenticationError(ctx.Ctx, "client verification failed",  ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}

	if !payload.Valid {
		logger.Error(errors.New("invalid attestation foken"), logger.LoggerOptions{
			Key: "token",
			Data: attestationToken,
		})
		apperrors.AuthenticationError(ctx.Ctx, "client verification failed",  ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	} else if payload.Header["alg"] != "RS256" {
		// Ensure the token's header uses the algorithm RS256
		logger.Error(errors.New("invalid attestation token algorithm"), logger.LoggerOptions{
			Key: "token",
			Data: attestationToken,
		})
		apperrors.AuthenticationError(ctx.Ctx, "client verification failed",  ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	} else if payload.Header["typ"] != "JWT" {
		// Ensure the token's header has type JWT
		logger.Error(errors.New("invalid attestation token type"), logger.LoggerOptions{
			Key: "token",
			Data: attestationToken,
		})
		apperrors.AuthenticationError(ctx.Ctx, "client verification failed",  ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	} else if !verifyAudClaim(payload.Claims.(jwt.MapClaims)["aud"].([]interface{})) {
		// Ensure the token's audience matches your project
		logger.Error(errors.New("invalid attestation token audience"), logger.LoggerOptions{
			Key: "token",
			Data: attestationToken,
		})
		apperrors.AuthenticationError(ctx.Ctx, "client verification failed",  ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	} else if !strings.Contains(payload.Claims.(jwt.MapClaims)["iss"].(string),
		"https://firebaseappcheck.googleapis.com/"+os.Getenv("PROJECT_NUMBER")) {
		// Ensure the token is issued by App Check
		logger.Error(errors.New("invalid attestation token issuer"), logger.LoggerOptions{
			Key: "token",
			Data: attestationToken,
		})
		apperrors.AuthenticationError(ctx.Ctx, "client verification failed",  ctx.GetHeader("Polymer-Device-Id"))
		return nil, false
	}
	jwks.EndBackground()
	return ctx, true
}

func verifyAudClaim(auds []interface{}) bool {
	for _, aud := range auds {
		if aud == "projects/"+os.Getenv("PROJECT_NUMBER") {
			return true
		}
	}
	return false
}