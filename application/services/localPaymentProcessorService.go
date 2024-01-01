package services

import (
	"fmt"
	"os"

	apperrors "kego.com/application/appErrors"
	paymentprocessor "kego.com/infrastructure/payment_processor"
	"kego.com/infrastructure/payment_processor/types"
)


func NameVerification(ctx any, accountNumber string, bankCode string) *string {
	response, statusCode, err := paymentprocessor.LocalPaymentProcessor.NameVerification(accountNumber, bankCode)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, os.Getenv("LOCAL_PAYMENT_PROCESSOR"), fmt.Sprintf("%d", statusCode), err)
		return nil
	}
	return &response.AccountName
}

func InitiateLocalPayment(ctx any, payload *types.InitiateLocalTransferPayload) *types.InitiateLocalTransferDataField {
	response, statusCode, err :=  paymentprocessor.LocalPaymentProcessor.InitiateLocalTransfer(payload)
	if err != nil {
		apperrors.ExternalDependencyError(ctx, os.Getenv("LOCAL_PAYMENT_PROCESSOR"), fmt.Sprintf("%d", statusCode), err)
		return nil
	}
	if response == nil {
		apperrors.UnknownError(ctx)
		return nil
	}
	if *statusCode >= 400 {
		apperrors.UnknownError(ctx)
		return nil
	}
	return response
}