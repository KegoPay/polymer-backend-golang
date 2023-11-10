package repository

import (
	"sync"

	"kego.com/entities"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/database/repository/mongo"
)


var businessOnce = sync.Once{}

var businessRepository mongo.MongoRepository[entities.Business]

func BusinessRepo() *mongo.MongoRepository[entities.Business] {
	businessOnce.Do(func() {
		businessRepository = mongo.MongoRepository[entities.Business]{Model: datastore.BusinessModel}
	})
	return &businessRepository
}
