package server

import (
	"fmt"
	"net/http"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/util"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterWhiteboardRoutes(r *gin.RouterGroup) {
	whiteBoard := r.Group("/whiteboard", s.userMiddleware)
	{
		whiteBoard.POST("/create", s.createWhiteboardHandler)
	}
}

type CreateWhiteboardRequest struct {
	Name string `json:"name" binding:"required"`
}

func (s *Server) createWhiteboardHandler(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return
	}

	user_id_uuid, err := util.ParseUUID(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	req, err := util.BindJSON[CreateWhiteboardRequest](c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	slug, err := util.GenerateRandomSlug(3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	fmt.Println(slug)
	err = s.db.CreateWhiteboard(user_id_uuid, req.Name, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "whiteboard created", "slug": slug})
}
