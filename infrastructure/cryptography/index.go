package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"usepolymer.co/application/utils"
	"usepolymer.co/infrastructure/database/repository/cache"
	"usepolymer.co/infrastructure/logger"
)

var CryptoHahser Hasher = argonHasher{}

func GeneratePublicKey(sessionID string, clientPubKey *ecdh.PublicKey) *ecdh.PublicKey {
	serverCurve := ecdh.P256()
	serverPrivKey, err := serverCurve.GenerateKey(rand.Reader)
	if err != nil {
		logger.Error(errors.New("error generating public key for key exchange"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil
	}
	serverPubKey := serverPrivKey.PublicKey()
	serverSecret, err := serverPrivKey.ECDH(clientPubKey)
	if err != nil {
		logger.Error(errors.New("error generating server secret for key exchange"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil
	}
	cache.Cache.CreateEntry(sessionID, string(serverSecret), time.Minute*20)
	return serverPubKey
}

func DecryptData(stringToDecrypt string, keyString *string) (string, error) {
	if keyString == nil {
		keyString = utils.GetStringPointer(os.Getenv("ENC_KEY"))
	}
	key, _ := hex.DecodeString(*keyString)
	ciphertext, _ := base64.URLEncoding.DecodeString(stringToDecrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil

}

func SymmetricEncryption(payload string, keyString *string) (encryptedString string, err error) {
	// convert key to bytes
	if keyString == nil {
		keyString = utils.GetStringPointer(os.Getenv("ENC_KEY"))
	}
	key, _ := hex.DecodeString(*keyString)
	plaintext := []byte(payload)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error(err)
		panic(err.Error())
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
