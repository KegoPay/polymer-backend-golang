package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "kego.com/application/appErrors"
	bankssupported "kego.com/application/banksSupported"
	"kego.com/application/constants"
	"kego.com/application/controllers/dto"
	"kego.com/application/interfaces"
	"kego.com/application/repository"
	"kego.com/application/services"
	"kego.com/application/utils"
	"kego.com/entities"
	currencyformatter "kego.com/infrastructure/currency_formatter"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
	pushnotification "kego.com/infrastructure/messaging/push_notifications"
	international_payment_processor "kego.com/infrastructure/payment_processor/chimoney"
	server_response "kego.com/infrastructure/serverResponse"
)

func InitiateBusinessInternationalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]){
	rates, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetExchangeRates(&ctx.Body.Amount)
	if err != nil {
		apperrors.ExternalDependencyError(ctx.Ctx, "chimoney", fmt.Sprintf("%d", statusCode), err)
		return
	}
	if statusCode != 200 {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("chimoney international payment returned with status code %d", statusCode))
		return
	}
	var fxRate *entities.ParsedExchangeRates = nil
	var usdRate float32 = 0
	for key, rate := range *rates {
		if strings.Contains(key, *ctx.Body.DestinationCountryCode) {
			fxRate = &rate
		}
		if strings.Contains(key, "(US)") {
			usdRate = rate.NGNRate
		}
		if usdRate != 0 && fxRate != nil {
			break
		}
	}
	var USDNGN = usdRate / float32(ctx.Body.Amount)
	if fxRate == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("country %s is not supported", *ctx.Body.DestinationCountryCode))
		return
	}
	if	fxRate.USDRate < constants.MINIMUM_INTERNATIONAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send less than $1 (₦%s)", currencyformatter.HumanReadableFloat32Currency(USDNGN)), nil)
		return
	}
	if fxRate.USDRate > constants.MAXIMUM_INTERNATIONAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send more than $20,000 (₦%s) at a time", currencyformatter.HumanReadableFloat32Currency(USDNGN)), nil)
		return
	}
	internationalProcessorFee, transactionFee  := utils.GetInternationalTransactionFee(fxRate.NGNRate)
	totalAmount := internationalProcessorFee + transactionFee + fxRate.NGNRate
	destinationCountry := utils.CountryCodeToCountryName(*ctx.Body.DestinationCountryCode)
	if os.Getenv("GIN_MODE") != "release" {
		destinationCountry = utils.CountryCodeToCountryName("NG")
	}
	if destinationCountry == "" {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("unsupported country code used %s", *ctx.Body.DestinationCountryCode))
		return
	}
	trxRef := utils.GenerateUUIDString()
	businessID := ctx.GetStringParameter("businessID") 
	wallet , err := services.InitiatePreAuth(ctx.Ctx, &businessID, ctx.GetStringContextData("UserID"), totalAmount, ctx.Body.Pin)
	if err != nil {
		return
	}
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(totalAmount), entities.ChimoneyDebitInternational, trxRef)
	if err != nil {
		return
	}
	response := services.InitiateInternationalPayment(ctx.Ctx, &international_payment_processor.InternationalPaymentRequestPayload{
		DestinationCountry: destinationCountry,
		AccountNumber: ctx.Body.AccountNumber,
		BankCode: ctx.Body.BankCode,
		ValueInUSD: fxRate.USDRate,
		Reference: trxRef,
	})
	if response == nil {
		services.ReverseLockFunds(ctx.Ctx, wallet.ID, trxRef)
		return
	}
	services.RemoveLockFunds(ctx.Ctx, wallet.ID, trxRef)
	transaction := entities.Transaction{
		TransactionReference: response.Chimoneys[0].ChiRef,
		MetaData: response.Chimoneys[0],
		AmountInUSD: utils.GetUInt64Pointer(utils.Float32ToUint64Currency(response.Chimoneys[0].ValueInUSD)),
		AmountInNGN: utils.Float32ToUint64Currency(totalAmount),
		Fee: utils.Float32ToUint64Currency(transactionFee),
		ProcessorFeeCurrency: "USD",
		ProcessorFee: utils.Float32ToUint64Currency(internationalProcessorFee),
		Amount: ctx.Body.Amount,
		Currency: utils.CurrencyCodeToCurrencySymbol(*ctx.Body.DestinationCountryCode),
		WalletID: wallet.ID,
		UserID: wallet.UserID,
		BusinessID: wallet.BusinessID,
		Description: func () string {
			if	ctx.Body.Description == nil {
				des := fmt.Sprintf("International transfer from %s %s to %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName"), *ctx.Body.FullName)
				return des
			}
			return *ctx.Body.Description
		}(),
		Location: entities.Location{
			IPAddress: ctx.Body.IPAddress,
		},
		Intent: entities.ChimoneyDebitInternational,
		DeviceInfo: &entities.DeviceInfo{
			IPAddress: ctx.Body.IPAddress,
			DeviceID: utils.GetStringPointer(ctx.GetStringContextData("DeviceID")),
			UserAgent: utils.GetStringPointer(ctx.GetStringContextData("UserAgent")),
		},
		Sender: entities.TransactionSender{
			BusinessName: wallet.BusinessName,
			FullName: fmt.Sprintf("%s %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName")),
			Email: utils.GetStringPointer(ctx.GetStringContextData("Email")),
		},
		Recepient: entities.TransactionRecepient{
			FullName: *ctx.Body.FullName,
			BankCode: &ctx.Body.BankCode,
			AccountNumber: ctx.Body.AccountNumber,
			BranchCode: ctx.Body.BranchCode,
			Country: ctx.Body.DestinationCountryCode,
		},
	}
	trxRepository := repository.TransactionRepo()
	trx, err := trxRepository.CreateOne(context.TODO(), transaction)
	if err != nil {
		logger.Error(errors.New("error creating transaction for international transfer"), logger.LoggerOptions{
			Key: "payload",
			Data: transaction,
		})
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}

	if ctx.GetBoolContextData("PushNotifOptions") {
		pushnotification.PushNotificationService.PushOne(ctx.GetStringContextData("PushNotificationToken"), "Your payment is on its way! 🚀",
			fmt.Sprintf("Your payment of %s%d to %s in %s is currently being processed.", utils.CurrencyCodeToCurrencySymbol(transaction.Currency), transaction.Amount, transaction.Recepient.FullName, utils.CountryCodeToCountryName(*transaction.Recepient.Country)))
	}

	if ctx.GetBoolContextData("EmailOptions") {
		emails.EmailService.SendEmail(ctx.GetStringContextData("Email"), "Your payment is on its way! 🚀", "payment_sent", map[string]any{
			"FIRSTNAME": transaction.Sender.FullName,
			"CURRENCY_CODE": utils.CurrencyCodeToCurrencySymbol(*ctx.Body.DestinationCountryCode),
			"AMOUNT": utils.UInt64ToFloat32Currency(ctx.Body.Amount),
			"RECEPIENT_NAME": transaction.Recepient.FullName,
			"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName(*transaction.Recepient.Country),
		})
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "Your payment is on its way! 🚀", trx, nil)
}

func InitiatePersonalInternationalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]){
	rates, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetExchangeRates(&ctx.Body.Amount)
	if err != nil {
		apperrors.ExternalDependencyError(ctx.Ctx, "chimoney", fmt.Sprintf("%d", statusCode), err)
		return
	}
	if statusCode != 200 {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("chimoney international payment returned with status code %d", statusCode))
		return
	}
	var fxRate *entities.ParsedExchangeRates = nil
	var usdRate float32 = 0
	for key, rate := range *rates {
		if strings.Contains(key, *ctx.Body.DestinationCountryCode) {
			fxRate = &rate
		}
		if strings.Contains(key, "(US)") {
			usdRate = rate.NGNRate
		}
		if usdRate != 0 && fxRate != nil {
			break
		}
	}
	var USDNGN = usdRate / float32(ctx.Body.Amount)
	if fxRate == nil {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("country %s is not supported", *ctx.Body.DestinationCountryCode))
		return
	}
	if	fxRate.USDRate < constants.MINIMUM_INTERNATIONAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send less than $1 (₦%s)", currencyformatter.HumanReadableFloat32Currency(USDNGN)), nil)
		return
	}
	if fxRate.USDRate > constants.MAXIMUM_INTERNATIONAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send more than $20,000 (₦%s) at a time", currencyformatter.HumanReadableFloat32Currency(USDNGN)), nil)
		return
	}
	internationalProcessorFee, transactionFee  := utils.GetInternationalTransactionFee(fxRate.NGNRate)
	totalAmount := internationalProcessorFee + transactionFee + fxRate.NGNRate
	destinationCountry := utils.CountryCodeToCountryName(*ctx.Body.DestinationCountryCode)
	if os.Getenv("GIN_MODE") != "release" {
		destinationCountry = utils.CountryCodeToCountryName("NG")
	}
	if destinationCountry == "" {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("unsupported country code used %s", *ctx.Body.DestinationCountryCode))
		return
	}
	trxRef := utils.GenerateUUIDString()
	wallet , err := services.InitiatePreAuth(ctx.Ctx, nil, ctx.GetStringContextData("UserID"), utils.Float32ToUint64Currency(totalAmount), ctx.Body.Pin)
	if err != nil {
		return
	}
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(totalAmount), entities.ChimoneyDebitInternational, trxRef)
	if err != nil {
		return
	}
	response := services.InitiateInternationalPayment(ctx.Ctx, &international_payment_processor.InternationalPaymentRequestPayload{
		DestinationCountry: destinationCountry,
		AccountNumber: ctx.Body.AccountNumber,
		BankCode: ctx.Body.BankCode,
		ValueInUSD: fxRate.USDRate,
		Reference: trxRef,
	})
	if response == nil {
		services.ReverseLockFunds(ctx.Ctx, wallet.ID, trxRef)
		return
	}
	services.RemoveLockFunds(ctx.Ctx, wallet.ID, trxRef)
	transaction := entities.Transaction{
		TransactionReference: response.Chimoneys[0].ChiRef,
		MetaData: response.Chimoneys[0],
		AmountInUSD: utils.GetUInt64Pointer(utils.Float32ToUint64Currency(response.Chimoneys[0].ValueInUSD)),
		AmountInNGN: utils.Float32ToUint64Currency(totalAmount),
		Fee: utils.Float32ToUint64Currency(transactionFee),
		ProcessorFeeCurrency: "USD",
		ProcessorFee: utils.Float32ToUint64Currency(internationalProcessorFee),
		Amount: ctx.Body.Amount,
		Currency: utils.CurrencyCodeToCurrencySymbol(*ctx.Body.DestinationCountryCode),
		WalletID: wallet.ID,
		UserID: wallet.UserID,
		Description: func () string {
			if	ctx.Body.Description == nil {
				des := fmt.Sprintf("International transfer from %s %s to %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName"), *ctx.Body.FullName)
				return des
			}
			return *ctx.Body.Description
		}(),
		Location: entities.Location{
			IPAddress: ctx.Body.IPAddress,
		},
		Intent: entities.ChimoneyDebitInternational,
		DeviceInfo: &entities.DeviceInfo{
			IPAddress: ctx.Body.IPAddress,
			DeviceID: utils.GetStringPointer(ctx.GetStringContextData("DeviceID")),
			UserAgent: utils.GetStringPointer(ctx.GetStringContextData("UserAgent")),
		},
		Sender: entities.TransactionSender{
			FullName: fmt.Sprintf("%s %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName")),
			Email: utils.GetStringPointer(ctx.GetStringContextData("Email")),
		},
		Recepient: entities.TransactionRecepient{
			FullName: *ctx.Body.FullName,
			BankCode: &ctx.Body.BankCode,
			AccountNumber: ctx.Body.AccountNumber,
			BranchCode: ctx.Body.BranchCode,
			Country: ctx.Body.DestinationCountryCode,
		},
	}
	trxRepository := repository.TransactionRepo()
	trx, err := trxRepository.CreateOne(context.TODO(), transaction)
	if err != nil {
		logger.Error(errors.New("error creating transaction for international transfer"), logger.LoggerOptions{
			Key: "payload",
			Data: transaction,
		})
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}

	if ctx.GetBoolContextData("PushNotifOptions") {
		pushnotification.PushNotificationService.PushOne(ctx.GetStringContextData("PushNotificationToken"), "Your payment is on its way! 🚀",
			fmt.Sprintf("Your payment of %s%d to %s in %s is currently being processed.", utils.CurrencyCodeToCurrencySymbol(transaction.Currency), transaction.Amount, transaction.Recepient.FullName, utils.CountryCodeToCountryName(*transaction.Recepient.Country)))
	}

	if ctx.GetBoolContextData("EmailOptions") {
		emails.EmailService.SendEmail(ctx.GetStringContextData("Email"), "Your payment is on its way! 🚀", "payment_sent", map[string]any{
			"FIRSTNAME": transaction.Sender.FullName,
			"CURRENCY_CODE": utils.CurrencyCodeToCurrencySymbol(*ctx.Body.DestinationCountryCode),
			"AMOUNT": utils.UInt64ToFloat32Currency(ctx.Body.Amount),
			"RECEPIENT_NAME": transaction.Recepient.FullName,
			"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName(*transaction.Recepient.Country),
		})
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "Your payment is on its way! 🚀", trx, nil)
}

func InitiateBusinessLocalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]){
	// if ctx.Body.Amount < constants.MINIMUM_LOCAL_TRANSFER_LIMIT {
	// 	apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send less than ₦%s", currencyformatter.HumanReadableIntCurrency(constants.MINIMUM_LOCAL_TRANSFER_LIMIT)), nil)
	// 	return
	// }
	// if ctx.Body.Amount >= constants.MAXIMUM_LOCAL_TRANSFER_LIMIT {
	// 	apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send more than ₦%s at a time", currencyformatter.HumanReadableIntCurrency(constants.MAXIMUM_LOCAL_TRANSFER_LIMIT)), nil)
	// 	return
	// }
	// localProcessorFee, polymerFee := utils.GetLocalTransactionFee(ctx.Body.Amount)
	// totalAmount := ctx.Body.Amount + utils.Float32ToUint64Currency(localProcessorFee) + utils.Float32ToUint64Currency(polymerFee)
	// bankName := ""
	// for _, bank := range bankssupported.SupportedLocalBanks {
	// 	if bank.Code == ctx.Body.BankCode {
	// 		bankName =  bank.Name
	// 		break
	// 	}
	// }
	// if bankName == "" {
	// 	apperrors.NotFoundError(ctx.Ctx, "Selected bank is not currently supported")
	// 	return
	// }
	// businessID := ctx.GetStringParameter("businessID") 
	// wallet , err := services.InitiatePreAuth(ctx.Ctx, &businessID, ctx.GetStringContextData("UserID"), totalAmount, ctx.Body.Pin)
	// if err != nil {
	// 	return
	// }
	// narration := func () string {
	// 	if	ctx.Body.Description == nil {
	// 		des := fmt.Sprintf("NGN Transfer from %s %s to %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName"), bankName)
	// 		return des
	// 	}
	// 	return *ctx.Body.Description
	// }()
	// reference := utils.GenerateUUIDString()
	// err = services.LockFunds(ctx.Ctx, wallet, totalAmount, entities.FlutterwaveDebitLocal, reference)
	// if err != nil {
	// 	return
	// }
	// if os.Getenv("GIN_MODE") != "release" {
	// 	reference = fmt.Sprintf("%s_PMCKDU_1", reference)
	// }
	// response := services.InitiateLocalPayment(ctx.Ctx, &types.InitiateLocalTransferPayload{
	// 	AccountNumber: ctx.Body.AccountNumber,
	// 	AccountBank: ctx.Body.BankCode,
	// 	Currency: "NGN",
	// 	Amount: utils.UInt64ToFloat32Currency(ctx.Body.Amount),
	// 	Narration: narration ,
	// 	Reference: reference,
	// 	DebitCurrency: "NGN",
	// 	CallbackURL: os.Getenv("LOCAL_TRANSFER_WEBHOOK_URL"),
	// 	Meta: types.InitiateLocalTransferMeta{
	// 		WalletID: wallet.ID,
	// 		UserID: wallet.UserID,
	// 	},
	// })
	// if response == nil {
	// 	services.ReverseLockFunds(ctx.Ctx, wallet.ID, reference)
	// 	return
	// }
	// transaction := entities.Transaction{
	// 	TransactionReference: reference,
	// 	MetaData: response,
	// 	AmountInNGN: totalAmount,
	// 	Fee: utils.Float32ToUint64Currency(polymerFee),
	// 	ProcessorFeeCurrency: "NGN",
	// 	ProcessorFee: utils.Float32ToUint64Currency(localProcessorFee),
	// 	Amount: totalAmount,
	// 	Currency: "NGN",
	// 	WalletID: wallet.ID,
	// 	UserID: wallet.UserID,
	// 	BusinessID: wallet.BusinessID,
	// 	Description: narration,
	// 	Location: entities.Location{
	// 		IPAddress: ctx.Body.IPAddress,
	// 	},
	// 	Intent: entities.FlutterwaveDebitLocal,
	// 	DeviceInfo: &entities.DeviceInfo{
	// 		IPAddress: ctx.Body.IPAddress,
	// 		DeviceID: utils.GetStringPointer(ctx.GetStringContextData("DeviceID")),
	// 		UserAgent: utils.GetStringPointer(ctx.GetStringContextData("UserAgent")),
	// 	},
	// 	Sender: entities.TransactionSender{
	// 		BusinessName: wallet.BusinessName,
	// 		FullName: fmt.Sprintf("%s %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName")),
	// 		Email: utils.GetStringPointer(ctx.GetStringContextData("Email")),
	// 	},
	// 	Recepient: entities.TransactionRecepient{
	// 		FullName: response.FullName,
	// 		BankCode: &ctx.Body.BankCode,
	// 		AccountNumber: ctx.Body.AccountNumber,
	// 		BranchCode: ctx.Body.BranchCode,
	// 		BankName: &bankName,
	// 		Country: utils.GetStringPointer("Nigeria"),
	// 	},
	// }
	// trxRepository := repository.TransactionRepo()
	// trx, err := trxRepository.CreateOne(context.TODO(), transaction)
	// if err != nil {
	// 	logger.Error(errors.New("error creating transaction for international transfer"), logger.LoggerOptions{
	// 		Key: "payload",
	// 		Data: transaction,
	// 	})
	// 	apperrors.FatalServerError(ctx.Ctx, err)
	// 	return
	// }
	// if ctx.GetBoolContextData("PushNotifOptions") {
	// 	pushnotification.PushNotificationService.PushOne(ctx.GetStringContextData("PushNotificationToken"), "Your payment is on its way! 🚀",
	// 		fmt.Sprintf("Your payment of %s%d to %s in %s is currently being processed.", utils.CurrencyCodeToCurrencySymbol(transaction.Currency), transaction.Amount, transaction.Recepient.FullName, utils.CountryCodeToCountryName(*transaction.Recepient.Country)))
	// }

	// if ctx.GetBoolContextData("EmailOptions") {
	// 	emails.EmailService.SendEmail(ctx.GetStringContextData("Email"), "Your payment is on its way! 🚀", "payment_sent", map[string]any{
	// 		"FIRSTNAME": transaction.Sender.FullName,
	// 		"CURRENCY_CODE": utils.CurrencyCodeToCurrencySymbol("NGN"),
	// 		"AMOUNT": utils.UInt64ToFloat32Currency(ctx.Body.Amount),
	// 		"RECEPIENT_NAME": transaction.Recepient.FullName,
	// 		"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName(*transaction.Recepient.Country),
	// 	})
	// }
	// server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "Your payment is on its way! 🚀", trx, nil)
}

