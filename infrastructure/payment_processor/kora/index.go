package kora_local_payment_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/network"
	"usepolymer.co/infrastructure/payment_processor/types"
)

var LocalPaymentProcessor *KoraPaymentProcessor = &KoraPaymentProcessor{}

type KoraPaymentProcessor struct {
	Network   *network.NetworkController
	SecretKey string
}

func (fpp *KoraPaymentProcessor) InitialisePaymentProcessor() {
	LocalPaymentProcessor.Network = &network.NetworkController{
		BaseUrl: os.Getenv("KORA_BASE_URL"),
	}
	LocalPaymentProcessor.SecretKey = os.Getenv("KORA_SECRET_KEY")
}

func (fpp *KoraPaymentProcessor) GenerateDVA(payload *CreateVirtualAccountPayload) (*types.VirtualAccountPayload, *int, error) {
logger.Info("here", logger.LoggerOptions{
	Key: "here",
	Data: payload.KYC,
})
	response, statusCode, err := fpp.Network.Post("/merchant/api/v1/virtual-bank-account", &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", fpp.SecretKey),
		"Content-Type":  "application/json",
	}, payload, nil, false, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while generating account number on kora"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, statusCode, errors.New("an error occured while generating account number")
	}

	var koraResponse types.CreateVirtualAccountResponse
	json.Unmarshal(*response, &koraResponse)

	if *statusCode != 200 {
		err = errors.New("failed to generate account number")
		var errRes any
		json.Unmarshal(*response, &errRes)
		payload.KYC.BVN = ""
		payload.KYC.NIN = ""
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "body",
			Data: errRes,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: payload,
		})
		return nil, statusCode, err
	}
	logger.Info(fmt.Sprintf("generated dva for %s", payload.Reference))
	return &koraResponse.Data, statusCode, nil
}

func (fpp *KoraPaymentProcessor) InitiateLocalTransfer(payload *types.InitiateLocalTransferPayload) (*types.InitiateLocalTransferDataField, *int, error) {
	response, statusCode, err := fpp.Network.Post("/merchant/api/v1/transactions/disburse", &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", fpp.SecretKey),
		"Content-Type":  "application/json",
	}, any(payload), nil, false, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while initiating local transfer on kora"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, statusCode, errors.New("an error occured while generating local transfer")
	}

	var koraResponse types.InitiateTransferPayloadResponse
	json.Unmarshal(*response, &koraResponse)
	if *statusCode != 200 {
		err = errors.New("failed to initiate local transfer")
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "body",
			Data: koraResponse,
		})
		return nil, statusCode, err
	}
	return &koraResponse.Data, statusCode, nil
}

func (fpp *KoraPaymentProcessor) InitiateMobileMoneyTransfer(payload *types.InitiateLocalTransferPayload) (*types.InitiateLocalTransferDataField, *int, error) {
	response, statusCode, err := fpp.Network.Post("/transfers", &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", fpp.SecretKey),
		"Content-Type":  "application/json",
	}, any(payload).(map[string]any), nil, false, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while initiating mobile money transfer on flutterwave"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, statusCode, errors.New("an error occured while generating mobile money transfer")
	}

	var koraResponse types.InitiateTransferPayloadResponse
	json.Unmarshal(*response, &koraResponse)
	if *statusCode != 200 {
		err = errors.New("failed to initiate mobile money transfer")
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "body",
			Data: koraResponse,
		})
		return nil, statusCode, err
	}
	return &koraResponse.Data, statusCode, nil
}

func (fpp *KoraPaymentProcessor) NameVerification(accountNumber string, bankCode string) (*types.NameVerificationResponseField, *int, error) {
	response, statusCode, err := fpp.Network.Post("/merchant/api/v1/misc/banks/resolve", &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", fpp.SecretKey),
		"Content-Type":  "application/json",
	}, map[string]any{
		"account": accountNumber,
		"bank":    bankCode,
	}, nil, false, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while initiating local transfer on kora"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, statusCode, errors.New("an error occured while generating local transfer")
	}

	var koraResponse types.NameVerificationResponseDTO
	json.Unmarshal(*response, &koraResponse)
	if *statusCode != 200 {
		err = fmt.Errorf("account number %s could not be verified at the specified bank", accountNumber)
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "body",
			Data: koraResponse,
		})
		return nil, statusCode, err
	}
	return &koraResponse.Data, statusCode, nil
}
