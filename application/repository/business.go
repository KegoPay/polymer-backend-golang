package repository

import (
	"sync"

	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/database/connection/datastore"
	"usepolymer.co/infrastructure/database/repository/mongo"
)

var businessOnce = sync.Once{}

var businessRepository mongo.MongoRepository[entities.Business]

func BusinessRepo() *mongo.MongoRepository[entities.Business] {
	businessOnce.Do(func() {
		businessRepository = mongo.MongoRepository[entities.Business]{Model: datastore.BusinessModel}
	})
	return &businessRepository
}
