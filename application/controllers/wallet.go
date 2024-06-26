package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/mongo/options"
	apperrors "usepolymer.co/application/appErrors"
	bankssupported "usepolymer.co/application/banksSupported"
	"usepolymer.co/application/constants"
	"usepolymer.co/application/controllers/dto"
	"usepolymer.co/application/interfaces"
	"usepolymer.co/application/repository"
	"usepolymer.co/application/services"
	"usepolymer.co/application/utils"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/background"
	currencyformatter "usepolymer.co/infrastructure/currency_formatter"
	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/messaging/emails"
	pushnotification "usepolymer.co/infrastructure/messaging/push_notifications"
	international_payment_processor "usepolymer.co/infrastructure/payment_processor/chimoney"
	"usepolymer.co/infrastructure/payment_processor/types"
	server_response "usepolymer.co/infrastructure/serverResponse"
)

func InitiateBusinessInternationalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]) {
	rates, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetExchangeRates(&ctx.Body.Amount)
	if err != nil {
		apperrors.ExternalDependencyError(ctx.Ctx, "chimoney", fmt.Sprintf("%d", statusCode), err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if statusCode != 200 {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("chimoney international payment returned with status code %d", statusCode), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	// var fxRate *entities.ParsedExchangeRates = nil
	var fxUSDRate float32
	var fxNGNRate float32
	var usdRate float32 = 0
	for key, rate := range *rates {
		if strings.Contains(key, *ctx.Body.DestinationCountryCode) {
			fxNGNRate = rate.NGNRate
			fxUSDRate = rate.USDRate
		}
		if strings.Contains(key, "(US)") {
			usdRate = rate.NGNRate
		}
		if usdRate != 0 && fxNGNRate != 0 {
			break
		}
	}
	// var USDNGN = usdRate / float32(ctx.Body.Amount)
	if fxNGNRate == 0 {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("country %s is not supported", *ctx.Body.DestinationCountryCode), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if fxUSDRate < constants.MINIMUM_INTERNATIONAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send less than $1 (₦%s)", currencyformatter.HumanReadableFloat32Currency((usdRate/utils.UInt64ToFloat32Currency(ctx.Body.Amount)*constants.MINIMUM_INTERNATIONAL_TRANSFER_LIMIT))), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if fxUSDRate > constants.MAXIMUM_INTERNATIONAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send more than $20,000 (₦%s) at a time", currencyformatter.HumanReadableFloat32Currency((usdRate/utils.UInt64ToFloat32Currency(ctx.Body.Amount)*constants.MAXIMUM_INTERNATIONAL_TRANSFER_LIMIT))), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	internationalProcessorFee, transactionFee, transactionFeeVat := utils.GetInternationalTransactionFee(fxNGNRate)
	totalAmount := internationalProcessorFee + transactionFee + fxNGNRate + transactionFeeVat
	destinationCountry := utils.CountryCodeToCountryName(*ctx.Body.DestinationCountryCode)
	if destinationCountry == "" {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("unsupported country code used %s", *ctx.Body.DestinationCountryCode), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if os.Getenv("ENV") != "production" {
		destinationCountry = "Nigeria"
	}
	trxRef := utils.GenerateUUIDString()
	businessID := ctx.GetStringParameter("businessID")
	wallet, err := services.InitiatePreAuth(ctx.Ctx, &businessID, ctx.GetStringContextData("UserID"), utils.Float32ToUint64Currency(totalAmount, true), ctx.Body.Pin, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}

	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(fxNGNRate, true), entities.InternationalDebit, trxRef, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking amount to be sent"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}
	// lock processor transaction fee
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(internationalProcessorFee, false), entities.InternationalDebitFee, trxRef, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking chimoney transaction fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// lock polymer transaction fee
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(transactionFee, false), entities.PolymerFee, trxRef, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking polymer processing fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// lock polymer vat
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(transactionFeeVat, false), entities.PolymerVAT, trxRef, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking polymer vat fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	response := services.InitiateInternationalPayment(ctx.Ctx, &international_payment_processor.InternationalPaymentRequestPayload{
		DestinationCountry: destinationCountry,
		AccountNumber:      ctx.Body.AccountNumber,
		BankCode:           ctx.Body.BankCode,
		ValueInUSD:         fxUSDRate,
		Reference:          trxRef,
		FullName:           ctx.Body.FullName,
	}, ctx.GetHeader("Polymer-Device-Id"))
	if response == nil {
		services.ReverseLockFunds(wallet.ID, trxRef)
		return
	}
	transaction := entities.Transaction{
		TransactionReference: trxRef,
		MetaData:             response.Chimoneys[0],
		TotalAmountInNGN:     utils.GetUInt64Pointer(utils.Float32ToUint64Currency(totalAmount, true)),
		AmountInUSD:          utils.GetUInt64Pointer(utils.Float32ToUint64Currency(fxUSDRate, false)),
		AmountInNGN:          utils.GetUInt64Pointer(utils.Float32ToUint64Currency(fxNGNRate, true)),
		Fee:                  utils.Float32ToUint64Currency(transactionFee, false),
		Vat:                  utils.Float32ToUint64Currency(transactionFeeVat, false),
		ProcessorFeeCurrency: "USD",
		ProcessorFee:         utils.Float32ToUint64Currency(internationalProcessorFee, false),
		Amount:               ctx.Body.Amount,
		Status:               entities.TransactionPending,
		Currency:             utils.CurrencyCodeToCurrencySymbol(*ctx.Body.DestinationCountryCode),
		WalletID:             wallet.ID,
		UserID:               wallet.UserID,
		BusinessID:           wallet.BusinessID,
		Description: func() string {
			if ctx.Body.Description == nil {
				des := fmt.Sprintf("International transfer from %s %s to %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName"), *ctx.Body.FullName)
				return des
			}
			return *ctx.Body.Description
		}(),
		Location: entities.Location{
			IPAddress: ctx.Body.IPAddress,
			Latitude:  ctx.GetFloat64ContextData("Latitude"),
			Longitude: ctx.GetFloat64ContextData("Longitude"),
		},
		Intent: entities.InternationalDebit,
		DeviceInfo: &entities.DeviceInfo{
			IPAddress: ctx.Body.IPAddress,
			DeviceID:  utils.GetStringPointer(ctx.GetStringContextData("DeviceID")),
			UserAgent: utils.GetStringPointer(ctx.GetStringContextData("UserAgent")),
		},
		Sender: entities.TransactionSender{
			BusinessName: wallet.BusinessName,
			FullName:     fmt.Sprintf("%s %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName")),
			Email:        utils.GetStringPointer(ctx.GetStringContextData("Email")),
		},
		Recepient: entities.TransactionRecepient{
			FullName:      *ctx.Body.FullName,
			BankCode:      &ctx.Body.BankCode,
			AccountNumber: ctx.Body.AccountNumber,
			BranchCode:    ctx.Body.BranchCode,
			Country:       ctx.Body.DestinationCountryCode,
		},
	}
	trxRepository := repository.TransactionRepo()
	trx, err := trxRepository.CreateOne(context.TODO(), transaction)
	if err != nil {
		logger.Error(errors.New("error creating transaction for international transfer"), logger.LoggerOptions{
			Key:  "payload",
			Data: transaction,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	if ctx.GetBoolContextData("PushNotifOptions") {
		pushnotification.PushNotificationService.PushOne(ctx.GetStringContextData("PushNotificationToken"), "Your payment is on its way! 🚀",
			fmt.Sprintf("Your payment of %s%d to %s in %s is currently being processed.", utils.CurrencyCodeToCurrencySymbol(transaction.Currency), transaction.Amount, transaction.Recepient.FullName, utils.CountryCodeToCountryName(*transaction.Recepient.Country)))
	}

	if ctx.GetBoolContextData("EmailOptions") {
		background.Scheduler.Emit("send_email", map[string]any{
			"email":        ctx.GetStringContextData("Email"),
			"subject":      "Your payment is on its way! 🚀",
			"templateName": "payment_sent",
			"opts": map[string]any{
				"FIRSTNAME":         transaction.Sender.FullName,
				"CURRENCY_CODE":     utils.CurrencyCodeToCurrencySymbol(*ctx.Body.DestinationCountryCode),
				"AMOUNT":            utils.UInt64ToFloat32Currency(ctx.Body.Amount),
				"RECEPIENT_NAME":    transaction.Recepient.FullName,
				"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName(*transaction.Recepient.Country),
			},
		})
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "Your payment is on its way! 🚀", trx, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func InitiatePersonalInternationalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]) {
	var wg sync.WaitGroup
	rates, statusCode, err := international_payment_processor.InternationalPaymentProcessor.GetExchangeRates(&ctx.Body.Amount)
	if err != nil {
		apperrors.ExternalDependencyError(ctx.Ctx, "chimoney", fmt.Sprintf("%d", statusCode), err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if statusCode != 200 {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("chimoney international payment returned with status code %d", statusCode), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	// var fxRate *entities.ParsedExchangeRates = nil
	var fxUSDRate float32
	var fxNGNRate float32
	var usdRate float32 = 0
	for key, rate := range *rates {
		if strings.Contains(key, *ctx.Body.DestinationCountryCode) {
			fxNGNRate = rate.NGNRate
			fxUSDRate = rate.USDRate
		}
		if strings.Contains(key, "(US)") {
			usdRate = rate.NGNRate
		}
		if usdRate != 0 && fxNGNRate != 0 {
			break
		}
	}
	userTier := ctx.GetIntContextData("Tier")
	if userTier == 1 && fxNGNRate > utils.UInt64ToFloat32Currency(constants.TIER_ONE_SINGLE_TRANSFER_LIMITS) {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You can only send up to %s at this moment. Upgrade your account by setting your phone number, address and next of kin to send up to %s at once!", currencyformatter.HumanReadableIntCurrency(constants.TIER_ONE_SINGLE_TRANSFER_LIMITS), currencyformatter.HumanReadableIntCurrency(constants.TIER_TWO_SINGLE_TRANSFER_LIMITS)), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if userTier == 2 && fxNGNRate > utils.UInt64ToFloat32Currency(constants.TIER_TWO_SINGLE_TRANSFER_LIMITS) {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You can only send up to %s at this moment. Upgrade your account by setting your NIN to send up to %s at once!", currencyformatter.HumanReadableIntCurrency(constants.TIER_TWO_SINGLE_TRANSFER_LIMITS), currencyformatter.HumanReadableIntCurrency(constants.TIER_THREE_SINGLE_TRANSFER_LIMITS)), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if userTier == 3 && fxNGNRate > utils.UInt64ToFloat32Currency(constants.TIER_THREE_SINGLE_TRANSFER_LIMITS) {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You can only send up to %s at this moment. Contact support to send more than %s at once!", currencyformatter.HumanReadableIntCurrency(constants.TIER_THREE_SINGLE_TRANSFER_LIMITS), currencyformatter.HumanReadableIntCurrency(constants.TIER_THREE_SINGLE_TRANSFER_LIMITS)), nil, &constants.REQUEST_HIGHER_TRANSFER_LIMIT, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	// var USDNGN = usdRate / float32(ctx.Body.Amount)
	if fxNGNRate == 0 {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("country %s is not supported", *ctx.Body.DestinationCountryCode), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if fxUSDRate < constants.MINIMUM_INTERNATIONAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send less than $1 (₦%s)", currencyformatter.HumanReadableFloat32Currency((usdRate/utils.UInt64ToFloat32Currency(ctx.Body.Amount)*constants.MINIMUM_INTERNATIONAL_TRANSFER_LIMIT))), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if fxUSDRate > constants.MAXIMUM_INTERNATIONAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send more than $20,000 (₦%s) at a time", currencyformatter.HumanReadableFloat32Currency((usdRate/utils.UInt64ToFloat32Currency(ctx.Body.Amount)*constants.MAXIMUM_INTERNATIONAL_TRANSFER_LIMIT))), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	internationalProcessorFee, transactionFee, transactionFeeVat := utils.GetInternationalTransactionFee(fxNGNRate)
	totalAmount := internationalProcessorFee + transactionFee + fxNGNRate + transactionFeeVat
	destinationCountry := utils.CountryCodeToCountryName(*ctx.Body.DestinationCountryCode)
	if destinationCountry == "" {
		apperrors.UnknownError(ctx.Ctx, fmt.Errorf("unsupported country code used %s", *ctx.Body.DestinationCountryCode), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if os.Getenv("ENV") != "production" {
		destinationCountry = "Ghana"
	}
	trxRef := utils.GenerateUUIDString()
	wallet, err := services.InitiatePreAuth(ctx.Ctx, nil, ctx.GetStringContextData("UserID"), utils.Float32ToUint64Currency(totalAmount, true), ctx.Body.Pin, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}

	errChan := make(chan error, 4)
	wg.Add(1)
	// lock amount to be sent
	go func() {
		defer wg.Done()
		err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(fxNGNRate, true), entities.InternationalDebit, trxRef, ctx.GetHeader("Polymer-Device-Id"))
		if err != nil {
			logger.Error(errors.New("error locking amount to be sent"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: ctx.Body,
			})
			errChan <- err
			return
		}
	}()
	// lock international transaction fee
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(internationalProcessorFee, false), entities.InternationalDebitFee, trxRef, ctx.GetHeader("Polymer-Device-Id"))
		if err != nil {
			logger.Error(errors.New("error locking chimoney transaction fee"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: ctx.Body,
			})
			errChan <- err
			return
		}
	}()
	// lock polymer transaction fee
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(transactionFee, false), entities.PolymerFee, trxRef, ctx.GetHeader("Polymer-Device-Id"))
		if err != nil {
			logger.Error(errors.New("error locking polymer processing fee"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: ctx.Body,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: ctx.Body,
			})
			errChan <- err
			return
		}
	}()
	// lock polymer vat
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(transactionFeeVat, false), entities.PolymerVAT, trxRef, ctx.GetHeader("Polymer-Device-Id"))
		if err != nil {
			logger.Error(errors.New("error locking polymer vat fee"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			}, logger.LoggerOptions{
				Key:  "payload",
				Data: ctx.Body,
			})
			errChan <- err
			return
		}
	}()
	wg.Wait()

	response := services.InitiateInternationalPayment(ctx.Ctx, &international_payment_processor.InternationalPaymentRequestPayload{
		DestinationCountry: destinationCountry,
		AccountNumber:      ctx.Body.AccountNumber,
		BankCode:           ctx.Body.BankCode,
		ValueInUSD:         fxUSDRate,
		Reference:          trxRef,
		FullName:           ctx.Body.FullName,
	}, ctx.GetHeader("Polymer-Device-Id"))
	if response == nil {
		services.ReverseLockFunds(wallet.ID, trxRef)
		return
	}
	transaction := entities.Transaction{
		TransactionReference: trxRef,
		MetaData:             response.Chimoneys[0],
		TotalAmountInNGN:     utils.GetUInt64Pointer(utils.Float32ToUint64Currency(totalAmount, true)),
		AmountInUSD:          utils.GetUInt64Pointer(utils.Float32ToUint64Currency(fxUSDRate, false)),
		AmountInNGN:          utils.GetUInt64Pointer(utils.Float32ToUint64Currency(fxNGNRate, true)),
		Fee:                  utils.Float32ToUint64Currency(transactionFee, false),
		Vat:                  utils.Float32ToUint64Currency(transactionFeeVat, false),
		ProcessorFeeCurrency: "USD",
		ProcessorFee:         utils.Float32ToUint64Currency(internationalProcessorFee, false),
		Amount:               ctx.Body.Amount,
		Currency:             utils.CurrencyCodeToCurrencySymbol(*ctx.Body.DestinationCountryCode),
		WalletID:             wallet.ID,
		UserID:               wallet.UserID,
		Status:               entities.TransactionPending,
		Description: func() string {
			if ctx.Body.Description == nil {
				des := fmt.Sprintf("International transfer from %s %s to %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName"), *ctx.Body.FullName)
				return des
			}
			return *ctx.Body.Description
		}(),
		Location: entities.Location{
			IPAddress: ctx.Body.IPAddress,
			Latitude:  ctx.GetFloat64ContextData("Latitude"),
			Longitude: ctx.GetFloat64ContextData("Longitude"),
		},
		Intent: entities.InternationalDebit,
		DeviceInfo: &entities.DeviceInfo{
			IPAddress: ctx.Body.IPAddress,
			DeviceID:  utils.GetStringPointer(ctx.GetStringContextData("DeviceID")),
			UserAgent: utils.GetStringPointer(ctx.GetStringContextData("UserAgent")),
		},
		Sender: entities.TransactionSender{
			FullName: fmt.Sprintf("%s %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName")),
			Email:    utils.GetStringPointer(ctx.GetStringContextData("Email")),
		},
		Recepient: entities.TransactionRecepient{
			FullName:      *ctx.Body.FullName,
			BankCode:      &ctx.Body.BankCode,
			AccountNumber: ctx.Body.AccountNumber,
			BranchCode:    ctx.Body.BranchCode,
			Country:       ctx.Body.DestinationCountryCode,
		},
	}
	trxRepository := repository.TransactionRepo()
	trx, err := trxRepository.CreateOne(context.TODO(), transaction)
	if err != nil {
		logger.Error(errors.New("error creating transaction for international transfer"), logger.LoggerOptions{
			Key:  "payload",
			Data: transaction,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}

	if ctx.GetBoolContextData("PushNotifOptions") {
		pushnotification.PushNotificationService.PushOne(ctx.GetStringContextData("PushNotificationToken"), "Your payment is on its way! 🚀",
			fmt.Sprintf("Your payment of %s%d to %s in %s is currently being processed.", utils.CurrencyCodeToCurrencySymbol(transaction.Currency), transaction.Amount, transaction.Recepient.FullName, utils.CountryCodeToCountryName(*transaction.Recepient.Country)))
	}

	if ctx.GetBoolContextData("EmailOptions") {
		background.Scheduler.Emit("send_email", map[string]any{
			"email":        ctx.GetStringContextData("Email"),
			"subject":      "Your payment is on its way! 🚀",
			"templateName": "payment_sent",
			"opts": map[string]any{
				"FIRSTNAME":         transaction.Sender.FullName,
				"CURRENCY_CODE":     utils.CurrencyCodeToCurrencySymbol(*ctx.Body.DestinationCountryCode),
				"AMOUNT":            utils.UInt64ToFloat32Currency(ctx.Body.Amount),
				"RECEPIENT_NAME":    transaction.Recepient.FullName,
				"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName(*transaction.Recepient.Country),
			},
		})
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "Your payment is on its way! 🚀", trx, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func InitiateBusinessLocalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]) {
	if ctx.Body.Amount < constants.MINIMUM_LOCAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send less than ₦%s", currencyformatter.HumanReadableIntCurrency(constants.MINIMUM_LOCAL_TRANSFER_LIMIT)), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if ctx.Body.Amount > constants.MAXIMUM_LOCAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send more than ₦%s at a time", currencyformatter.HumanReadableIntCurrency(constants.MAXIMUM_LOCAL_TRANSFER_LIMIT)), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	localProcessorFee, localProcessorVAT, polymerFee, polymerVAT := utils.GetLocalTransactionFee(ctx.Body.Amount)
	totalAmount := ctx.Body.Amount + utils.Float32ToUint64Currency(localProcessorFee, false) + utils.Float32ToUint64Currency(polymerFee, false) + utils.Float32ToUint64Currency(localProcessorVAT, false) + utils.Float32ToUint64Currency(polymerVAT, false)
	bankName := ""
	for _, bank := range bankssupported.SupportedLocalBanks {
		if bank.Code == ctx.Body.BankCode {
			bankName = bank.Name
			break
		}
	}
	if bankName == "" {
		apperrors.NotFoundError(ctx.Ctx, "Selected bank is not currently supported", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	businessID := ctx.GetStringParameter("businessID")
	wallet, err := services.InitiatePreAuth(ctx.Ctx, &businessID, ctx.GetStringContextData("UserID"), totalAmount, ctx.Body.Pin, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}
	narration := func() string {
		if ctx.Body.Description == nil {
			des := fmt.Sprintf("NGN Transfer from %s %s to %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName"), bankName)
			return des
		}
		return *ctx.Body.Description
	}()
	reference := utils.GenerateUUIDString()

	// lock amount to be sent
	err = services.LockFunds(ctx.Ctx, wallet, ctx.Body.Amount, entities.LocalDebit, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking amount to be sent"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// lock local transaction fee
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(localProcessorFee, false), entities.LocalDebitFee, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking chimoney transaction fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
	}

	// lock local debit vat
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(localProcessorVAT, false), entities.LocalDebitVAT, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking polymer processing fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// lock polymer fee
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(polymerFee, false), entities.PolymerFee, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking polymer vat fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// lock polymer vat
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(polymerVAT, false), entities.PolymerVAT, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking polymer vat fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// err = services.LockFunds(ctx.Ctx, wallet, totalAmount, entities.FlutterwaveDebitLocal, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}
	if os.Getenv("ENV") != "prod" {
		reference = fmt.Sprintf("%s_PMCKDU_1", reference)
	}
	response := services.InitiateLocalPayment(ctx.Ctx, &types.InitiateLocalTransferPayload{
		Reference: reference,
		Destination: types.LocalTransferDestination{
			Amount:    utils.UInt64ToFloat32Currency(ctx.Body.Amount),
			Narration: narration,
			BankAccount: types.LocalTransferDestinationBankAccount{
				Bank:    ctx.Body.BankCode,
				Account: ctx.Body.AccountNumber,
			},
			Customer: types.LocalTransferDestinationCustomer{
				Name: *ctx.Body.FullName,
			},
		},
	}, ctx.GetHeader("Polymer-Device-Id"))
	if response == nil {
		services.ReverseLockFunds(wallet.ID, reference)
		return
	}
	transaction := entities.Transaction{
		TransactionReference: reference,
		MetaData:             response,
		Fee:                  utils.Float32ToUint64Currency(polymerFee, false),
		Vat:                  utils.Float32ToUint64Currency(polymerVAT, false),
		ProcessorFeeCurrency: "NGN",
		ProcessorFee:         utils.Float32ToUint64Currency(localProcessorFee, false),
		ProcessorFeeVAT:      utils.Float32ToUint64Currency(localProcessorVAT, false),
		Amount:               totalAmount,
		Currency:             "NGN",
		WalletID:             wallet.ID,
		UserID:               wallet.UserID,
		BusinessID:           wallet.BusinessID,
		Status:               entities.TransactionPending,
		Description:          narration,
		Location: entities.Location{
			IPAddress: ctx.Body.IPAddress,
		},
		Intent: entities.LocalDebit,
		DeviceInfo: &entities.DeviceInfo{
			IPAddress: ctx.Body.IPAddress,
			DeviceID:  utils.GetStringPointer(ctx.GetStringContextData("DeviceID")),
			UserAgent: utils.GetStringPointer(ctx.GetStringContextData("UserAgent")),
		},
		Sender: entities.TransactionSender{
			BusinessName: wallet.BusinessName,
			FullName:     fmt.Sprintf("%s %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName")),
			Email:        utils.GetStringPointer(ctx.GetStringContextData("Email")),
		},
		Recepient: entities.TransactionRecepient{
			FullName:      response.FullName,
			BankCode:      &ctx.Body.BankCode,
			AccountNumber: ctx.Body.AccountNumber,
			BranchCode:    ctx.Body.BranchCode,
			BankName:      &bankName,
			Country:       utils.GetStringPointer("Nigeria"),
		},
	}
	trxRepository := repository.TransactionRepo()
	trx, err := trxRepository.CreateOne(context.TODO(), transaction)
	if err != nil {
		logger.Error(errors.New("error creating transaction for international transfer"), logger.LoggerOptions{
			Key:  "payload",
			Data: transaction,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if ctx.GetBoolContextData("PushNotifOptions") {
		pushnotification.PushNotificationService.PushOne(ctx.GetStringContextData("PushNotificationToken"), "Your payment is on its way! 🚀",
			fmt.Sprintf("Your payment of %s%d to %s in %s is currently being processed.", utils.CurrencyCodeToCurrencySymbol(transaction.Currency), transaction.Amount, transaction.Recepient.FullName, utils.CountryCodeToCountryName(*transaction.Recepient.Country)))
	}

	if ctx.GetBoolContextData("EmailOptions") {
		emails.EmailService.SendEmail(ctx.GetStringContextData("Email"), "Your payment is on its way! 🚀", "payment_sent", map[string]any{
			"FIRSTNAME":         transaction.Sender.FullName,
			"CURRENCY_CODE":     utils.CurrencyCodeToCurrencySymbol("NGN"),
			"AMOUNT":            utils.UInt64ToFloat32Currency(ctx.Body.Amount),
			"RECEPIENT_NAME":    transaction.Recepient.FullName,
			"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName(*transaction.Recepient.Country),
		})
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "Your payment is on its way! 🚀", trx, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func InitiatePersonalLocalPayment(ctx *interfaces.ApplicationContext[dto.SendPaymentDTO]) {
	userTier := ctx.GetIntContextData("Tier")
	if userTier == 1 && ctx.Body.Amount > constants.TIER_ONE_SINGLE_TRANSFER_LIMITS {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You can only send up to %s at this moment. Upgrade your account by setting your phone number, address and next of kin to send up to %s at once!", currencyformatter.HumanReadableIntCurrency(constants.TIER_ONE_SINGLE_TRANSFER_LIMITS), currencyformatter.HumanReadableIntCurrency(constants.TIER_TWO_SINGLE_TRANSFER_LIMITS)), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if userTier == 2 && ctx.Body.Amount > constants.TIER_TWO_SINGLE_TRANSFER_LIMITS {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You can only send up to %s at this moment. Upgrade your account by setting your NIN to send up to %s at once!", currencyformatter.HumanReadableIntCurrency(constants.TIER_TWO_SINGLE_TRANSFER_LIMITS), currencyformatter.HumanReadableIntCurrency(constants.TIER_THREE_SINGLE_TRANSFER_LIMITS)), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if userTier == 3 && ctx.Body.Amount > constants.TIER_THREE_SINGLE_TRANSFER_LIMITS {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You can only send up to %s at this moment. Contact support to send more than %s at once!", currencyformatter.HumanReadableIntCurrency(constants.TIER_THREE_SINGLE_TRANSFER_LIMITS), currencyformatter.HumanReadableIntCurrency(constants.TIER_THREE_SINGLE_TRANSFER_LIMITS)), nil, &constants.REQUEST_HIGHER_TRANSFER_LIMIT, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if ctx.Body.Amount < constants.MINIMUM_LOCAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send less than ₦%s", currencyformatter.HumanReadableIntCurrency(constants.MINIMUM_LOCAL_TRANSFER_LIMIT)), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if ctx.Body.Amount > constants.MAXIMUM_LOCAL_TRANSFER_LIMIT {
		apperrors.ClientError(ctx.Ctx, fmt.Sprintf("You cannot send more than ₦%s at a time", currencyformatter.HumanReadableIntCurrency(constants.MAXIMUM_LOCAL_TRANSFER_LIMIT)), nil, nil, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	localProcessorFee, localProcessorVAT, polymerFee, polymerVAT := utils.GetLocalTransactionFee(ctx.Body.Amount)
	totalAmount := ctx.Body.Amount + utils.Float32ToUint64Currency(localProcessorFee, false) + utils.Float32ToUint64Currency(polymerFee, false) + utils.Float32ToUint64Currency(localProcessorVAT, false) + utils.Float32ToUint64Currency(polymerVAT, false)
	bankName := ""
	for _, bank := range bankssupported.SupportedLocalBanks {
		if bank.Code == ctx.Body.BankCode {
			bankName = bank.Name
			break
		}
	}
	if bankName == "" {
		apperrors.NotFoundError(ctx.Ctx, "Selected bank is not currently supported", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	wallet, err := services.InitiatePreAuth(ctx.Ctx, nil, ctx.GetStringContextData("UserID"), totalAmount, ctx.Body.Pin, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		return
	}
	narration := func() string {
		if ctx.Body.Description == nil {
			des := fmt.Sprintf("NGN Transfer from %s %s to %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName"), bankName)
			return des
		}
		return *ctx.Body.Description
	}()
	reference := utils.GenerateUUIDString()

	// lock amount to be sent
	err = services.LockFunds(ctx.Ctx, wallet, ctx.Body.Amount, entities.LocalDebit, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking amount to be sent"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// lock local transaction fee
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(localProcessorFee, false), entities.LocalDebitFee, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking chimoney transaction fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// lock local debit vat
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(localProcessorVAT, false), entities.LocalDebitVAT, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking polymer processing fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}
	// lock polymer fee
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(polymerFee, false), entities.PolymerFee, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking polymer vat fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	// lock polymer vat
	err = services.LockFunds(ctx.Ctx, wallet, utils.Float32ToUint64Currency(polymerVAT, false), entities.PolymerVAT, reference, ctx.GetHeader("Polymer-Device-Id"))
	if err != nil {
		logger.Error(errors.New("error locking polymer vat fee"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "payload",
			Data: ctx.Body,
		})
		return
	}

	if os.Getenv("ENV") != "prod" {
		reference = fmt.Sprintf("%s_PMCKDU_1", reference)
	}
	response := services.InitiateLocalPayment(ctx.Ctx, &types.InitiateLocalTransferPayload{
		Reference: reference,
		Destination: types.LocalTransferDestination{
			Amount:    utils.UInt64ToFloat32Currency(ctx.Body.Amount),
			Narration: narration,
			BankAccount: types.LocalTransferDestinationBankAccount{
				Bank:    ctx.Body.BankCode,
				Account: ctx.Body.AccountNumber,
			},
			Customer: types.LocalTransferDestinationCustomer{
				Name: *ctx.Body.FullName,
			},
		},
	}, ctx.GetHeader("Polymer-Device-Id"))
	if response == nil {
		services.ReverseLockFunds(wallet.ID, reference)
		return
	}
	transaction := entities.Transaction{
		TransactionReference: reference,
		MetaData:             response,
		Fee:                  utils.Float32ToUint64Currency(polymerFee, false),
		Vat:                  utils.Float32ToUint64Currency(polymerVAT, false),
		ProcessorFeeCurrency: "NGN",
		ProcessorFee:         utils.Float32ToUint64Currency(localProcessorFee, false),
		ProcessorFeeVAT:      utils.Float32ToUint64Currency(localProcessorVAT, false),
		Amount:               totalAmount,
		Currency:             "NGN",
		WalletID:             wallet.ID,
		UserID:               wallet.UserID,
		Description:          narration,
		Location: entities.Location{
			IPAddress: ctx.Body.IPAddress,
		},
		Status: entities.TransactionPending,
		Intent: entities.LocalDebit,
		DeviceInfo: &entities.DeviceInfo{
			IPAddress: ctx.Body.IPAddress,
			DeviceID:  utils.GetStringPointer(ctx.GetStringContextData("DeviceID")),
			UserAgent: utils.GetStringPointer(ctx.GetStringContextData("UserAgent")),
		},
		Sender: entities.TransactionSender{
			FullName: fmt.Sprintf("%s %s", ctx.GetStringContextData("FirstName"), ctx.GetStringContextData("LastName")),
			Email:    utils.GetStringPointer(ctx.GetStringContextData("Email")),
		},
		Recepient: entities.TransactionRecepient{
			FullName:      response.FullName,
			BankCode:      &ctx.Body.BankCode,
			AccountNumber: ctx.Body.AccountNumber,
			BranchCode:    ctx.Body.BranchCode,
			BankName:      &bankName,
			Country:       utils.GetStringPointer("Nigeria"),
		},
	}
	trxRepository := repository.TransactionRepo()
	trx, err := trxRepository.CreateOne(context.TODO(), transaction)
	if err != nil {
		logger.Error(errors.New("error creating transaction for international transfer"), logger.LoggerOptions{
			Key:  "payload",
			Data: transaction,
		})
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if ctx.GetBoolContextData("PushNotifOptions") {
		pushnotification.PushNotificationService.PushOne(ctx.GetStringContextData("PushNotificationToken"), "Your payment is on its way! 🚀",
			fmt.Sprintf("Your payment of %s%d to %s in %s is currently being processed.", utils.CurrencyCodeToCurrencySymbol(transaction.Currency), transaction.Amount, transaction.Recepient.FullName, utils.CountryCodeToCountryName(*transaction.Recepient.Country)))
	}

	if ctx.GetBoolContextData("EmailOptions") {
		emails.EmailService.SendEmail(ctx.GetStringContextData("Email"), "Your payment is on its way! 🚀", "payment_sent", map[string]any{
			"FIRSTNAME":         transaction.Sender.FullName,
			"CURRENCY_CODE":     utils.CurrencyCodeToCurrencySymbol("NGN"),
			"AMOUNT":            utils.UInt64ToFloat32Currency(ctx.Body.Amount),
			"RECEPIENT_NAME":    transaction.Recepient.FullName,
			"RECEPIENT_COUNTRY": utils.CountryCodeToCountryName(*transaction.Recepient.Country),
		})
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusCreated, "Your payment is on its way! 🚀", trx, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func VerifyLocalAccountName(ctx *interfaces.ApplicationContext[dto.NameVerificationDTO]) {
	bankCode := ""
	for _, bank := range bankssupported.SupportedLocalBanks {
		if bank.Name == ctx.Body.BankName {
			bankCode = bank.Code
			break
		}
	}
	if bankCode == "" {
		apperrors.NotFoundError(ctx.Ctx, fmt.Sprintf("%s is not a supported bank on our platform yet.", ctx.Body.BankName), ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	name := services.NameVerification(ctx.Ctx, ctx.Body.AccountNumber, bankCode, ctx.GetHeader("Polymer-Device-Id"))
	if name == nil {
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "name verification complete", map[string]string{
		"name": *name,
	}, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func FetchPastBusinessTransactions(ctx *interfaces.ApplicationContext[any]) {
	transactionsRepo := repository.TransactionRepo()
	transactions, err := transactionsRepo.FindMany(map[string]interface{}{
		"userID":     ctx.GetStringContextData("UserID"),
		"businessID": ctx.GetStringParameter("businessID"),
	}, &options.FindOptions{
		Limit: utils.GetInt64Pointer(15),
		Sort: map[string]any{
			"createdAt": -1,
		},
	}, options.Find().SetProjection(map[string]any{
		"transactionReference": 1,
		"amount":               1,
		"amountInNGN":          1,
		"fee":                  1,
		"description":          1,
		"amountInUSD":          1,
		"transactionRecepient": 1,
		"transactionSender":    1,
		"status":               1,
		"currency":             1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "transctions fetched", transactions, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func FetchPastPersonalTransactions(ctx *interfaces.ApplicationContext[any]) {
	transactionsRepo := repository.TransactionRepo()
	transactions, err := transactionsRepo.FindManyStripped(map[string]interface{}{
		"userID": ctx.GetStringContextData("UserID"),
	}, &options.FindOptions{
		Limit: utils.GetInt64Pointer(15),
		Sort: map[string]any{
			"createdAt": -1,
		},
	}, options.Find().SetProjection(map[string]any{
		"transactionReference": 1,
		"amount":               1,
		"amountInNGN":          1,
		"fee":                  1,
		"description":          1,
		"amountInUSD":          1,
		"transactionRecepient": 1,
		"transactionSender":    1,
		"status":               1,
		"currency":             1,
	}))
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "transctions fetched", transactions, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}

func RequestAccountStatement(ctx *interfaces.ApplicationContext[dto.RequestAccountStatementDTO]) {
	walletRepo := repository.WalletRepo()
	exists, err := walletRepo.CountDocs(map[string]interface{}{
		"_id":    ctx.Body.WalletID,
		"userID": ctx.GetStringContextData("UserID"),
	})
	if err != nil {
		apperrors.FatalServerError(ctx.Ctx, err, ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	if exists != 1 {
		apperrors.NotFoundError(ctx.Ctx, "wallet not found", ctx.GetHeader("Polymer-Device-Id"))
		return
	}
	background.Scheduler.Emit("generate_account_statement", map[string]any{
		"walletID": ctx.Body.WalletID,
		"email":    ctx.Body.Email,
		"start":    ctx.Body.Start,
		"end":      ctx.Body.End,
	})

	server_response.Responder.Respond(ctx.Ctx, http.StatusOK, "statement is being generated", nil, nil, nil, ctx.GetHeader("Polymer-Device-Id"))
}
