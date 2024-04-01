package services

import (
	"errors"
	"fmt"
	"os"

	apperrors "kego.com/application/appErrors"
	"kego.com/entities"
	"kego.com/infrastructure/logger"
	paymentprocessor "kego.com/infrastructure/payment_processor"
	international_payment_processor "kego.com/infrastructure/payment_processor/chimoney"
	"kego.com/infrastructure/payment_processor/types"
)


func FetchInternationalBanks(ctx any, countryCode string, device_id *string) *[]entities.Bank {
	response, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetSupportedInternationalBanks(countryCode)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, "chimoney", fmt.Sprintf("%d", statusCode), err, device_id)
		return nil
	}
	if len(*response) == 0 {
		apperrors.ClientError(ctx, fmt.Sprintf("No banks found for %s", countryCode), nil, nil, device_id)
		return nil
	}
	if statusCode >= 400 {
		apperrors.ClientError(ctx, fmt.Sprintf("No banks found for %s", countryCode), nil, nil, device_id)
		return nil
	}
	return response
}

func InitiateInternationalPayment(ctx any, payload *international_payment_processor.InternationalPaymentRequestPayload, device_id *string) *international_payment_processor.InternationalPaymentRequestResponseDataPayload {
	response, statusCode, err :=  international_payment_processor.InternationalPaymentProcessor.InitiateInternationalPayment(payload)
	if err != nil {
		logger.Error(errors.New("unexpected error from chimoney"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "status code",
			Data: statusCode,
		}, logger.LoggerOptions{
			Key: "response",
			Data: response,
		})
		apperrors.ExternalDependencyError(ctx, "chimoney", fmt.Sprintf("%d", statusCode), err, device_id)
		return nil
	}
	if response == nil {
		apperrors.UnknownError(ctx, errors.New("chimoney initiate international payment returned a nil response"), device_id)
		return nil
	}
	if statusCode >= 400 {
		apperrors.UnknownError(ctx, fmt.Errorf("chimoney initiate international payment returned with status code %d", statusCode), device_id)
		return nil
	}
	return response
}


func InitiateMobileMoneyPayment(ctx any, payload *types.InitiateLocalTransferPayload, device_id *string) *types.InitiateLocalTransferDataField {
	response, statusCode, err :=  paymentprocessor.LocalPaymentProcessor.InitiateMobileMoneyTransfer(payload)
	if err != nil {
		logger.Error(fmt.Errorf("unexpected error from %s", os.Getenv("LOCAL_PAYMENT_PROCESSOR")), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "status code",
			Data: statusCode,
		}, logger.LoggerOptions{
			Key: "response",
			Data: response,
		})
		apperrors.ExternalDependencyError(ctx, os.Getenv("LOCAL_PAYMENT_PROCESSOR"), fmt.Sprintf("%d", statusCode), err, device_id)
		return nil
	}
	if response == nil {
		logger.Error(fmt.Errorf("nil response from %s",  os.Getenv("LOCAL_PAYMENT_PROCESSOR")), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "status code",
			Data: statusCode,
		})
		apperrors.UnknownError(ctx, fmt.Errorf("nil response from %s initate local payment staus code", os.Getenv("LOCAL_PAYMENT_PROCESSOR")), device_id)
		return nil
	}
	if *statusCode >= 400 {
		logger.Error(errors.New("flutterwave failed to initate local transfer"), logger.LoggerOptions{
			Key: "response",
			Data: response,
		})
		apperrors.UnknownError(ctx, fmt.Errorf("%s initiate local payment returned status code %d", os.Getenv("LOCAL_PAYMENT_PROCESSOR"), *statusCode), device_id)
		return nil
	}
	return response
}
