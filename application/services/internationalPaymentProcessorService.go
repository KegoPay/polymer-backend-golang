package services

import (
	"errors"
	"fmt"

	apperrors "kego.com/application/appErrors"
	"kego.com/entities"
	international_payment_processor "kego.com/infrastructure/payment_processor/chimoney"
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