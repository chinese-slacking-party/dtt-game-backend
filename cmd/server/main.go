package main

import (
	"net/http"
	"os"

	"github.com/chinese-slacking-party/dtt-game-backend/config"
	"github.com/chinese-slacking-party/dtt-game-backend/db/mongo"
	"github.com/chinese-slacking-party/dtt-game-backend/handlers/album"
	"github.com/chinese-slacking-party/dtt-game-backend/handlers/game"
	"github.com/chinese-slacking-party/dtt-game-backend/handlers/session"
	"github.com/chinese-slacking-party/dtt-game-backend/handlers/users"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get config from environment.
	config.ReplicateAPIKey = os.Getenv("REPLICATE_API_TOKEN")
	if config.ReplicateAPIKey == "" {
		panic("REPLICATE_API_TOKEN not set")
	}

	// Initialize the Gin engine.
	r := gin.Default()

	// Set client options and connect to MongoDB
	mongo.Init(config.DBName)

	// Define a route for the index page.
	r.GET("/api/v1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Gin backend service with MongoDB!",
		})
	})

	// Define a route for a health check.
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// Add business routes to engine.
	addRoutes(r)

	// Start serving the application on port 8080.
	r.Run(":8080")
}

func addRoutes(e *gin.Engine) {
	g := e.Group(config.APIPrefix)
	g.POST("/session", session.Login)
	g.POST("/users", users.Register)
	g.POST("/album/new", album.AddPhoto)
	g.POST("/game/match/:level/new", game.Start)

	g.POST("/users/files/:filename", users.UploadFile)
	g.POST("/users/:name/files/:filename", users.UploadFile)

	// TODO: Deprecate the default file server; write something with authentication
	g.Static(config.APIFiles, config.PhotoDir)
}
