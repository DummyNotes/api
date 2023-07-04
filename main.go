package main

import (
	"github.com/dummynotes/notes/internal/middleware"
	"log"
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

	public := router.Group("/")

	public.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	notesGroup := router.Group("/notes")
	notesGroup.Use(middleware.ExtractPayload())

	notesGroup.GET("/logheaders", func(c *gin.Context) {
		for k, v := range c.Request.Header {
			log.Printf("%s: %v", k, v)
		}

		c.Status(200)
	})

	notesGroup.PUT("/", noteHandlers.CreateNote)
	notesGroup.GET("/", noteHandlers.ListNote)
	notesGroup.GET("/:id", noteHandlers.GetNote)
	notesGroup.DELETE("/:id", noteHandlers.DeleteNote)

	router.Run(":8080")
}
