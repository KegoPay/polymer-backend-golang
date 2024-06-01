package paystack_local_payment_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/network"
)

var LocalPaymentProcessor *PaystackPaymentProcessor

type PaystackPaymentProcessor struct {
	Network   *network.NetworkController
	AuthToken string
}

func (paystackPP *PaystackPaymentProcessor) InitialisePaymentProcessor() {
	LocalPaymentProcessor = &PaystackPaymentProcessor{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("PAYSTACK_BASE_URL"),
		},
		AuthToken: os.Getenv("PAYSTACK_ACCESS_TOKEN"),
	}
}

func (paystackPP *PaystackPaymentProcessor) NameVerification(accountNumber string, bankCode string) (*PaystackNameVerificationResponseDTO, int, error) {
	response, statusCode, err := paystackPP.Network.Get(fmt.Sprintf("/bank/resolve?account_number=%s&bank_code=%s", accountNumber, bankCode), &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", paystackPP.AuthToken),
		"Content-Type":  "application/json",
	}, nil)
	var paystackResponse PaystackNameVerificationResponseDTO
	json.Unmarshal(*response, &paystackResponse)
	if err != nil {
		logger.Error(errors.New("an error occured while verifying account number on paystack"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, *statusCode, errors.New("an error occured while verifying account number")
	}
	if *statusCode != 200 {
		err = errors.New("failed to verify  account name")
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "body",
			Data: paystackResponse,
		})
		return &paystackResponse, *statusCode, nil
	}
	return &paystackResponse, *statusCode, nil
}
