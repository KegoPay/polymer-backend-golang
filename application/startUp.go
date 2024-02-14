package application

import (
	"errors"
	"fmt"
	"os"

	"kego.com/application/repository"
	"kego.com/application/utils"
	"kego.com/entities"
	"kego.com/infrastructure/logger"
)

// Performs default operations on the database that need to be performed
// Only when the project is first initialised
func DBGenesis() {
	walletRepo := repository.WalletRepo()
	count, err := walletRepo.CountDocs(map[string]interface{}{
		"userID": os.Getenv("POLYMER_WALLET_USER_ID"),
	})
	if err != nil {
		logger.Error(errors.New("error counting genesis wallets"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		panic("error counting genesis wallets")
	}
	if count == 4 {
		return
	}
	if  count > 0 {
		logger.Warning(fmt.Sprintf("default wallet count - %d", count))
		panic(fmt.Sprintf("default wallet count - %d", count))
	}
	created, err := walletRepo.CreateBulk([]entities.Wallet{
		{
			Currency: "NGN",
			BusinessName: utils.GetStringPointer("Polymer Intl Fee Wallet"),
			BusinessID: nil,
			UserID: os.Getenv("POLYMER_WALLET_USER_ID"),
		},
		{
			Currency: "NGN",
			BusinessName: utils.GetStringPointer("Polymer Local VAT Wallet"),
			BusinessID: nil,
			UserID: os.Getenv("POLYMER_WALLET_USER_ID"),
		},
		{
			Currency: "NGN",
			BusinessName: utils.GetStringPointer("Polymer Local Fee Wallet"),
			BusinessID: nil,
			UserID: os.Getenv("POLYMER_WALLET_USER_ID"),
		},
		{
			Currency: "NGN",
			BusinessName: utils.GetStringPointer("Polymer Intl VAT Wallet"),
			BusinessID: nil,
			UserID: os.Getenv("POLYMER_WALLET_USER_ID"),
		},
	})
	if err != nil {
		logger.Error(errors.New("error creating genesis wallets"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		panic("error creating genesis wallets")
	}
	if len(*created) != 4 {
		logger.Error(errors.New("failed to create the expeted 4 genesis wallets"), logger.LoggerOptions{
			Key: "created",
			Data: created,
		})
		panic("failed to create the expeted 4 genesis wallets")
	}
	logger.Info("genesis wallets created successfully", logger.LoggerOptions{
		Key: "number",
		Data: len(*created),
	})
}