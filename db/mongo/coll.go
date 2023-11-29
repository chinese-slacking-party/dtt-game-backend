package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	tblUsers  = "users"
	tblPhotos = "photos"
	tblGames  = "games"
)

var (
	CollUsers  *mongo.Collection
	CollPhotos *mongo.Collection
	CollGames  *mongo.Collection
)

func initCollections() {
	CollUsers = database.Collection(tblUsers)
	idxUserName := mongo.IndexModel{
		Keys:    bson.D{{"name", 1}},
		Options: options.Index().SetUnique(true),
	}
	if createIndexResult, err := CollUsers.Indexes().CreateOne(context.TODO(),
		idxUserName); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Creating indexes for %s... %s", tblUsers, createIndexResult)
	}

	CollPhotos = database.Collection(tblPhotos)
	idxUserKey := mongo.IndexModel{
		Keys:    bson.D{{"user", 1}, {"key", 1}},
		Options: options.Index().SetUnique(true),
	}
	idxUserSeq := mongo.IndexModel{
		Keys:    bson.D{{"user", 1}, {"seq", 1}},
		Options: options.Index().SetUnique(true),
	}
	if createIndexResult, err := CollPhotos.Indexes().CreateMany(context.TODO(),
		[]mongo.IndexModel{idxUserKey, idxUserSeq}); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Creating indexes for %s... %s", tblPhotos, createIndexResult)
	}

	CollGames = database.Collection(tblGames)
	idxUserLevel := mongo.IndexModel{
		Keys:    bson.D{{"user", 1}, {"level", 1}},
		Options: options.Index().SetUnique(true),
	}
	if createIndexResult, err := CollGames.Indexes().CreateOne(context.TODO(),
		idxUserLevel); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Creating indexes for %s... %s", tblGames, createIndexResult)
	}
}
