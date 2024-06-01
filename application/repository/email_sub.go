package repository

import (
	"sync"

	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/database/connection/datastore"
	"usepolymer.co/infrastructure/database/repository/mongo"
)

var emailSubOnce = sync.Once{}

var emailSubRepository mongo.MongoRepository[entities.Subscriptions]

func EmailSubRepo() *mongo.MongoRepository[entities.Subscriptions] {
	emailSubOnce.Do(func() {
		emailSubRepository = mongo.MongoRepository[entities.Subscriptions]{Model: datastore.EmailSubs}
	})
	return &emailSubRepository
}
