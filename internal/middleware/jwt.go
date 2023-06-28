package middleware

import (
	"fmt"
	"net/http"

	"github.com/dummynotes/notes/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractPayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &models.User{}

		token, _, err := new(jwt.Parser).ParseUnverified(c.GetHeader("Authorization"), jwt.MapClaims{})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Something went wrong"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			user.UserID = fmt.Sprint(claims["userid"])
		}

		if user.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Something went wrong"})
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
