package repository

import (
	"sync"

	"kego.com/entities"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/database/repository/mongo"
)


var walletOnce = sync.Once{}

var walletRepository mongo.MongoRepository[entities.Wallet]

func WalletRepo() *mongo.MongoRepository[entities.Wallet] {
	walletOnce.Do(func() {
		walletRepository = mongo.MongoRepository[entities.Wallet]{Model: datastore.WalletModel}
	})
	return &walletRepository
}
