package services

import (
	"fmt"

	apperrors "kego.com/application/appErrors"
	payment_processor "kego.com/infrastructure/payment_processor/paystack"
)


func NameVerification(ctx any, accountNumber string, bankCode string) *string {
	response, statusCode, err := payment_processor.LocalPaymentProcessor.NameVerification(accountNumber,bankCode)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, "paystack", fmt.Sprintf("%d", statusCode), err)
		return nil
	}
	if statusCode == 422 {
		apperrors.ClientError(ctx, fmt.Sprintf("Account number %s could not be verified at the specified bank", accountNumber), nil)
		return nil
	}
	return &response.Data.AccountName
}