package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// TODO: config
	hardcodedDBAddr = "mongodb://localhost:27017"
)

var (
	client   *mongo.Client
	database *mongo.Database
)

func Init(db string) (err error) {
	clientOptions := options.Client().ApplyURI(hardcodedDBAddr)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		return
	}
	database = client.Database(db)
	return
}

func GetDB() *mongo.Database {
	return database
}
