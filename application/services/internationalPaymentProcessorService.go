package services

import (
	"errors"
	"fmt"

	apperrors "kego.com/application/appErrors"
	"kego.com/entities"
	international_payment_processor "kego.com/infrastructure/payment_processor/chimoney"
)


func FetchInternationalBanks(ctx any, countryCode string) *[]entities.Bank {
	response, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetSupportedInternationalBanks(countryCode)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, "chimoney", fmt.Sprintf("%d", statusCode), err)
		return nil
	}
	if len(*response) == 0 {
		apperrors.ClientError(ctx, fmt.Sprintf("No banks found for %s", countryCode), nil)
		return nil
	}
	if statusCode >= 400 {
		apperrors.ClientError(ctx, fmt.Sprintf("No banks found for %s", countryCode), nil)
		return nil
	}
	return response
}

func InitiateInternationalPayment(ctx any, payload *international_payment_processor.InternationalPaymentRequestPayload) *international_payment_processor.InternationalPaymentRequestResponseDataPayload {
	response, statusCode, err :=  international_payment_processor.InternationalPaymentProcessor.InitiateInternationalPayment(payload)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, "chimoney", fmt.Sprintf("%d", statusCode), err)
		return nil
	}
	if response == nil {
		apperrors.UnknownError(ctx, errors.New("chimoney initiate international payment returned a nil response"))
		return nil
	}
	if statusCode >= 400 {
		apperrors.UnknownError(ctx, fmt.Errorf("chimoney initiate international payment returned with status code %s", statusCode))
		return nil
	}
	return response
}