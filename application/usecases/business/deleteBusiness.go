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

func DeleteBusiness(ctx any, id string, device_id *string) error {
	businessRepo := repository.BusinessRepo()
	var e error
	businessRepo.StartTransaction(func(sc mongo.Session, c context.Context) error {
		deleted, err := businessRepo.DeleteOne(c, map[string]interface{}{
			"_id": id,
		})
		if err != nil {
			(sc).AbortTransaction(c)
			logger.Error(errors.New("error deleting business"), logger.LoggerOptions{
				Key: "error",
				Data: err,
			}, logger.LoggerOptions{
				Key: "id",
				Data: id,
			})
			e =  err
			apperrors.FatalServerError(ctx, err, device_id)
			return err
		}
		if deleted == 0 {
			(sc).AbortTransaction(c)
			err = errors.New("business does not exist")
			e =  err
			apperrors.NotFoundError(ctx, err.Error(), device_id)
			return err
		}
		err = walletUsecases.DeleteWallet(ctx, c, id, device_id)
		if err != nil {
			e =  err
		(sc).AbortTransaction(c)
			return err
		}
		
		(sc).CommitTransaction(c)
		return err
	})


	return e
}