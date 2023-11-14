package mongo

import (
	"context"
	"sync"

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

	// TODO: use lock instead (allow de-init and re-init)
	initDBOnce sync.Once
)

func Init(db string) (err error) {
	initDBOnce.Do(func() {
		err = doInit(db)
	})
	return
}

func doInit(db string) (err error) {
	clientOptions := options.Client().ApplyURI(hardcodedDBAddr)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		return
	}
	database = client.Database(db)
	initCollections()
	return
}

func GetDB() *mongo.Database {
	return database
}
