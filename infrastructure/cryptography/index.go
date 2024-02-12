package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"

	"kego.com/infrastructure/database/repository/cache"
	"kego.com/infrastructure/logger"
)

var CryptoHahser Hasher = argonHasher{}


func GeneratePublicKey(sessionID string, clientPubKey *ecdh.PublicKey) *ecdh.PublicKey {
	serverCurve := ecdh.P256()
	serverPrivKey, err := serverCurve.GenerateKey(rand.Reader)
	if err != nil {
		logger.Error(errors.New("error generating public key for key exchange"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}
	serverPubKey := serverPrivKey.PublicKey()
	serverSecret, err := serverPrivKey.ECDH(clientPubKey)
	if err != nil {
		logger.Error(errors.New("error generating server secret for key exchange"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}
	cache.Cache.CreateEntry(sessionID, string(serverSecret), time.Minute * 20)
	return serverPubKey
}

// Encrypts data using a secret generated from the Epileptic Curve Diffie Hellman protocol
func EncryptData(secret []byte, data any) (encryptedData *string, err error) {
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    block, err := aes.NewCipher(secret)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
	marshaledData, err := json.Marshal(data)
	if err != nil {
		e := errors.New("failed to marshal payload for encryption")
		logger.Error(e, logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, err
	}
    ciphertext := gcm.Seal(nil, iv, marshaledData, nil)
    combined := append(iv, ciphertext...)
	encodedData := base64.StdEncoding.EncodeToString(combined)
    return &encodedData, nil
}