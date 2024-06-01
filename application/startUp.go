package application

import (
	"errors"
	"fmt"
	"os"

	"usepolymer.co/application/repository"
	"usepolymer.co/application/utils"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/logger"
)

// Performs default operations on the database that need to be performed
// only when the project is first initialised
func DBGenesis() {
	walletRepo := repository.WalletRepo()
	count, err := walletRepo.CountDocs(map[string]interface{}{
		"userID": os.Getenv("POLYMER_WALLET_USER_ID"),
	})
	if err != nil {
		logger.Error(errors.New("error counting genesis wallets"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		panic("error counting genesis wallets")
	}
	if count == 2 {
		return
	}
	if count > 0 {
		logger.Warning("wrong default wallet count", logger.LoggerOptions{
			Key:  "wallet count",
			Data: count,
		})
		panic(fmt.Sprintf("default wallet count - %d", count))
	}
	created, err := walletRepo.CreateBulk([]entities.Wallet{
		{
			Currency:     "NGN",
			BusinessName: utils.GetStringPointer("Polymer Fee Wallet"),
			BusinessID:   nil,
			UserID:       os.Getenv("POLYMER_WALLET_USER_ID"),
		},
		{
			Currency:     "NGN",
			BusinessName: utils.GetStringPointer("Polymer VAT Wallet"),
			BusinessID:   nil,
			UserID:       os.Getenv("POLYMER_WALLET_USER_ID"),
		},
	})
	if err != nil {
		logger.Error(errors.New("error creating genesis wallets"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		panic("error creating genesis wallets")
	}
	if len(*created) != 2 {
		logger.Error(errors.New("failed to create the expeted 2 genesis wallets"), logger.LoggerOptions{
			Key:  "created",
			Data: created,
		})
		panic("failed to create the expeted 5 genesis wallets")
	}
	logger.Info("genesis wallets created successfully", logger.LoggerOptions{
		Key:  "number",
		Data: len(*created),
	})
}
