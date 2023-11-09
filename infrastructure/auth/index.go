package auth

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
	"kego.com/infrastructure/logger"
)

func GenerateAuthToken(claimsData ClaimsData) (*string, error) {
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":        os.Getenv("JWT_ISSUER"),
		"userID":     claimsData.UserID,
		"exp":        claimsData.ExpiresAt,
		"email":      claimsData.Email,
		"phone":      claimsData.Phone,
		"iat":        claimsData.IssuedAt,
		"deviceID":   claimsData.DeviceID,
		"deviceType": claimsData.DeviceType,
	}).SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func DecodeAuthToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Warning("this will trigger a wallet lock")
			err = errors.New("invalid token signature used")
			return nil, err
		}
		logger.Error(errors.New("error decoding jwt"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, err
	}
	if !token.Valid {
		err := errors.New("invalid token used")
		logger.Error(err)
		return nil, err
	}
	return token, nil
}