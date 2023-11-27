package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	Name   string    `bson:"name" json:"name"`
	Nick   string    `bson:"nick" json:"nickname"`
	Avatar string    `bson:"avatar" json:"avatar"`
	Stats  UserStats `bson:"stats" json:"stats"`

	NextPicSeq int32   `bson:"next_pic_seq" json:"-"`
	Album      []Photo `bson:"-" json:"album"` // Filled by querying the `photos` collection
}

type UserStats struct {
	GamesPlayed int32 `bson:"played" json:"games_played"`
	GamesWon    int32 `bson:"won" json:"wins"`
}

type Photo struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	Seq      int32             `bson:"seq" json:"id"`
	Key      string            `bson:"key" json:"key"`
	Desc     string            `bson:"desc" json:"desc"`
	Original string            `bson:"original" json:"original"`
	URLs     map[string]string `bson:"urls" json:"urls,omitempty"`

	// ObjectId for user-uploaded, or some predefined value for system-generated
	UserID string `bson:"user" json:"-"`
}

type Level struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	// TODO
}

type Game struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	// These could be from `levels` collection instead, but since we're using
	// hard-coded levels for this version, it's reasonable to put them here.
	Title       string   `bson:"title" json:"title,omitempty"`
	SuccessText []string `bson:"success_text" json:"success_text,omitempty"`
	FailureText []string `bson:"failure_text" json:"failure_text,omitempty"`
	Prompt      string   `bson:"prompt" json:"prompt,omitempty"`
	Width       int32    `bson:"width" json:"width,omitempty"`
	Height      int32    `bson:"height" json:"height,omitempty"`

	Map [][]MapTile `bson:"map" json:"map"`

	LevelID   string    `bson:"level" json:"-"`
	CreatedAt time.Time `bson:"created_at" json:"-"`
	UpdatedAt time.Time `bson:"updated_at" json:"-"`
	Status    string    `bson:"status" json:"status"`

	UserID string `bson:"user" json:"-"`
}

type MapTile struct {
	ImageID  string `bson:"image_id" json:"image_id"`
	ImageTag string `bson:"image_tag" json:"image_tag"`
	URL      string `bson:"-" json:"url"`
}

var (
	CookieLife = 168 * time.Hour
)
