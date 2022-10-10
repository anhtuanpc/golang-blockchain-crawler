package db

import (
	"cbridgewrapper/config"
	"cbridgewrapper/entity"
	"cbridgewrapper/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var relayTransaction *mongo.Collection

func init() {
	logger.Logger.Infof("Initailize mongo db.")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbHost := config.GetConfig("dbhost")
	if dbHost == "" {
		logger.Logger.Fatal("Invalid db host")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbHost))
	if err != nil {
		logger.Logger.Fatal("Can't connect %v", err)
	}
	relayTransaction = client.Database("parse").Collection("RelayTransaction")
}

func UpdateRelayStatus(m primitive.M, d primitive.D) error {
	_, err := relayTransaction.UpdateOne(context.Background(), m, d)
	return err
}
func FindRelayTransaction(d bson.D) (*mongo.Cursor, error) {
	return relayTransaction.Find(context.Background(), d)
}
func FindAllPendingRelayTransaction() ([]entity.RelayTransaction, error) {
	cursor, err := relayTransaction.Find(context.Background(), bson.D{{"status", "PENDING"}})
	if err != nil {
		return nil, err
	}
	var trans []entity.RelayTransaction
	for cursor.Next(context.Background()) {
		var rs entity.RelayTransaction
		err := cursor.Decode(&rs)
		if err != nil {
			logger.Logger.Errorf("Can't parse record %v", err)
			continue
		}
		trans = append(trans, rs)
	}
	return trans, err
}
