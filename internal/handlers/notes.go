package handlers

import (
	"net/http"
	"time"

	"github.com/dummynotes/notes/internal/config"
	"github.com/dummynotes/notes/internal/database"
	"github.com/dummynotes/notes/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotesHandlers struct {
	envConfig *config.EnvConfigStruct
	dbClient  *database.DynamoDBClient
}

func NewNotesAPIHandlers(client *database.DynamoDBClient, envConfig *config.EnvConfigStruct) *NotesHandlers {
	return &NotesHandlers{
		dbClient:  client,
		envConfig: envConfig,
	}
}

type CreateNoteRequestBody struct {
	Title string `json:"Title" binding:"required"`
	Text  string `json:"Text" binding:"required"`
}

func (h *NotesHandlers) CreateNote(c *gin.Context) {
	body := CreateNoteRequestBody{}

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	timeNow := time.Now().UTC().Format(time.RFC3339)

	note := &models.Note{
		NoteId:    uuid.New().String(),
		UserId:    uuid.New().String(),
		Title:     body.Title,
		Text:      body.Text,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	_, err := h.dbClient.Create(note)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *NotesHandlers) GetNote(c *gin.Context) {
	id := c.Param("id")

	response, err := h.dbClient.Get(id)
	if err.Error() == "GetItem: Data not found." {
		c.Status(http.StatusNotFound)
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": response})
}

func (h *NotesHandlers) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	_, err := h.dbClient.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
