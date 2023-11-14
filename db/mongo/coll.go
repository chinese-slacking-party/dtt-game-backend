package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	tblUsers = "users"
)

var (
	CollUsers *mongo.Collection
)

func initCollections() {
	CollUsers = database.Collection(tblUsers)
}
