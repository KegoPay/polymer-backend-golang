package repository

import (
	"sync"

	"kego.com/entities"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/database/repository/mongo"
)


var userOnce = sync.Once{}

var userRepository mongo.MongoRepository[entities.User]

func UserRepo() *mongo.MongoRepository[entities.User] {
	userOnce.Do(func() {
		userRepository = mongo.MongoRepository[entities.User]{Model: datastore.UserModel}
	})
	return &userRepository
}
