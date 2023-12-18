package international_payment_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"kego.com/entities"
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

func (chimoneyPP *ChimoneyPaymentProcessor)GetExchangeRates(currency any, amount any) (*map[string]float32, int, error){
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
		return nil, *statusCode, nil
	}
	rateResponse := map[string]float32{}
	formattedRates := chimoneyResponse.Data.FormatAllRates()
	if currency != "" {
		for key, rate := range formattedRates {
			if strings.Contains(key, currency.(string)) {
				rateResponse = map[string]float32{
					"rate": rate,
				}
				break
			}
		}
		if rateResponse["rate"] == 0 {
			return nil, *statusCode, fmt.Errorf("Currency %s is not supported", currency)
		}

		if amount != "" {
				rateResponse["convertToUSD"] = (rateResponse["rate"] * float32(amount.(uint64))) / chimoneyResponse.Data.USDNGN
				rateResponse["convertedValue"] = rateResponse["rate"] * float32(amount.(uint64))
		}
	}else {
		rateResponse = formattedRates
	}
	return &rateResponse, *statusCode, nil
}

func (chimoneyPP *ChimoneyPaymentProcessor)GetSupportedInternationalBanks(countryCode string) (*[]entities.Bank,  int, error) {
	response, statusCode, err := chimoneyPP.Network.Get(fmt.Sprintf("/info/country-banks?countryCode=%s", countryCode), &map[string]string{
		"X-API-KEY": chimoneyPP.AuthToken,
		"Content-Type": "application/json",
	}, nil)

	var chimoneyResponse ChimoneySupportedBanksDTO
	json.Unmarshal(*response, &chimoneyResponse)
	if err != nil {
		logger.Error(errors.New("an error occured while fetching supported banks on chimoney"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, *statusCode, errors.New("an error occured while fetching supported banks on chimoney")
	}
	if *statusCode != 200 {
		err = errors.New("failed to fetch supported banks")
		logger.Error(err, logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: chimoneyResponse,
		})
		return &chimoneyResponse.Data, *statusCode, nil
	}
	return &chimoneyResponse.Data, *statusCode, nil
}


func (chimoneyPP *ChimoneyPaymentProcessor)InitiateInternationalPayment(payload *InternationalPaymentRequestPayload) (*InternationalPaymentRequestResponseDataPayload,  int, error) {
	response, statusCode, err := chimoneyPP.Network.Post("/payouts/bank", &map[string]string{
		"X-API-KEY": chimoneyPP.AuthToken,
	}, map[string]interface{}{
		"banks": []InternationalPaymentRequestPayload{
			*payload,
		},
	}, nil)

	var chimoneyResponse InternationalPaymentRequestResponsePayload
	json.Unmarshal(*response, &chimoneyResponse)
	if err != nil {
		logger.Error(errors.New("an error occured while initiating international payment on chimoney"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, *statusCode, errors.New("an error occured while initiating international payment on chimoney")
	}
	if *statusCode != 200 {
		err = errors.New("failed to initiate international payment")
		logger.Error(err, logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: payload,
		}, logger.LoggerOptions{
			Key: "response",
			Data: chimoneyResponse,
		})
		return &chimoneyResponse.Data, *statusCode, nil
	}
	return &chimoneyResponse.Data, *statusCode, nil
}