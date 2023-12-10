package repository

import (
	"sync"

	"kego.com/entities"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/database/repository/mongo"
)


var frozenWalletLogOnce = sync.Once{}

var frozenWalletLogRepository mongo.MongoRepository[entities.FrozenWalletLog]

func FrozenWalletLogRepo() *mongo.MongoRepository[entities.FrozenWalletLog] {
	frozenWalletLogOnce.Do(func() {
		frozenWalletLogRepository = mongo.MongoRepository[entities.FrozenWalletLog]{Model: datastore.FrozenWalletLogModel}
	})
	return &frozenWalletLogRepository
}
