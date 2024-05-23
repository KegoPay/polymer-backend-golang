package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kego.com/infrastructure/database"
	"kego.com/infrastructure/logger"
)

func (repo *MongoRepository[T]) CreateOne(ctx context.Context, payload T, opts ...*options.InsertOneOptions) (*T, error) {
	var cancel context.CancelFunc
	if ctx == nil {
		c, ctxCancel := repo.createCtx()
		ctx = c
		cancel = ctxCancel
	}

	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	parsedPayload := interface{}(payload).(database.BaseModel).ParseModel()
	_, err := repo.Model.InsertOne(ctx, parsedPayload, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running CreateOne"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "payload",
			Data: payload,
		})
		if errParts := strings.Split(err.Error(), "E11000 duplicate key error collection:"); len(errParts) == 2 {
			return nil, fmt.Errorf("%s already exists", repo.Model.Name())
		}
		return nil, err
	}
	logger.Info("mongo CreateOne complete")
	return parsedPayload.(*T), err
}

func (repo *MongoRepository[T]) CreateBulk(payload []T, opts ...*options.InsertManyOptions) (*[]string, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()
	marshaled := []interface{}{}
	for _, i := range payload {
		marshaled = append(marshaled, interface{}(i).(database.BaseModel).ParseModel())
	}
	response, err := repo.Model.InsertMany(c, marshaled, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running CreateBulk"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "payload",
			Data: payload,
		})
		return nil, err
	}
	var ids []string
	for _, id := range response.InsertedIDs {
		ids = append(ids, id.(string))
	}
	logger.Info("CreateBulk complete")
	return &ids, err
}

func (repo *MongoRepository[T]) CreateBulkAndReturnPayload(payload []T, opts ...*options.InsertManyOptions) (*[]T, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()
	marshaled := []interface{}{}
	for _, i := range payload {
		marshaled = append(marshaled, interface{}(i).(database.BaseModel).ParseModel())
	}
	_, err := repo.Model.InsertMany(c, marshaled, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running CreateBulkAndReturnPayload"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "payload",
			Data: payload,
		})
		return nil, err
	}
	logger.Info("CreateBulkAndReturnPayload complete")
	return &payload, err
}

func (repo *MongoRepository[T]) FindOneByFilter(filter map[string]interface{}, opts ...*options.FindOneOptions) (*T, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()
	var result T
	doc := repo.Model.FindOne(c, filter, opts...)
	err := doc.Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		logger.Error(errors.New("mongo error occured while running FindOneByFilter"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return nil, err
	}
	logger.Info("FindOneByFilter complete")
	return &result, nil
}

func (repo *MongoRepository[T]) FindMany(filter map[string]interface{}, opts ...*options.FindOptions) (*[]T, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()
	var result []T
	cursor, err := repo.Model.Find(c, filter, opts...)
	if err != nil {
		return nil, err
	}
	err = cursor.All(c, &result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("no documents found")
		}
		logger.Error(errors.New("mongo error occured while running FindMany"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return nil, err
	}
	logger.Info("FindMany complete")
	return &result, nil
}

// FindManyStripped
// Strips all unneeded parts of the payload that does is nil
func (repo *MongoRepository[T]) FindManyStripped(filter map[string]interface{}, opts ...*options.FindOptions) (*[]map[string]interface{}, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()
	var result []map[string]interface{}
	cursor, err := repo.Model.Find(c, filter, opts...)
	if err != nil {
		return nil, err
	}
	err = cursor.All(c, &result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("no documents found")
		}
		logger.Error(errors.New("mongo error occured while running FindManyStripped"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return nil, err
	}
	logger.Info("FindManyStripped complete")
	return &result, nil
}

func (repo *MongoRepository[T]) FindByID(id string, opts ...*options.FindOneOptions) (*T, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()
	var result T
	err := repo.Model.FindOne(c, bson.M{"_id": id}, opts...).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		logger.Error(errors.New("mongo error occured while running FindById"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "resourceID",
			Data: id,
		})
		return nil, err
	}
	logger.Info("FindById complete")
	return &result, nil
}

func (repo *MongoRepository[T]) CountDocs(filter map[string]interface{}, opts ...*options.CountOptions) (int64, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()
	count, err := repo.Model.CountDocuments(c, filter, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running CountDocs"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return 0, err
	}
	logger.Info("CountDocs complete")
	return count, nil
}

func (repo *MongoRepository[T]) FindLast(opts ...*options.FindOptions) (*T, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	var lastRecord T
	err := repo.Model.FindOne(c, bson.M{}, options.FindOne().SetSort(bson.M{"$natural": -1})).Decode(&lastRecord)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running FindLast"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, err
	}
	logger.Info("FindLast complete")
	return &lastRecord, nil
}

func (repo *MongoRepository[T]) DeleteOne(ctx context.Context,  filter map[string]interface{}) (int64, error) {
	var cancel context.CancelFunc
	if ctx == nil {
		c, ctxCancel := repo.createCtx()
		ctx = c
		cancel = ctxCancel
	}

	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	result, err := repo.Model.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running DeleteOne"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return 0, err
	}
	logger.Info("DeleteOne complete")
	return result.DeletedCount, err
}

