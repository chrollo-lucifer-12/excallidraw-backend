package server

import (
	"net/http"
	"strings"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/util"
	"github.com/gin-gonic/gin"
)

func (s *Server) userMiddleware(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	token := strings.Split(bearerToken, " ")[1]

	userId, err := util.ParseToken(token, "sahil")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	findUser, err := s.db.FindUserByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	if findUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.Set("user_id", userId)
	c.Next()
}
