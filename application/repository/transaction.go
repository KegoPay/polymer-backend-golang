package repository

import (
	"sync"

	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/database/connection/datastore"
	"usepolymer.co/infrastructure/database/repository/mongo"
)

var transactionOnce = sync.Once{}

var transactionRepository mongo.MongoRepository[entities.Transaction]

func TransactionRepo() *mongo.MongoRepository[entities.Transaction] {
	transactionOnce.Do(func() {
		transactionRepository = mongo.MongoRepository[entities.Transaction]{Model: datastore.TransactionModel}
	})
	return &transactionRepository
}
