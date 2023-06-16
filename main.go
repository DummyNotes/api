package main

import (
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

	router.PUT("/", noteHandlers.CreateNote)
	router.GET("/:id", noteHandlers.GetNote)
	router.DELETE("/:id", noteHandlers.DeleteNote)

	router.Run(":8080")
}
