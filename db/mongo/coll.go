package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	tblUsers = "users"
)

var (
	CollUsers *mongo.Collection
)

func initCollections() {
	CollUsers = database.Collection(tblUsers)
	idxUserName := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}
	if createIndexResult, err := CollUsers.Indexes().CreateOne(context.TODO(), idxUserName); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Creating idxUserName...", createIndexResult)
	}
}
