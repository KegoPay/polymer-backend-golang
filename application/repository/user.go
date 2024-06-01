package repository

import (
	"sync"

	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/database/connection/datastore"
	"usepolymer.co/infrastructure/database/repository/mongo"
)

var userOnce = sync.Once{}

var userRepository mongo.MongoRepository[entities.User]

func UserRepo() *mongo.MongoRepository[entities.User] {
	userOnce.Do(func() {
		userRepository = mongo.MongoRepository[entities.User]{Model: datastore.UserModel}
	})
	return &userRepository
}
