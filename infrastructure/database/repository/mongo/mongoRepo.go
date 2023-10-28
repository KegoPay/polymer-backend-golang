package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kego.com/infrastructure/logger"
)

func (repo *MongoRepository[T]) CreateOne(payload T, opts ...*options.InsertOneOptions) (*T, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()
	parsed_payload := parsePayload(payload)
	_, err := repo.Model.InsertOne(c, parsed_payload, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running CreateOne"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "payload",
			Data: payload,
		})
		return nil, err
	}
	logger.Info("mongo CreateOne complete")
	return parsed_payload, err
}

func (repo *MongoRepository[T]) CreateBulk(payload []T, opts ...*options.InsertManyOptions) (*[]string, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()
	parsed_payload := parseMultiple(payload)
	marshaled := []interface{}{}
	for _, i := range parsed_payload {
		interface{}(i).(ModelMethods).MarshalBSON()
		interface{}(i).(ModelMethods).MarshalBinary()
		marshaled = append(marshaled, i)
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
		ids = append(ids, id.(primitive.ObjectID).Hex())
	}
	logger.Info("CreateBulk complete")
	return &ids, err
}

func (repo *MongoRepository[T]) CreateBulkAndReturnPayload(payload []T, opts ...*options.InsertManyOptions) ([]*T, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()
	parsed_payload := parseMultiple(payload)
	marshaled := []interface{}{}
	for _, i := range parsed_payload {
		interface{}(i).(ModelMethods).MarshalBSON()
		interface{}(i).(ModelMethods).MarshalBinary()
		marshaled = append(marshaled, i)
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
	return parsed_payload, err
}

func (repo *MongoRepository[T]) FindOneByFilter(filter map[string]interface{}, opts ...*options.FindOneOptions) (*T, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()
	var result T
	f := parseFilter(filter)
	doc := repo.Model.FindOne(c, f, opts...)
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
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()
	var result []T
	f := parseFilter(filter)
	cursor, err := repo.Model.Find(c, f, opts...)
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
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()
	var result []map[string]interface{}
	f := parseFilter(filter)
	cursor, err := repo.Model.Find(c, f, opts...)
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
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()
	var result T
	i := parseStringToMongoID(&id)
	err := repo.Model.FindOne(c, bson.M{"_id": i}, opts...).Decode(&result)
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
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()
	cc := parseFilter(filter)
	count, err := repo.Model.CountDocuments(c, cc, opts...)
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
	c, cancel := createCtx()

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

func (repo *MongoRepository[T]) DeleteOne(filter map[string]interface{}) (bool, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.DeleteOne(c, parseFilter(filter))
	if err != nil {
		logger.Error(errors.New("mongo error occured while running DeleteOne"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return false, err
	}
	logger.Info("DeleteOne complete")
	return true, err
}

func (repo *MongoRepository[T]) DeleteByID(id string) (bool, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.DeleteOne(c, bson.M{"_id": parseStringToMongoID(&id)})
	if err != nil {
		logger.Error(errors.New("mongo error occured while running DeleteByID"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "resourceID",
			Data: id,
		})
		return false, err
	}
	logger.Info("DeleteByID complete")
	return true, err
}

func (repo *MongoRepository[T]) DeleteMany(filter map[string]interface{}) (int64, error) {
	c, cancel := createCtx()
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
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateOne(c, parseFilter(filter), payload, opts...)
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
	c, cancel := createCtx()

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

func (repo *MongoRepository[T]) UpdateManyWithOperator(filter map[string]interface{}, payload map[string]interface{}, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateMany(c, filter, payload, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdateManyWithOperator"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "filter",
			Data: filter,
		})
		return false, err
	}
	logger.Info("UpdateManyWithOperator complete")
	return true, err
}

func (repo *MongoRepository[T]) UpdateOrCreateByField(filter map[string]interface{}, payload map[string]interface{}, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateOne(c, parseFilter(filter), bson.D{primitive.E{Key: "$set", Value: payload}}, opts...)
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
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	result, err := repo.Model.UpdateOne(c, parseFilter(filter), bson.D{primitive.E{Key: "$set", Value: &payload}}, opts...)
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
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateByID(c, parseStringToMongoID(&id), bson.D{primitive.E{Key: "$set", Value: *payload}}, opts...)
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

func (repo *MongoRepository[T]) UpdatePartialByID(id string, payload interface{}, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	_, err := repo.Model.UpdateByID(c, parseStringToMongoID(&id), bson.D{primitive.E{Key: "$set", Value: payload}}, opts...)
	if err != nil {
		logger.Error(errors.New("mongo error occured while running UpdatePartialByID"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "resourceID",
			Data: id,
		})
		return false, err
	}
	logger.Info("UpdatePartialByID complete")
	return true, err
}

func (repo *MongoRepository[T]) UpdatePartialByFilter(filter map[string]interface{}, payload interface{}, opts ...*options.UpdateOptions) (bool, error) {
	c, cancel := createCtx()

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

func (repo MongoRepository[T]) StartTransaction(payload func(sc *mongo.SessionContext, c *context.Context) error) error {
	c, cancel := createCtx()

	defer func() {
		cancel()
	}()

	if err := repo.Model.Database().Client().UseSession(c, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}
		return payload(&sc, &c)
	}); err != nil {
		logger.Error(errors.New("mongo error occured while running StartTransaction"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	logger.Info("StartTransaction complete")
	return nil
}

func parseFilter(f interface{}) interface{} {
	filter := (f).(map[string]interface{})
	if filter["_id"] != nil && reflect.TypeOf(filter["_id"]).String() == "string" {
		id := fmt.Sprintf("%v", filter["_id"])
		filter["_id"] = parseStringToMongoID(&id)
	}
	return filter
}

func createCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}

func parsePayload[T MongoModels](payload T) *T {
	byteA := dataToByteArray(payload)
	payload_map := *byteArrayToData[map[string]interface{}](byteA)
	if payload_map["Id"] == "000000000000000000000000" {
		payload_map["id"] = primitive.NewObjectID()
	} else if payload_map["Id"] != nil {
		payload_map["id"] = parseStringToMongoID(payload_map["Id"].(*string))
	} else if payload_map["Id"] == nil {
		payload_map["id"] = primitive.NewObjectID()
	}
	return byteArrayToData[T](dataToByteArray(payload_map))
}

func parseMultiple[T MongoModels](payload []T) []*T {
	var result []*T
	for _, data := range payload {
		result = append(result, parsePayload(data))
	}
	return result
}

// turns a byte array to the specified generic type
func byteArrayToData[T interface{}](payload []byte) *T {
	var data T
	json.Unmarshal(payload, &data)
	return &data
}

func dataToByteArray(payload interface{}) []byte {
	data, _ := json.Marshal(payload)
	return data
}

func parseStringToMongoID(id *string) primitive.ObjectID {
	objId, _ := primitive.ObjectIDFromHex(*id)
	return objId
}
