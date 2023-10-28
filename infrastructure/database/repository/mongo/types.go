package mongo

import "go.mongodb.org/mongo-driver/mongo"


type MongoModels interface {
	MongoDBName() string
}

type ModelMethods interface {
	MarshalBSON() ([]byte, error)
	MarshalBinary() ([]byte, error)
}

type MongoRepository[T MongoModels] struct {
	Model   *mongo.Collection
	Payload interface{}
}

type FindOptions struct{
	Projection *interface{}
	Sort *interface{}
	Skip *int64
}