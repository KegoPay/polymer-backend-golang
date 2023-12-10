package datastore

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kego.com/infrastructure/logger"
)

var (
	UserModel *mongo.Collection
	WalletModel *mongo.Collection
	FrozenWalletLogModel *mongo.Collection
	BusinessModel *mongo.Collection
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
		logger.Warning("an error occured while starting the database", logger.LoggerOptions{Key: "error", Data: err})
		return &cancel
	}

	db := client.Database(os.Getenv("DB_NAME"))
	setUpIndexes(ctx, db)

	logger.Info("connected to mongodb successfully")
	return &cancel
}

// Set up the indexes for the database
func setUpIndexes(ctx context.Context, db *mongo.Database) {
	UserModel = db.Collection("Users")
	UserModel.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}, {
		Keys:    bson.D{{Key: "bvn", Value: 1}},
		Options: options.Index().SetUnique(true),
	},})

	WalletModel = db.Collection("Wallets")
	WalletModel.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys:    bson.D{{Key: "businessID", Value: 1}},
		Options: options.Index().SetUnique(true),
	},{
		Keys:    bson.D{{Key: "userID", Value: 1}},
		Options: options.Index(),
	}})

	BusinessModel = db.Collection("Businesses")
	BusinessModel.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys:    bson.D{{Key: "walletID", Value: 1}},
		Options: options.Index().SetUnique(true),
	},{
		Keys:    bson.D{{Key: "userID", Value: 1}},
		Options: options.Index(),
	}})

	FrozenWalletLogModel = db.Collection("FrozenWalletLogs")
	
	logger.Info("mongodb indexes set up successfully")
}
