package chimoney_international_payment_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/database/repository/cache"
	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/network"
)

var InternationalPaymentProcessor *ChimoneyPaymentProcessor = &ChimoneyPaymentProcessor{}

type ChimoneyPaymentProcessor struct {
	Network   *network.NetworkController
	AuthToken string
}

func (chimoneyPP *ChimoneyPaymentProcessor) InitialisePaymentProcessor() {
	InternationalPaymentProcessor.Network = &network.NetworkController{
		BaseUrl: os.Getenv("CHIMONEY_BASE_URL"),
	}
	InternationalPaymentProcessor.AuthToken = os.Getenv("CHIMONEY_ACCESS_TOKEN")
}

func (chimoneyPP *ChimoneyPaymentProcessor) GetExchangeRates(amount *uint64) (*map[string]entities.ParsedExchangeRates, int, error) {
	// cachedRate := cache.Cache.FindOne("fx_rates")
	// if cachedRate != nil {
	// 	var chimoneyResponse map[string]entities.ParsedExchangeRates
	// 	json.Unmarshal([]byte(*cachedRate), &chimoneyResponse)
	// 	return &chimoneyResponse, 200, nil
	// }
	response, statusCode, err := chimoneyPP.Network.Get("/info/exchange-rates", &map[string]string{
		"X-API-KEY":    chimoneyPP.AuthToken,
		"Content-Type": "application/json",
	}, nil)
	var chimoneyResponse ChimoneyExchangeRateDTO
	json.Unmarshal(*response, &chimoneyResponse)
	if err != nil {
		logger.Error(errors.New("an error occured while fetching exchange rates on chimoney"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, *statusCode, errors.New("an error occured while fetching exchange rates on chimoney")
	}
	if *statusCode != 200 {
		err = errors.New("failed to fetch exchange rates")
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "body",
			Data: chimoneyResponse,
		})
		return nil, *statusCode, nil
	}
	rates := chimoneyResponse.Data.FormatAllRates(amount)
	durationSeconds := chimoneyResponse.ValidTill / 1000
	duration := (time.Duration(durationSeconds) * time.Second) - 1*time.Minute // expire 1 min before just to be safe
	r, _ := json.Marshal(rates)
	cache.Cache.CreateEntry("fx_rates", r, duration)
	return rates, *statusCode, nil
}

func (chimoneyPP *ChimoneyPaymentProcessor) GetSupportedInternationalBanks(countryCode string) (*[]entities.Bank, int, error) {
	response, statusCode, err := chimoneyPP.Network.Get(fmt.Sprintf("/info/country-banks?countryCode=%s", countryCode), &map[string]string{
		"X-API-KEY":    chimoneyPP.AuthToken,
		"Content-Type": "application/json",
	}, nil)

	var chimoneyResponse ChimoneySupportedBanksDTO
	json.Unmarshal(*response, &chimoneyResponse)
	if err != nil {
		logger.Error(errors.New("an error occured while fetching supported banks on chimoney"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, *statusCode, errors.New("an error occured while fetching supported banks on chimoney")
	}
	if *statusCode != 200 {
		err = errors.New("failed to fetch supported banks")
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "body",
			Data: chimoneyResponse,
		})
		return &chimoneyResponse.Data, *statusCode, nil
	}
	return &chimoneyResponse.Data, *statusCode, nil
}

func (chimoneyPP *ChimoneyPaymentProcessor) InitiateInternationalPayment(payload *InternationalPaymentRequestPayload) (*InternationalPaymentRequestResponseDataPayload, int, error) {
	response, statusCode, err := chimoneyPP.Network.Post("/payouts/bank", &map[string]string{
		"X-API-KEY": chimoneyPP.AuthToken,
	}, map[string]interface{}{
		"banks": []InternationalPaymentRequestPayload{
			*payload,
		},
	}, nil, false, nil)

	var chimoneyResponse InternationalPaymentRequestResponsePayload
	json.Unmarshal(*response, &chimoneyResponse)
	if err != nil {
		logger.Error(errors.New("an error occured while initiating international payment on chimoney"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, *statusCode, errors.New("an error occured while initiating international payment on chimoney")
	}
	if *statusCode != 200 {
		err = errors.New("failed to initiate international payment")
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "body",
			Data: payload,
		}, logger.LoggerOptions{
			Key:  "response",
			Data: chimoneyResponse,
		})
		return &chimoneyResponse.Data, *statusCode, nil
	}
	return &chimoneyResponse.Data, *statusCode, nil
}
