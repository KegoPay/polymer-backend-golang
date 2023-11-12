package identityverification

import (
	"encoding/json"
	"errors"
	"fmt"

	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/network"
)

type PaystackVerification struct {
	Network *network.NetworkController
	AuthToken string
}

func (pv *PaystackVerification) CreateAndVerifyUser(payload CustomerPayload, account AccountPayload) (customerID *int, customerCode *string, failureReason string) {
	customer, msg := pv.createCustomer(payload)
	if msg != "" {
		return &customer.Data.ID, &customer.Data.CustomerCode, msg
	}
	msg = pv.VerifyUser(&CustomerVerificationPayload{
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Country: "NG",
		Type: "bank_account",
		BVN: account.BVN,
		BankCode: account.BankCode,
		AccountNumber: account.AccountNumber,
	}, customer.Data.CustomerCode)
	return &customer.Data.ID, &customer.Data.CustomerCode, msg
}

func (pv *PaystackVerification) createCustomer(payload CustomerPayload) (*CustomerCreationDTO, string) {
	response, statusCode, err := pv.Network.Post("/customer", &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", pv.AuthToken),
		"Content-Type": "application/json",
	}, payload, nil)
	var paystackResponse CustomerCreationDTO
	json.Unmarshal(*response, &paystackResponse)
	if err != nil {
		logger.Error(errors.New("error creating customer on paystack"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, "something went wrong with our verification service"
	}
	if *statusCode != 200 {
		logger.Error(errors.New("failed to create customer on paystack"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: paystackResponse,
		})
		return nil, paystackResponse.Message
	}
	return &paystackResponse, ""
}
func (pv *PaystackVerification) VerifyUser(payload *CustomerVerificationPayload, customerCode string) string {
	response, statusCode, err := pv.Network.Post(fmt.Sprintf("/customer/%s/identification", customerCode), &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", pv.AuthToken),
		"Content-Type": "application/json",
	}, payload, nil)
	var paystackResponse CustomerCreationDTO
	json.Unmarshal(*response, &paystackResponse)
	if err != nil {
		logger.Error(errors.New("error verifying customer on paystack"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return "error verifying provided information"
	}
	if *statusCode != 202 {
		logger.Error(errors.New("failed to verify customer on paystack"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: paystackResponse,
		})
		return paystackResponse.Message
	}
	return ""
}