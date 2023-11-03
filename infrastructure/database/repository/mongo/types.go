package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"kego.com/infrastructure/database"
)


type MongoRepository[T database.BaseModel] struct {
	Model   *mongo.Collection
}


type FindOptions struct{
	Projection *interface{}
	Sort *interface{}
	Skip *int64
}