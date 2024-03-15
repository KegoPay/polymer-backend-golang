package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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

func DecryptData(encryptedData string) (*string, error) {
    // Split the IV and ciphertext
    parts := strings.Split(encryptedData, ":")
    if len(parts) != 2 {
        return nil, errors.New("invalid encrypted data format")
    }
    iv, err := hex.DecodeString(parts[0])
    if err != nil {
        return nil, err
    }
    ciphertext, err := hex.DecodeString(parts[1])
    if err != nil {
        return nil, err
    }

    // Initialize the block cipher
    block, err := newCipherBlock(os.Getenv("ENC_KEY"))
    if err != nil {
        return nil, err
    }

    // Check if the ciphertext length is a multiple of block size
    if len(ciphertext)%aes.BlockSize != 0 {
        return nil, errors.New("ciphertext is not a multiple of the block size")
    }

    // Create a CBC decrypter
    mode := cipher.NewCBCDecrypter(block, iv)

    // Decrypt the ciphertext
    mode.CryptBlocks(ciphertext, ciphertext)

    // Unpad the plaintext
    plaintext, err := pkcs7Unpad(ciphertext, aes.BlockSize)
    if err != nil {
        return nil, err
    }

    return utils.GetStringPointer(string(plaintext)), nil
}

func SymmetricEncryption(data string) (*string, error) {
	block, err := newCipherBlock(os.Getenv("ENC_KEY"))
	if err != nil {
	   return nil, err
	}
 
   //pad plaintext
	ptbs, _ := pkcs7Pad([]byte(data), block.BlockSize())
 
	if len(ptbs)%aes.BlockSize != 0 {
	   return nil, errors.New("plaintext is not a multiple of the block size")
	}
 
	ciphertext := make([]byte, len(ptbs))
 
   //create an Initialisation vector which is the length of the block size for AES
	var iv []byte = make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	   return nil, err
	}
 
	mode := cipher.NewCBCEncrypter(block, iv)
 
   //encrypt plaintext
	mode.CryptBlocks(ciphertext, ptbs)
 
   //concatenate initialisation vector and ciphertext
	encryptedData :=  hex.EncodeToString(iv) + ":" + hex.EncodeToString(ciphertext)

	return &encryptedData, nil
}

 func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
	   return nil, errors.New("invalid blocksize")
	}
	if len(b) == 0 {
	   return nil, errors.New("invalid PKCS7 data (empty or not padded)")
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
 }
 
 func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
	   return nil, errors.New("invalid blocksize")
	}
	if len(b) == 0 {
	   return nil, errors.New("invalid PKCS7 data (empty or not padded)")
	}
 
	if len(b)%blocksize != 0 {
	   return nil, errors.New("invalid padding on input")
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
	   fmt.Println("here", n)
	   return nil, errors.New("invalid padding on input")
	}
	for i := 0; i < n; i++ {
	   if b[len(b)-n+i] != c {
		  fmt.Println("hereeee")
		  return nil, errors.New("invalid padding on input")
	   }
	}
	return b[:len(b)-n], nil
 }

 func hashWithSha256(plaintext string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, plaintext);err != nil{
	   return "", err
	}
	r := h.Sum(nil)
	return hex.EncodeToString(r), nil
 }
 
func newCipherBlock(key string) (cipher.Block, error){
	hashedKey, err := hashWithSha256(key)
	if err != nil{
	   return nil, err
	}
	bs, err := hex.DecodeString(hashedKey)
	if err != nil{
	   return nil, err
	}
	return aes.NewCipher(bs[:])
 }
