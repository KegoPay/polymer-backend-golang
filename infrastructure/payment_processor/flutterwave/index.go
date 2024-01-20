package flutterwave_local_payment_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/network"
	"kego.com/infrastructure/payment_processor/types"
)


var LocalPaymentProcessor *FlutterwavePaymentProcessor = &FlutterwavePaymentProcessor{}


type FlutterwavePaymentProcessor struct {
	Network *network.NetworkController
	AuthToken string
}

func (fpp *FlutterwavePaymentProcessor) InitialisePaymentProcessor() {
	LocalPaymentProcessor.Network =  &network.NetworkController{
		BaseUrl: os.Getenv("FLUTTERWAVE_BASE_URL"),
	}
	LocalPaymentProcessor.AuthToken =  os.Getenv("FLUTTERWAVE_ACCESS_TOKEN")
}

func (fpp *FlutterwavePaymentProcessor) GenerateDVA(payload *types.CreateVirtualAccountPayload) (*types.VirtualAccountPayload, *int,  error) {
	response, statusCode, err := fpp.Network.Post("/virtual-account-numbers", &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", fpp.AuthToken),
		"Content-Type": "application/json",
	}, payload, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while generating account number on flutterwave"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, statusCode, errors.New("an error occured while generating account number")
	}

	var flwResponse types.CreateVirtualAccountResponse
	json.Unmarshal(*response, &flwResponse)

	if *statusCode != 200 {
		err = errors.New("failed to generate account number")
		logger.Error(err, logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: flwResponse,
		})
		return nil, statusCode, err
	}
	logger.Info(fmt.Sprintf("generated dva for %s", payload.TransactionReference))
	return &flwResponse.Data, statusCode, nil
}


func (fpp *FlutterwavePaymentProcessor) InitiateLocalTransfer(payload *types.InitiateLocalTransferPayload) (*types.InitiateLocalTransferDataField, *int, error) {
	response, statusCode, err := fpp.Network.Post("/transfers", &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", fpp.AuthToken),
		"Content-Type": "application/json",
	}, payload, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while initiating local transfer on flutterwave"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, statusCode,  errors.New("an error occured while generating local transfer")
	}

	var flwResponse types.InitiateLocalTransferPayloadResponse
	json.Unmarshal(*response, &flwResponse)
	if *statusCode != 200 {
		err = errors.New("failed to initiate local transfer")
		logger.Error(err, logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: flwResponse,
		})
		return nil, statusCode, err
	}
	return &flwResponse.Data, statusCode, nil
}

func (fpp *FlutterwavePaymentProcessor) NameVerification(accountNumber string, bankCode string) (*types.NameVerificationResponseField, *int, error) {
	response, statusCode, err := fpp.Network.Post("/accounts/resolve", &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", fpp.AuthToken),
		"Content-Type": "application/json",
	}, map[string]string{
		"account_number": accountNumber,
		"account_bank": bankCode,
	}, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while initiating local transfer on flutterwave"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, statusCode,  errors.New("an error occured while generating local transfer")
	}

	var flwResponse types.NameVerificationResponseDTO
	json.Unmarshal(*response, &flwResponse)
	if *statusCode != 200 {
		err = fmt.Errorf("Account number %s could not be verified at the specified bank", accountNumber)
		logger.Error(err, logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: flwResponse,
		})
		return nil, statusCode, err
	}
	return &flwResponse.Data, statusCode, nil
}
