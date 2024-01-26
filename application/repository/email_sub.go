package repository

import (
	"sync"

	"kego.com/entities"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/database/repository/mongo"
)


var emailSubOnce = sync.Once{}

var emailSubRepository mongo.MongoRepository[entities.Subscriptions]

func EmailSubRepo() *mongo.MongoRepository[entities.Subscriptions] {
	emailSubOnce.Do(func() {
		emailSubRepository = mongo.MongoRepository[entities.Subscriptions]{Model: datastore.EmailSubs}
	})
	return &emailSubRepository
}
