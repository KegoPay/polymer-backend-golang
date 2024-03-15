package authusecases

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"time"

	apperrors "kego.com/application/appErrors"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/database/repository/cache"
)

func InitiateKetExchange(ctx any, deviceID string, clientPublicKey *ecdh.PublicKey) []byte {
	serverPrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		apperrors.FatalServerError(ctx, err)
		return nil
	}
	sharedSecret, err := serverPrivateKey.ECDH(clientPublicKey)
	if err != nil {
		apperrors.FatalServerError(ctx, err)
		return nil
	}
	serverPublicKey := serverPrivateKey.PublicKey()
	parsedSharedSecret := string(sharedSecret)
	encryptedSecret, err := cryptography.SymmetricEncryption(string(parsedSharedSecret))
	if err != nil {
		apperrors.FatalServerError(ctx, err)
		return nil
	}
	success := cache.Cache.CreateEntry(deviceID, encryptedSecret, time.Minute * 15)
	if !success {
		apperrors.FatalServerError(ctx, nil)
		return nil
	}
	fmt.Println(cache.Cache.FindOne(deviceID))
	fmt.Println(*cache.Cache.FindOne(deviceID))
	fmt.Println(parsedSharedSecret)
	return serverPublicKey.Bytes()
}