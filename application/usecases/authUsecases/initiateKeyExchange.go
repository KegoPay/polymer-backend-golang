package authusecases

import (
	"crypto/ecdh"
	"crypto/rand"
	"time"

	apperrors "kego.com/application/appErrors"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/database/repository/cache"
)

func InitiateKetExchange(ctx any, deviceID string, clientPublicKey *ecdh.PublicKey, device_id *string) ([]byte, error) {
	serverPrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, err
	}
	sharedSecret, err := serverPrivateKey.ECDH(clientPublicKey)
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, err
	}
	serverPublicKey := serverPrivateKey.PublicKey()
	parsedSharedSecret := string(sharedSecret)
	encryptedSecret, err := cryptography.SymmetricEncryption(string(parsedSharedSecret), nil)
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, err
	}
	success := cache.Cache.CreateEntry(deviceID, *encryptedSecret, time.Minute * 15)
	if !success {
		apperrors.FatalServerError(ctx, nil, device_id)
		return nil, err
	}
	return serverPublicKey.Bytes(), nil
}