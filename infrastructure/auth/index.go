package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/database/repository/cache"
	"kego.com/infrastructure/logger"
)

const otpChars = "1234567890"

func GenerateOTP(length int, channel string) (*string, error) {
	var otp string
	if os.Getenv("ENV") == "staging" || os.Getenv("ENV") == "development" {
		otp = "000000"
	} else {
		buffer := make([]byte, length)
		_, err := rand.Read(buffer)
		if err != nil {
			return nil, err
		}
		otpCharsLength := len(otpChars)
		for i := 0; i < length; i++ {
			buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
		}
		otp = string(buffer)
	}
	otpSaved := saveOTP(channel, otp)
	if !otpSaved {
		return nil, errors.New("could not save otp")
	}
	return &otp, nil
}

func saveOTP(channel string, otp string) bool {
	hashedOTP, err := cryptography.CryptoHahser.HashString(otp)
	if err != nil {
		logger.Error(errors.New("auth module error - error while saving otp"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return false
	}
	return cache.Cache.CreateEntry(fmt.Sprintf("%s-otp", channel), string(hashedOTP), 5*time.Minute) // otp is valid for 5 mins
}

func VerifyOTP(key string, otp string) (string, bool) {
	data := cache.Cache.FindOne(fmt.Sprintf("%s-otp", key))
	if data == nil {
		logger.Info(fmt.Sprintf("%s otp not found", key))
		return "this otp has expired", false
	}
	success := cryptography.CryptoHahser.VerifyData(*data, otp)
	if !success {
		return "wrong otp provided", false
	}
	cache.Cache.DeleteOne(fmt.Sprintf("%s-otp", key))
	return "", true
}

func GenerateAuthToken(claimsData ClaimsData) (*string, error) {
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":                   os.Getenv("JWT_ISSUER"),
		"userID":                claimsData.UserID,
		"businessID":            claimsData.BusinessID,
		"exp":                   claimsData.ExpiresAt,
		"email":                 claimsData.Email,
		"phoneNum":              claimsData.PhoneNum,
		"phone":                 claimsData.Phone,
		"firstName":             claimsData.FirstName,
		"lastName":              claimsData.LastName,
		"iat":                   claimsData.IssuedAt,
		"deviceID":              claimsData.DeviceID,
		"userAgent":             claimsData.UserAgent,
		"appVersion":            claimsData.AppVersion,
		"otpIntent":             claimsData.OTPIntent,
		"pushNotificationToken": claimsData.PushNotificationToken,
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
			Key:  "error",
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

func SignOutUser(ctx any, id string, reason string) {
	logger.Info("system user signout initiated", logger.LoggerOptions{
		Key:  "reason",
		Data: reason,
	})
	deleted := cache.Cache.DeleteOne(id)
	if !deleted {
		logger.Error(errors.New("failed to sign out user"), logger.LoggerOptions{
			Key:  "id",
			Data: id,
		})
	}
}