func BusinessLocalPaymentFee(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]){
	if ctx.Body.Amount < 1000 {
		apperrors.ClientError(ctx.Ctx, "You cannot send less than ₦10", nil)
		return
	}
	if ctx.Body.Amount >= 30000000000 {
		apperrors.ClientError(ctx.Ctx, "You cannot send more than ₦300,000,000 at a time", nil)
		return
	}
	localProcessorFee, polymerFee := utils.GetLocalTransactionFee(ctx.Body.Amount)
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "fee calculated", map[string]any{
		"processorFee": utils.Float32ToUint64Currency(localProcessorFee),
		"polymerFee": utils.Float32ToUint64Currency(polymerFee),
	}, nil)
}

func BusinessInternationalPaymentFee(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]){
	if ctx.Body.Amount < 1000 {
		apperrors.ClientError(ctx.Ctx, "You cannot send less than ₦10", nil)
		return
	}
	if ctx.Body.Amount >= 30000000000 {
		apperrors.ClientError(ctx.Ctx, "You cannot send more than ₦300,000,000 at a time", nil)
		return
	}
	internationalProcessorFee, polymerFee := utils.GetInternationalTransactionFee(utils.UInt64ToFloat32Currency(ctx.Body.Amount))
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "fee calculated", map[string]any{
		"processorFee": utils.Float32ToUint64Currency(internationalProcessorFee),
		"polymerFee": utils.Float32ToUint64Currency(polymerFee),
	}, nil)
}

