package server

import (
	"net/http"
	"time"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterUserRoutes(r *gin.RouterGroup) {
	user := r.Group("/user", s.userMiddleware)
	{
		user.POST("/create-profile", s.createProfileHandler)
	}
}

type CreateProfileRequest struct {
	BirthDate string `json:"birthdate"`
	AvatarUrl string `json:"avatarUrl"`
	Fullname  string `json:"fullname" binding:"required"`
	Username  string `json:"username" binding:"required"`
}

func (s *Server) createProfileHandler(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return
	}

	req, err := util.BindJSON[CreateProfileRequest](c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	findUser, err := s.db.FindUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	if findUser != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username already exists"})
		return
	}

	user_id_uuid, err := uuid.Parse(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	layout := "2006-01-02"
	birthdate_time, err := time.Parse(layout, req.BirthDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	err = s.db.CreateUserProfile(birthdate_time, req.AvatarUrl, req.Fullname, req.Username, user_id_uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Profile created successfully"})
}
