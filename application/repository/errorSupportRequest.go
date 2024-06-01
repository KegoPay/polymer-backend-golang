package repository

import (
	"sync"

	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/database/connection/datastore"
	"usepolymer.co/infrastructure/database/repository/mongo"
)

var errorSupportRequestOnce = sync.Once{}

var ErrorSupportRequestRepository mongo.MongoRepository[entities.ErrorSupportRequest]

func ErrorSupportRequestRepo() *mongo.MongoRepository[entities.ErrorSupportRequest] {
	errorSupportRequestOnce.Do(func() {
		ErrorSupportRequestRepository = mongo.MongoRepository[entities.ErrorSupportRequest]{Model: datastore.ErrorSupportRequestModel}
	})
	return &ErrorSupportRequestRepository
}
