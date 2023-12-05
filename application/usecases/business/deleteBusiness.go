package business

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	apperrors "kego.com/application/appErrors"
	"kego.com/application/repository"
	walletUsecases "kego.com/application/usecases/wallet"
	"kego.com/infrastructure/logger"
)

func DeleteBusiness(ctx any, id string) error {
	businessRepo := repository.BusinessRepo()

	businessRepo.StartTransaction(func(sc *mongo.SessionContext, c *context.Context) error {
		deleted, err := businessRepo.DeleteOne(c, map[string]interface{}{
			"_id": id,
		})
		if err != nil {
			logger.Error(errors.New("error deleting business"), logger.LoggerOptions{
				Key: "error",
				Data: err,
			}, logger.LoggerOptions{
				Key: "id",
				Data: id,
			})
			apperrors.FatalServerError(ctx)
			return err
		}
		if deleted == 0 {
			apperrors.NotFoundError(ctx, "business does not exist")
			return errors.New("")
		}
		err = walletUsecases.DeleteWallet(ctx, c, id)
		return err
	})


	return nil
}