func VerifyLocalAccountName(ctx *interfaces.ApplicationContext[dto.NameVerificationDTO]){
	bankCode := ""
	for _, bank := range bankssupported.SupportedLocalBanks {
		if bank.Name == ctx.Body.BankName {
			bankCode = bank.Code
			break
		}
	}
	if bankCode  == "" {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("%s is not a supported bank on our platform yet.", ctx.Body.BankName))
		return
	}
	name := services.NameVerification(ctx.Ctx, ctx.Body.AccountNumber, bankCode)
	if name == nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "name verification complete", map[string]string{
		"name": *name,
	}, nil)
}

func FetchPastTransactions(ctx *interfaces.ApplicationContext[any]){
	transactionsRepo := repository.TransactionRepo()
	transactions, err := transactionsRepo.FindMany(map[string]interface{}{
		"userID": ctx.GetStringContextData("UserID"),
	}, &options.FindOptions{
		Limit: utils.GetInt64Pointer(15),
		Sort: map[string]any{
			"createdAt": -1,
		},
	}, options.Find().SetProjection(map[string]any{
		"transactionReference": 1,
		"amount": 1,
		"amountInNGN": 1,
		"fee": 1,
		"description": 1,
		"amountInUSD": 1,
		"transactionRecepient": 1,
		"currency": 1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err)
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "transctions fetched", transactions, nil)
}