func (repo *MongoRepository[T]) DeleteByID(id string) (int64, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	result, err := repo.Model.DeleteOne(c, bson.M{"_id": &id})
	if err != nil {
		logger.Error(errors.New("mongo error occured while running DeleteByID"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "resourceID",
			Data: id,
		})
		return 0, err
	}
	logger.Info("DeleteByID complete")
	return result.DeletedCount, err
}

func (repo *MongoRepository[T]) DeleteMany(filter map[string]interface{}) (int64, error) {
	c, cancel := repo.createCtx()
	defer func() {
		cancel()
	}()

	count, err := repo.Model.DeleteMany(c, filter)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running DeleteMany"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return 0, err
	}
	logger.Info("DeleteMany complete")
	return count.DeletedCount, err
}

func (repo *MongoRepository[T]) UpdateByField(filter map[string]interface{}, payload *T, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateOne(c, filter, payload, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdateByField"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return false, err
	}
	logger.Info("UpdateByField complete")
	return true, err
}

func (repo *MongoRepository[T]) UpdateWithOperator(filter map[string]interface{}, payload map[string]interface{}, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateOne(c, filter, payload, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdateWithOperator"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return false, err
	}
	logger.Info("UpdateWithOperator complete")
	return true, err
}

func (repo *MongoRepository[T]) UpdateManyWithOperator(filter map[string]interface{}, payload map[string]interface{}, opts ...*options.UpdateOptions) (int64, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	affected, err := repo.Model.UpdateMany(c, filter, payload, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdateManyWithOperator"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return 0, err
	}
	return affected.ModifiedCount, err
}

func (repo *MongoRepository[T]) UpdateOrCreateByField(filter map[string]interface{}, payload map[string]interface{}, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateOne(c, filter, bson.D{primitive.E{Key: "$set", Value: payload}}, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdateOrCreateByField"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return false, err
	}
	logger.Info("UpdateOrCreateByField complete")
	return true, err
}

func (repo *MongoRepository[T]) UpdateOrCreateByFieldAndReturn(filter map[string]interface{}, payload T, opts ...*options.UpdateOptions) (*string, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	result, err := repo.Model.UpdateOne(c, filter, bson.D{primitive.E{Key: "$set", Value: &payload}}, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdateOrCreateByFieldAndReturn"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return nil, err
	}
	if result.UpsertedID == nil {
		return nil, nil
	}
	id := result.UpsertedID.(primitive.ObjectID).Hex()
	logger.Info("UpdateOrCreateByFieldAndReturn complete")
	return &id, err
}

func (repo *MongoRepository[T]) UpdateByID(id string, payload *T, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateByID(c, id, bson.D{primitive.E{Key: "$set", Value: *payload}}, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdateByID"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "resourceID",
			Data: id,
		})
		return false, err
	}
	logger.Info("UpdateByID complete")
	return true, err
}

func (repo *MongoRepository[T]) UpdatePartialByID(id string, payload interface{}, opts ...*options.UpdateOptions) (int64, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	result , err := repo.Model.UpdateByID(c, id, bson.D{primitive.E{Key: "$set", Value: payload}}, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdatePartialByID"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "resourceID",
			Data: id,
		})
		return 0, err
	}
	logger.Info("UpdatePartialByID complete")
	return result.MatchedCount, err
}

func (repo *MongoRepository[T]) UpdatePartialByFilter(filter map[string]interface{}, payload interface{}, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := repo.createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateMany(c, filter, bson.D{primitive.E{Key: "$set", Value: payload}}, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdatePartialByFilter"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return false, err
	}
	logger.Info("UpdatePartialByFilter complete")
	return true, err
}

func (repo MongoRepository[T]) StartTransaction(payload func(sc mongo.Session, c context.Context) error) error {
	session, err := repo.Model.Database().Client().StartSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.EndSession(context.Background())
	ctx := mongo.NewSessionContext(context.Background(), session)
	err = session.StartTransaction()
	if err != nil {
		log.Fatal(err)
	}
	if err := payload(session, ctx); err != nil {
		logger.Error(errors.New("mongo error occured while running StartTransaction"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	logger.Info("StartTransaction complete")
	return nil
}

func (repo *MongoRepository[T]) createCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}

// func (repo *MongoRepository[T])  marshalBSON(payload T) ([]byte, error) {
	
// }
