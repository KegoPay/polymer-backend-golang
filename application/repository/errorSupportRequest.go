package repository

import (
	"sync"

	"kego.com/entities"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/database/repository/mongo"
)


var errorSupportRequestOnce = sync.Once{}

var ErrorSupportRequestRepository mongo.MongoRepository[entities.ErrorSupportRequest]

func ErrorSupportRequestRepo() *mongo.MongoRepository[entities.ErrorSupportRequest] {
	errorSupportRequestOnce.Do(func() {
		ErrorSupportRequestRepository = mongo.MongoRepository[entities.ErrorSupportRequest]{Model: datastore.ErrorSupportRequestModel}
	})
	return &ErrorSupportRequestRepository
}
