package services

import (
	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/services/types"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/cryptography"
)

func VerifyPin(ctx any, account *entities.User, pin string, pinType *types.PinSelectionType, device_id *string) bool {
	if pinType.Password {
		passwordMatch := cryptography.CryptoHahser.VerifyData(account.Password, pin)
		if !passwordMatch {
			apperrors.AuthenticationError(ctx, "wrong password", device_id)
			return false
		} else {
			return true
		}
	} else if pinType.TransactionPin {
		passwordMatch := cryptography.CryptoHahser.VerifyData(account.TransactionPin, pin)
		if !passwordMatch {
			apperrors.AuthenticationError(ctx, "wrong transaction pin", device_id)
			return false
		} else {
			return true
		}
	}
	apperrors.AuthenticationError(ctx, "wrong pin type selected", device_id)
	return false
}
