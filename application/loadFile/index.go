package loadfile

import (
	"errors"

	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/network"
)

func LoadImageFromURL(url string) *[]byte {
	n := network.NetworkController{}
	data, statusCode, err := n.Get(url, nil, nil)
	if err != nil {
		logger.Error(errors.New("could not load file from url"), logger.LoggerOptions{
			Key:  "error",
			Data: url,
		})
		return nil
	}
	if *statusCode != 200 {
		logger.Error(errors.New("could not load file from url"), logger.LoggerOptions{
			Key:  "url",
			Data: url,
		}, logger.LoggerOptions{
			Key:  "status code",
			Data: statusCode,
		})
		return nil
	}
	return data
}
