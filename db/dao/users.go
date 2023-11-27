package dao

import (
	"context"
	"log"

	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	driver "go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(ctx context.Context, name string) (*db.User, error) {
	var newUser = db.User{
		Name: name,
	}
	if result, err := mongo.CollUsers.InsertOne(ctx, newUser); err != nil {
		log.Printf("Create user %s error: %+v", name, err)
		if driver.IsDuplicateKeyError(err) {
			return nil, &db.ErrDuplicateKey{Internal: err}
		}
		return nil, err
	} else {
		log.Printf("Create user %s result: %+v", name, result)
		newUser.ID = result.InsertedID.(primitive.ObjectID)
	}
	return &newUser, nil
}

func LoadUser(name string) (*db.User, error) {
	// TODO
	return nil, nil
}
