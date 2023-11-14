package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user record in the database
type User struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

func main() {
	// Initialize the Gin engine.
	r := gin.Default()

	// Set client options and connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get a handle for your collection
	collection := client.Database("mydatabase").Collection("users")

	// Define a route for the index page.
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Gin backend service with MongoDB!",
		})
	})

	// Define a route for a health check.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// Define a POST route to receive and add a new user to the MongoDB.
	r.POST("/user", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert a new user into the collection
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, err := collection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully!",
			"user":    result,
		})
	})

	// Define a GET route to retrieve all users from the MongoDB.
	r.GET("/users", func(c *gin.Context) {
		var users []User

		// Find all users in the collection
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var user User
			cursor.Decode(&user)
			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	})

	// Start serving the application on port 8080.
	r.Run(":8080")
}
