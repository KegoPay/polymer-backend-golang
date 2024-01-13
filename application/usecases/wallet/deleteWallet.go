package wallet

import (
	"context"
	"errors"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/repository"
	"kego.com/infrastructure/logger"
)

func DeleteWallet(ctx any, transactionCtx context.Context, businessID string) error {
	walletRepo := repository.WalletRepo()
	deleted, err := walletRepo.DeleteOne(transactionCtx, map[string]interface{}{
		"businessID": businessID,
	})
	if err != nil {
		logger.Error(errors.New("error deleting wallet"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		apperrors.FatalServerError(ctx, err)
		return err
	}
	if deleted == 0 {
		apperrors.NotFoundError(ctx, "wallet does not exist")
		return errors.New("")
	}
	return nil
}