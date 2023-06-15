package notes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tiborhercz/notes/internal/config"
	"github.com/tiborhercz/notes/internal/database"
	"github.com/tiborhercz/notes/internal/models"
)

type Handlers struct {
	envConfig *config.EnvConfigStruct
	dbClient  *database.DynamoDBClient
}

func NewAPIHandlers(client *database.DynamoDBClient, envConfig *config.EnvConfigStruct) *Handlers {
	return &Handlers{
		dbClient:  client,
		envConfig: envConfig,
	}
}

type CreateNoteRequestBody struct {
	Title string `json:"Title" binding:"required"`
	Text  string `json:"Text" binding:"required"`
}

func (h *Handlers) CreateNote(c *gin.Context) {
	body := CreateNoteRequestBody{}

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	timeNow := time.Now().Format(time.RFC3339)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

type GetNoteRequestBody struct {
	NoteId string `json:"NoteId" binding:"required"`
}

func (h *Handlers) GetNote(c *gin.Context) {
	id := c.Param("id")

	response, err := h.dbClient.Get(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"note": response})
}
