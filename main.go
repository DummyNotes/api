package main

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/tiborhercz/notes/internal/config"
	"github.com/tiborhercz/notes/internal/database"
	"github.com/tiborhercz/notes/internal/notes"
)

func main() {
	envConfig := config.GetConfig()
	client := database.DynamoDBClient{}

	client.Configure(database.GetDynamodbClient(envConfig.Region))

	noteHandlers := notes.NewAPIHandlers(&client, envConfig)

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
