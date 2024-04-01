package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"kego.com/application/utils"
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

func DecryptData(encryptedData string, enc_key *string) (*string, error) {
	if enc_key == nil {
		enc_key = utils.GetStringPointer(os.Getenv("ENC_KEY"))
	}
	encryptedDataByte := []byte(encryptedData)
    c, err := aes.NewCipher([]byte(*enc_key))
    if err != nil {
		logger.Error(errors.New("error generating new cipher to decryot data"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
        return nil, err
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
		logger.Error(errors.New("error generating new gcm to decryot data"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    nonce, encryptedDataByte := encryptedDataByte[:nonceSize], encryptedDataByte[nonceSize:]

	decryptedData, err := gcm.Open(nil, nonce, encryptedDataByte, nil)
	if err != nil {
		logger.Error(errors.New("error decryptingg encrypted data"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, err
	}
    return  utils.GetStringPointer(string(decryptedData)), nil
}

func SymmetricEncryption(data string, enc_key *string) (*string, error) {
	if enc_key == nil {
		enc_key = utils.GetStringPointer(os.Getenv("ENC_KEY"))
	}
	c, err := aes.NewCipher([]byte(*enc_key))
	if err != nil {
		logger.Error(errors.New("error generating cipher to encrypt data"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, err
	}
    gcm, err := cipher.NewGCM(c)
    if err != nil {
		logger.Error(errors.New("error generating new gcm to encrypt data"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    _, err = io.ReadFull(rand.Reader, nonce)
    if err != nil {
        return nil, err
    }
	encryptedData := string(gcm.Seal(nonce, nonce, []byte(data), nil))
	fmt.Println("--")
	fmt.Println("--")
	fmt.Println("--")
	fmt.Println(encryptedData)
    return &encryptedData, nil
}