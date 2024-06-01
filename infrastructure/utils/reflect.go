package utils

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"usepolymer.co/infrastructure/logger"
)

func ReflectMapToStruct[T any](payload any, s *T) error {
	err := mapstructure.Decode(payload, s)
	if err != nil {
		logger.Error(errors.New("error reflecting map to struct"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return err
	}
	return err
}
