package repository

import (
	"sync"

	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/database/connection/datastore"
	"usepolymer.co/infrastructure/database/repository/mongo"
)

var frozenWalletLogOnce = sync.Once{}

var frozenWalletLogRepository mongo.MongoRepository[entities.FrozenWalletLog]

func FrozenWalletLogRepo() *mongo.MongoRepository[entities.FrozenWalletLog] {
	frozenWalletLogOnce.Do(func() {
		frozenWalletLogRepository = mongo.MongoRepository[entities.FrozenWalletLog]{Model: datastore.FrozenWalletLogModel}
	})
	return &frozenWalletLogRepository
}
