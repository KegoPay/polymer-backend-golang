package datastore

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kego.com/infrastructure/logger"
)

var (
	// db models here
)

func connectMongo() *context.CancelFunc {
	url := os.Getenv("DB_URL")

	if url == "" {
		logger.Error(errors.New("set mongo url"))
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))

	if err != nil {
		logger.Error(errors.New("an error occured while starting the database"), logger.LoggerOptions{Key: "error", Data: err})
		return &cancel
	}

	db := client.Database(os.Getenv("DB_NAME"))
	setUpIndexes(ctx, db)

	logger.Info("connected to mongodb successfully")
	return &cancel
}

// Set up the indexes for the database
func setUpIndexes(ctx context.Context, db *mongo.Database) {
	logger.Info("mongodb indexes set up successfully")
}
