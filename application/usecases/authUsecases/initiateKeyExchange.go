package authusecases

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/infrastructure/cryptography"
	"usepolymer.co/infrastructure/database/repository/cache"
)

func InitiateKeyExchange(ctx any, deviceID string, clientPublicKey *ecdh.PublicKey, device_id *string) ([]byte, *string, error) {
	serverPrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, nil, err
	}
	sharedSecret, err := serverPrivateKey.ECDH(clientPublicKey)
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, nil, err
	}

	serverPublicKey := serverPrivateKey.PublicKey()

	parsedSharedSecret := hex.EncodeToString(sharedSecret)[:32]

	encryptedSecret, err := cryptography.SymmetricEncryption(parsedSharedSecret, nil)
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, nil, err
	}
	success := cache.Cache.CreateEntry(fmt.Sprintf("%s-key", deviceID), encryptedSecret, time.Minute*15)
	if !success {
		apperrors.FatalServerError(ctx, nil, device_id)
		return nil, nil, err
	}
	return serverPublicKey.Bytes(), &parsedSharedSecret, nil
}
