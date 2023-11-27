package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestSaveLoadUser(t *testing.T) {
	dbName := fmt.Sprintf("dtt_test_sluser_%d", time.Now().UnixMicro())
	mongo.Init(dbName)

	{
		_, err := GetUserByName(context.Background(), "user001")
		assert.ErrorIs(t, err, db.ErrNotFound)
	}

	{
		_, err := CreateUser(context.Background(), "user002")
		assert.NoError(t, err)
		_, err = GetUserByName(context.Background(), "user002")
		assert.NoError(t, err)
	}

	assert.NoError(t, mongo.GetDB().Drop(context.Background()), "Unable to drop database - is your MongoDB sane?")
}

func TestAddRetrievePhoto(t *testing.T) {
	dbName := fmt.Sprintf("dtt_test_arphoto_%d", time.Now().UnixMicro())
	mongo.Init(dbName)

	{
		noPhoto, err := LoadUserPhotos(context.Background(), "user003")
		assert.NoError(t, err)
		assert.Empty(t, noPhoto)
	}

	{
		var x = db.Photo{
			Seq: 1,
			Key: "dad",
			URLs: map[string]string{
				"red":  "https://",
				"blue": "https://",
			},
			UserID: "user004",
		}
		assert.NoError(t, AddPhoto(context.Background(), &x))
		assert.NotEmpty(t, x.ID)
		x.ID = primitive.NilObjectID
		assert.Error(t, AddPhoto(context.Background(), &x))
		onePhoto, err := LoadUserPhotos(context.Background(), "user004")
		assert.NoError(t, err)
		assert.Len(t, onePhoto, 1)
	}

	assert.NoError(t, mongo.GetDB().Drop(context.Background()), "Unable to drop database - is your MongoDB sane?")
}
