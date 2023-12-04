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

func AddPhoto(ctx context.Context, fields *db.Photo) error {
	if result, err := mongo.CollPhotos.InsertOne(ctx, fields); err != nil {
		log.Printf("Create photo %+v error: %+v", fields, err)
		if driver.IsDuplicateKeyError(err) {
			return &db.ErrDuplicateKey{Internal: err}
		}
		return err
	} else {
		log.Printf("Create photo %+v result: %+v", fields, result)
		fields.ID = result.InsertedID.(primitive.ObjectID)
	}
	return nil
}

func UpdatePhotoInitProgress(ctx context.Context, photoID primitive.ObjectID, status string, progress int) error {
	initComplete := (status == "succeeded")
	if status == "succeeded" || status == "failed" {
		progress = 100
	}
	updatePredicate := bson.M{
		"init_complete": initComplete,
		"init_status":   status,
		"init_progress": progress,
	}
	count, err := mongo.CollPhotos.UpdateByID(ctx, photoID, bson.M{
		"$set": updatePredicate,
	})
	log.Println("UpdatePhotoInitProgress for", photoID, "result:", count)
	return err
}

func LoadUserPhotos(ctx context.Context, userID string) ([]db.Photo, error) {
	cursor, err := mongo.CollPhotos.Find(ctx, bson.M{"user": userID})
	if err != nil {
		log.Println("Error acquiring cursor:", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var ret []db.Photo
	// TODO: performance
	if err = cursor.All(ctx, &ret); err != nil {
		log.Println("Error executing All():", err)
		return nil, err
	}
	return ret, nil
}
