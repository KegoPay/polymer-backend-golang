package cryptography

import (
	"errors"

	"github.com/matthewhartstonge/argon2"
	"usepolymer.co/infrastructure/logger"
)

type argonHasher struct{}

func (ah argonHasher) HashString(data string) ([]byte, error) {
	config := argon2.DefaultConfig()
	raw, err := config.Hash([]byte(data), nil)
	if err != nil {
		logger.Error(errors.New("argon - error while hashing data"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, err
	}

	return raw.Encode(), nil
}

func (ah argonHasher) VerifyData(hash string, data string) bool {
	raw, err := argon2.Decode([]byte(hash))
	if err != nil {
		logger.Error(errors.New("argon - could not decode data "), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "data",
			Data: hash,
		})
	}
	ok, err := raw.Verify([]byte(data))
	if err != nil {
		logger.Error(errors.New("argon - error while verifying data"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "data",
			Data: data,
		}, logger.LoggerOptions{
			Key:  "hash",
			Data: hash,
		})
	}

	return ok
}
