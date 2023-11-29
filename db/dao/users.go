package dao

import (
	"context"
	"log"

	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	driver "go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(ctx context.Context, name string) (*db.User, error) {
	var newUser = db.User{
		Name:       name,
		NextPicSeq: 1,
	}
	if result, err := mongo.CollUsers.InsertOne(ctx, &newUser); err != nil {
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

func IncrPhotoSeq(ctx context.Context, id primitive.ObjectID) error {
	count, err := mongo.CollUsers.UpdateByID(ctx, id, bson.M{
		"$inc": bson.M{
			"next_pic_seq": 1,
		},
	})
	log.Println("IncrPhotoSeq for", id, "result:", count)
	return err
}

func GetUserByName(ctx context.Context, name string) (*db.User, error) {
	var existingUser db.User
	if err := mongo.CollUsers.FindOne(ctx, bson.M{"name": name}).Decode(&existingUser); err != nil {
		if err == driver.ErrNoDocuments {
			return nil, db.ErrNotFound
		}
		return nil, err
	} else {
		return &existingUser, err
	}
}

func GetUserByID(ctx context.Context, id string) (*db.User, error) {
	var existingUser db.User
	if err := mongo.CollUsers.FindOne(ctx, bson.M{"_id": id}).Decode(&existingUser); err != nil {
		if err == driver.ErrNoDocuments {
			return nil, db.ErrNotFound
		}
		return nil, err
	} else {
		return &existingUser, err
	}
}
