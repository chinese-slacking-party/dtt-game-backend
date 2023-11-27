package mongo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chinese-slacking-party/dtt-game-backend/config"
)

var (
	client   *mongo.Client
	database *mongo.Database

	moduleLock sync.Mutex
)

func Init(db string) error {
	moduleLock.Lock()
	defer moduleLock.Unlock()
	if client != nil {
		client.Disconnect(context.TODO())
		database = nil
		client = nil
	}
	return doInit(db)
}

func doInit(db string) (err error) {
	clientOptions := options.Client().ApplyURI(config.DBAddr)
	if client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
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
