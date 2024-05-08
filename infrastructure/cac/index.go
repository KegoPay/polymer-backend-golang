package cac_service

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/network"
)

type CACService struct {
	Network *network.NetworkController
}

func (cacs *CACService) FetchBusinessDetailsByName(name string) (*[]CACBusiness, error) {
	response, _, err := cacs.Network.Post("/searchapp/api/public-search/company-business-name-it", nil, map[string]any{
		"searchTerm": name,
	}, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while retrieving data from cac server"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, errors.New("an error occured")
	}
	var business CACBusinessNameSearchResponse
	err = json.Unmarshal(*response, &business)
	if err != nil {
		logger.Error(errors.New("an error occured while parsing data from cac server"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "response",
			Data: *response,
		})
		return nil, errors.New("something went wrong while retireving bvn data from prembly")
	}
	if business.Data == nil {
		logger.Error(errors.New("error with cac API"), logger.LoggerOptions{
			Key: "message",
			Data: business.Message,
		}, logger.LoggerOptions{
			Key: "errorCode",
			Data: business.Error,
		})
		return nil, nil
	}
	data := *business.Data
	for i, d := range data {
		d.FullAddress = strings.ReplaceAll(strings.TrimSpace(d.FullAddress), " ", " ")
		data[i] = d
	}
	return &data, nil
}

var CACServiceInstance CACService

func CreateCACServiceInstance() {
	CACServiceInstance = CACService{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("CAC_BASE_URL"),
		},
	}
}