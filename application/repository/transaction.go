package repository

import (
	"sync"

	"kego.com/entities"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/database/repository/mongo"
)


var transactionOnce = sync.Once{}

var transactionRepository mongo.MongoRepository[entities.Transaction]

func TransactionRepo() *mongo.MongoRepository[entities.Transaction] {
	transactionOnce.Do(func() {
		transactionRepository = mongo.MongoRepository[entities.Transaction]{Model: datastore.TransactionModel}
	})
	return &transactionRepository
}
