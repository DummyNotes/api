package main

import (
	"net/http"

	"github.com/dummynotes/notes/internal/config"
	"github.com/dummynotes/notes/internal/database"
	"github.com/dummynotes/notes/internal/handlers"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	envConfig := config.GetConfig()
	client := database.DynamoDBClient{}

	client.Configure(database.GetDynamodbClient(envConfig.Region))

	noteHandlers := handlers.NewNotesAPIHandlers(&client, envConfig)

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	notesGroup := router.Group("/notes")

	notesGroup.PUT("/", noteHandlers.CreateNote)
	notesGroup.GET("/:id", noteHandlers.GetNote)
	notesGroup.DELETE("/:id", noteHandlers.DeleteNote)

	router.Run(":8080")
}
