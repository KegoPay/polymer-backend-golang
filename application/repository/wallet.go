package repository

import (
	"sync"

	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/database/connection/datastore"
	"usepolymer.co/infrastructure/database/repository/mongo"
)

var walletOnce = sync.Once{}

var walletRepository mongo.MongoRepository[entities.Wallet]

func WalletRepo() *mongo.MongoRepository[entities.Wallet] {
	walletOnce.Do(func() {
		walletRepository = mongo.MongoRepository[entities.Wallet]{Model: datastore.WalletModel}
	})
	return &walletRepository
}
