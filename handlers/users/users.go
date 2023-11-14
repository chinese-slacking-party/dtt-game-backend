package users

import (
	"context"
	"net/http"
	"time"

	"github.com/chinese-slacking-party/dtt-game-backend/db/mongo"

	"github.com/gin-gonic/gin"
)

type User struct {
	// TODO: constraint
	Name string `json:"name" bson:"name"`
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := doRegister(&user)
	// TODO: "user already exists" error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		// TODO: profile struct
		"profile": map[string]interface{}{},
		"debug":   result,
	})
}

func doRegister(user *User) (map[string]interface{}, error) {
	// Insert a new user into the collection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := mongo.CollUsers.InsertOne(ctx, user)
	return map[string]interface{}{"insertedID": result.InsertedID}, err
}
