package services

import (
	"errors"
	"fmt"
	"os"

	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/infrastructure/logger"
	paymentprocessor "usepolymer.co/infrastructure/payment_processor"
	kora_local_payment_processor "usepolymer.co/infrastructure/payment_processor/kora"
	"usepolymer.co/infrastructure/payment_processor/types"
)

func NameVerification(ctx any, accountNumber string, bankCode string, device_id *string) *string {
	response, statusCode, err := paymentprocessor.LocalPaymentProcessor.NameVerification(accountNumber, bankCode)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, os.Getenv("LOCAL_PAYMENT_PROCESSOR"), fmt.Sprintf("%d", statusCode), err, device_id)
		return nil
	}
	return &response.AccountName
}

func InitiateLocalPayment(ctx any, payload *types.InitiateLocalTransferPayload, device_id *string) *types.InitiateLocalTransferDataField {
	payload.Destination.Type = "bank_account"
	payload.Destination.Currency = "NGN"
	response, statusCode, err := paymentprocessor.LocalPaymentProcessor.InitiateLocalTransfer(payload)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, os.Getenv("LOCAL_PAYMENT_PROCESSOR"), fmt.Sprintf("%d", statusCode), err, device_id)
		return nil
	}
	if response == nil {
		apperrors.UnknownError(ctx, fmt.Errorf("nil response from %s initate local payment staus code - %d", os.Getenv("LOCAL_PAYMENT_PROCESSOR"), *statusCode), device_id)
		return nil
	}
	if *statusCode >= 400 {
		logger.Error(errors.New("flutterwave failed to initate local transfer"), logger.LoggerOptions{
			Key:  "response",
			Data: response,
		})
		apperrors.UnknownError(ctx, fmt.Errorf("%s initiate local payment returned status code %d", os.Getenv("LOCAL_PAYMENT_PROCESSOR"), *statusCode), device_id)
		return nil
	}
	return response
}

func GenerateDVA(ctx any, payload *kora_local_payment_processor.CreateVirtualAccountPayload, device_id *string) *types.VirtualAccountPayload {
	response, statusCode, err := paymentprocessor.LocalPaymentProcessor.GenerateDVA(payload)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, "kora", fmt.Sprintf("%d", *statusCode), err, device_id)
		return nil
	}
	return response
}
