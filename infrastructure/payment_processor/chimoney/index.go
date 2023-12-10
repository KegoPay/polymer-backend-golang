package international_payment_processor

import (
	"encoding/json"
	"errors"
	"os"

	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/network"
)



var InternationalPaymentProcessor *ChimoneyPaymentProcessor

func InitialiseChimoneyPaymentProcessor() {
	InternationalPaymentProcessor = &ChimoneyPaymentProcessor{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("CHIMONEY_BASE_URL"),
		},
		AuthToken: os.Getenv("CHIMONEY_ACCESS_TOKEN"),
	}
}

type ChimoneyPaymentProcessor struct {
	Network *network.NetworkController
	AuthToken string
}

func (chimoneyPP *ChimoneyPaymentProcessor)GetExchangeRates() (*ChimoneyExchangeRateDTO, int, error){
	response, statusCode, err := chimoneyPP.Network.Get("/info/exchange-rates", &map[string]string{
		"X-API-KEY": chimoneyPP.AuthToken,
		"Content-Type": "application/json",
	}, nil)
	var chimoneyResponse ChimoneyExchangeRateDTO
	json.Unmarshal(*response, &chimoneyResponse)
	if err != nil {
		logger.Error(errors.New("an error occured while fetching exchange rates on chimoney"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, *statusCode, errors.New("an error occured while fetching exchange rates on chimoney")
	}
	if *statusCode != 200 {
		err = errors.New("failed to fetch exchange rates")
		logger.Error(err, logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: chimoneyResponse,
		})
		return &chimoneyResponse, *statusCode, nil
	}
	return &chimoneyResponse, *statusCode, nil
}